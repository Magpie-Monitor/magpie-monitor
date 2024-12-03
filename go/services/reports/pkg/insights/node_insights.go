package insights

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IBM/fp-go/array"
	"github.com/IBM/fp-go/option"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/openai"
	"go.uber.org/zap"
	"golang.org/x/exp/maps"
	"time"
)

type NodeLogsInsight struct {
	NodeName       string   `json:"nodeName"`
	Title          string   `json:"title"`
	Category       string   `json:"category"`
	Summary        string   `json:"summary"`
	Recommendation string   `json:"recommendation"`
	Urgency        Urgency  `json:"urgency"`
	SourceLogIds   []string `json:"sourceLogIds"`
}

type NodeInsightConfiguration struct {
	NodeName     string `json:"nodeName"`
	Accuracy     string `json:"accuracy"`
	CustomPrompt string `json:"customPrompt"`
}

type ScheduledNodeInsights struct {
	ScheduledJobIds   []string                    `bson:"scheduledJobIds" json:"scheduledJobIds"`
	SinceMs           int64                       `bson:"sinceMs" json:"sinceMs"`
	ToMs              int64                       `bson:"toMs" json:"toMs"`
	ClusterId         string                      `bson:"clusterId" json:"clusterId"`
	NodeConfiguration []*NodeInsightConfiguration `bson:"nodeConfiguration" json:"nodeConfiguration"`
}

type NodeInsightMetadata struct {
	ClusterId     string `json:"clusterId"`
	NodeName      string `json:"nodeName"`
	CollectedAtMs int64  `json:"collectedAtMs"`
	Filename      string `json:"filename"`
	SourceLogId   string `json:"sourceLogId"`
	SourceLog     string `json:"sourceLog"`
}

type NodeInsightsWithMetadata struct {
	Insight  *NodeLogsInsight      `json:"insight"`
	Metadata []NodeInsightMetadata `json:"metadata"`
}

type NodeInsightsGenerator interface {
	ScheduleNodeInsights(
		groupedLogs map[string][]*repositories.NodeLogsDocument,
		configurationByNode map[string]*NodeInsightConfiguration,
		cluster string,
		fromDate int64,
		toDate int64,
	) (*ScheduledNodeInsights, error)

	GetScheduledNodeInsights(
		sheduledInsights *ScheduledNodeInsights,
	) ([]NodeInsightsWithMetadata, error)

	AwaitScheduledNodeInsights(
		sheduledInsights *ScheduledNodeInsights,
	) ([]NodeInsightsWithMetadata, error)
}

type NodeInsightsResponseDto struct {
	Insights []NodeLogsInsight
}

func MapNodeNameToConfiguration(configurations []*NodeInsightConfiguration) map[string]*NodeInsightConfiguration {
	groupedConfigurations := make(map[string]*NodeInsightConfiguration)
	for _, conf := range configurations {
		groupedConfigurations[conf.NodeName] = conf
	}

	return groupedConfigurations
}

func (g *OpenAiInsightsGenerator) getNodeLogById(logId string, logs []*repositories.NodeLogsDocument) (*repositories.NodeLogsDocument, error) {

	firstById := array.FindFirst(func(log *repositories.NodeLogsDocument) bool {
		return log.Id == logId
	})

	first, isSome := option.Unwrap(firstById(logs))
	if !isSome {
		g.logger.Error("Failed to find log by id", zap.String("id", logId))
		return nil, errors.New("No logs by id")
	}

	return first, nil
}

func (g *OpenAiInsightsGenerator) getInsightsForSingleNode(
	logs []*repositories.NodeLogsDocument,
	configuration *NodeInsightConfiguration) ([]NodeLogsInsight, error) {

	messages, err := g.createMessagesFromNodeLogs(logs, configuration)
	if err != nil {
		g.logger.Error("Failed to create messages", zap.Error(err))
		return nil, err
	}

	openAiResponse, err := g.client.Complete(messages,
		openai.CreateJsonReponseFormat("node_insights", NodeInsightsResponseDto{}))

	if err != nil {
		g.logger.Error("Failed to get node logs insights from openai client", zap.Error(err))
		return nil, err
	}

	var insights NodeInsightsResponseDto

	if len(openAiResponse.Choices) == 0 {
		g.logger.Error("No insight choices were returned for", zap.Any("node", configuration.NodeName))
		return nil, errors.New(fmt.Sprintf("No insight choices were returned for %s", configuration.NodeName))
	}

	err = json.Unmarshal([]byte(openAiResponse.Choices[0].Message.Content), &insights)
	if err != nil {
		g.logger.Error("Failed to decode node insights from openai client", zap.Error(err))
		return nil, err
	}

	return insights.Insights, nil
}

func (g *OpenAiInsightsGenerator) ScheduleNodeInsights(
	groupedLogs map[string][]*repositories.NodeLogsDocument,
	configurationByNode map[string]*NodeInsightConfiguration,
	clusterId string,
	sinceMs int64,
	toMs int64,
) (*ScheduledNodeInsights, error) {

	completionRequests := make(map[string]*openai.CompletionRequest, len(groupedLogs))

	// Generate insights for each application separately.
	for nodeName, logs := range groupedLogs {

		// Split completion request to packets not greater than OpenAi model's context
		logPackets := repositories.SplitLogsIntoPackets(logs, g.client.ContextSizeBytes)

		for idx, logPacket := range logPackets {
			messages, err := g.createMessagesFromNodeLogs(
				logPacket,
				configurationByNode[nodeName],
			)
			if err != nil {
				g.logger.Error("Failed to messages from logs", zap.Error(err), zap.String("node", nodeName))
			}

			// Assign each request a customId with a nodeName and a packetId
			completionRequests[fmt.Sprintf("%s-%d", nodeName, idx)] = &openai.CompletionRequest{
				Messages:       messages,
				Temperature:    g.client.Temperature,
				ResponseFormat: openai.CreateJsonReponseFormat("insights", NodeInsightsResponseDto{}),
				Model:          g.client.Model(),
			}
		}
	}

	completionRequestsPerJob, err := g.client.SplitCompletionReqestsByBatchSize(completionRequests)
	if err != nil {
		g.logger.Error("Failed to split node insights completion requests", zap.Error(err))
		return nil, err
	}

	scheduledJobs := array.Map(openai.NewOpenAiJob)(completionRequestsPerJob)

	ids, repositoryErr := g.scheduledJobsRepository.InsertScheduledJobs(context.Background(), scheduledJobs)
	if repositoryErr != nil {
		g.logger.Error("Failed to split application insights completion requests", zap.Error(err))
		return nil, err
	}

	nodeConfigurationList := maps.Values(configurationByNode)

	return &ScheduledNodeInsights{
		ScheduledJobIds:   ids,
		ClusterId:         clusterId,
		SinceMs:           sinceMs,
		ToMs:              toMs,
		NodeConfiguration: nodeConfigurationList,
	}, nil
}

func (g *OpenAiInsightsGenerator) createMessagesFromNodeLogs(
	logs []*repositories.NodeLogsDocument,
	configuration *NodeInsightConfiguration) ([]*openai.Message, error) {

	encodedLogs, err := json.Marshal(logs)
	if err != nil {
		g.logger.Error("Failed to encode application logs for application insights", zap.Error(err))
		return nil, err
	}

	customPrompt := ""
	if configuration != nil {
		customPrompt = configuration.CustomPrompt
	}

	messages := []*openai.Message{
		{
			Role: "system",
			Content: fmt.Sprintf(`You are a Kubernetes cluster system administrator.
        Analyze the provided Kubernetes node logs and identify any meaningful insights indicating potential errors or issues. 
        For each identified issue, provide:
        - A concise title (few words) summarizing the insight.
        - A detailed summary (max 50 words) that includes the probable cause, category of the issue, and context.
        - An urgency level (LOW, MEDIUM, HIGH).
        - A recommended action for resolution.

        Each insight should:
        - Include relevant, unmodified source log entries in the sourceIds array, aggregating similar logs across nodes if they pertain to the same issue. 
       	 	Use _id field of provided log as its id.
        - Avoid redundant insights; if a source has already been referenced, do not repeat it in another insight.
        - Ignore logs that do not clearly suggest an issue or represent routine node activities.
        - Exclude insights entirely if no errors or warnings are present.

        Here is the additional configuration to consider when generating insights: %s`, customPrompt),
		},
		{
			Role: "user",
			Content: fmt.Sprintf(`These are logs from my Kubernetes nodes.
        Please identify and describe any issues or anomalies indicated by the logs:
        %s`, encodedLogs),
		},
	}

	return messages, nil
}

// TODO: Remove once the AwaitScheduledNodeInsights is working
func (g *OpenAiInsightsGenerator) GetScheduledNodeInsights(
	scheduledInsights *ScheduledNodeInsights,
) ([]NodeInsightsWithMetadata, error) {

	batches, err := g.client.Batches(scheduledInsights.ScheduledJobIds)
	if err != nil {
		g.logger.Error("Failed to get batch from id", zap.Error(err))
		return nil, err
	}

	completionResponses, err := g.client.CompletionResponseEntriesFromBatches(batches)
	if err != nil {
		g.logger.Error("Failed to get node batch completion responses", zap.Error(err))
		return nil, err
	}

	insights, err := g.GetNodeInsightsFromBatchEntries(completionResponses, scheduledInsights)
	if err != nil {
		g.logger.Error("Failed to transform batch entries into node insights")
		return nil, err
	}

	return insights, nil
}

func (g *OpenAiInsightsGenerator) AwaitScheduledNodeInsights(
	scheduledInsights *ScheduledNodeInsights,
) ([]NodeInsightsWithMetadata, error) {

	jobs, failedBatches, err := g.batchPoller.AwaitPendingJobs(scheduledInsights.ScheduledJobIds)
	if err != nil {
		g.logger.Error("Failed to get batch from id", zap.Error(err))
		return nil, err
	}

	// Ignoring failed batches
	if len(failedBatches) > 0 {
		g.logger.Error("Some of the node batches have failed", zap.Any("failedBatches", failedBatches))
	}

	batches, err := g.batchPoller.BatchesFromJobs(jobs)
	if err != nil {
		g.logger.Error("Failed to fetch node batches from scheduled jobs", zap.Error(err))
		return nil, err
	}

	completionResponses, err := g.client.CompletionResponseEntriesFromBatches(batches)
	if err != nil {
		g.logger.Error("Failed to get node batch completion responses", zap.Error(err))
		return nil, err
	}

	insights, err := g.GetNodeInsightsFromBatchEntries(completionResponses, scheduledInsights)
	if err != nil {
		g.logger.Error("Failed to transform batch entries into node insights")
		return nil, err
	}

	return insights, nil
}

func (g OpenAiInsightsGenerator) GetNodeInsightsFromBatchEntries(
	batchEntries []*openai.BatchFileCompletionResponseEntry,
	scheduledNodeInsights *ScheduledNodeInsights) ([]NodeInsightsWithMetadata, error) {

	// Each jsonl entry contains insights for a single node
	insights := []NodeInsightsWithMetadata{}
	for _, response := range batchEntries {
		var nodeInsights NodeInsightsResponseDto
		if len(response.Response.Body.Choices) == 0 {
			return nil, errors.New("Failed to get insights from batch completion choices")
		}

		messageContent := response.Response.Body.Choices[0].Message.Content
		err := json.Unmarshal([]byte(messageContent), &nodeInsights)
		if err != nil {
			//OpenAI returned incorrectly structured output (skipping to minimize impact)
			g.logger.Error("OpenAI returned incorrectly formatted node insights output", zap.Error(err), zap.Any("content", messageContent))
			continue
		}

		insightsWithMetadata := make([]NodeInsightsWithMetadata, 0, 0)
		for _, insight := range nodeInsights.Insights {

			insightWithMetadata, err := g.addMetadataToNodeInsight(insight, scheduledNodeInsights)
			if err != nil {
				g.logger.Error("Failed to add metadata to node insight, skipping")
				continue
			}
			insightsWithMetadata = append(insightsWithMetadata, *insightWithMetadata)
		}

		insights = append(insights, insightsWithMetadata...)
	}

	return insights, nil
}

func (g *OpenAiInsightsGenerator) addMetadataToNodeInsight(
	insight NodeLogsInsight,
	scheduledNodeInsights *ScheduledNodeInsights,
) (*NodeInsightsWithMetadata, error) {

	nodeInsightsMetadata := make([]NodeInsightMetadata, 0, len(insight.SourceLogIds))

	logs, err := g.nodeLogsRepository.GetLogsByIds(context.Background(),
		scheduledNodeInsights.ClusterId,
		time.UnixMilli(scheduledNodeInsights.SinceMs),
		time.UnixMilli(scheduledNodeInsights.ToMs),
		insight.SourceLogIds)

	if err != nil {
		g.logger.Error("Failed to get logs by ids to insight metadata")
		return nil, err
	}

	for _, sourceLogId := range insight.SourceLogIds {

		log, err := g.getNodeLogById(sourceLogId, logs)
		if err != nil {
			g.logger.Error("Failed to source node insights", zap.Error(err))

			// In case of fake ids skip the source
			continue
		}

		nodeInsightsMetadata = append(nodeInsightsMetadata, NodeInsightMetadata{
			CollectedAtMs: log.CollectedAtMs,
			NodeName:      log.Name,
			SourceLog:     log.Content,
			SourceLogId:   log.Id,
			ClusterId:     log.ClusterId,
			Filename:      log.Filename,
		},
		)
	}

	return &NodeInsightsWithMetadata{
		Insight:  &insight,
		Metadata: nodeInsightsMetadata,
	}, nil
}

func GroupInsightsByNode(nodeInsights []NodeInsightsWithMetadata) map[string][]NodeInsightsWithMetadata {
	insightsByNode := make(map[string][]NodeInsightsWithMetadata)

	for _, insight := range nodeInsights {
		if len(insight.Metadata) > 0 {
			nodeName := insight.Metadata[0].NodeName
			insightsByNode[nodeName] = append(insightsByNode[nodeName], insight)
		}
	}
	return insightsByNode
}

func GroupNodeConfigurationByName(nodeConfigurations []*NodeInsightConfiguration) map[string]*NodeInsightConfiguration {
	nodeConfigurationsByName := make(map[string]*NodeInsightConfiguration)

	for _, configuration := range nodeConfigurations {
		nodeConfigurationsByName[configuration.NodeName] = configuration
	}

	return nodeConfigurationsByName
}

var _ NodeInsightsGenerator = &OpenAiInsightsGenerator{}
