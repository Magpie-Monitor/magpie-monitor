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
	"sync"
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
	Source        string `json:"source"`
}

type NodeInsightsWithMetadata struct {
	Insight  *NodeLogsInsight      `json:"insight"`
	Metadata []NodeInsightMetadata `json:"metadata"`
}

type NodeInsightsGenerator interface {
	OnDemandNodeInsights(
		logs []*repositories.NodeLogsDocument,
		configuration []*NodeInsightConfiguration) ([]NodeInsightsWithMetadata, error)

	ScheduleNodeInsights(
		logs []*repositories.NodeLogsDocument,
		configuration []*NodeInsightConfiguration,
		scheduledTime time.Time,
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

type nodeInsightsResponseDto struct {
	Insights []NodeLogsInsight
}

func MapNodeNameToConfiguration(configurations []*NodeInsightConfiguration) map[string]*NodeInsightConfiguration {
	groupedConfigurations := make(map[string]*NodeInsightConfiguration)
	for _, conf := range configurations {
		groupedConfigurations[conf.NodeName] = conf
	}

	return groupedConfigurations
}

func FilterByNodesAccuracy(logsByNode map[string][]*repositories.NodeLogsDocument, configurationByNode map[string]*NodeInsightConfiguration) {
	for node, logs := range logsByNode {
		config, ok := configurationByNode[node]
		var accuracy Accuracy
		if !ok {
			// By default the node is not included
			delete(logsByNode, node)
		} else {
			accuracy = config.Accuracy
			filter := NewAccuracyFilter[*repositories.NodeLogsDocument](accuracy)
			logsByNode[node] = filter.Filter(logs)
		}
	}
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

func (g *OpenAiInsightsGenerator) OnDemandNodeInsights(
	logs []*repositories.NodeLogsDocument,
	configurations []*NodeInsightConfiguration) ([]NodeInsightsWithMetadata, error) {

	groupedLogs := GroupNodeLogsByName(logs)
	configurationsByNode := MapNodeNameToConfiguration(configurations)

	var wg sync.WaitGroup

	allInsights := make([]NodeInsightsWithMetadata, 0, len(groupedLogs))
	insightsChannel := make(chan []NodeInsightsWithMetadata, len(groupedLogs))

	for nodeName, logs := range groupedLogs {
		wg.Add(1)
		go func() {

			defer wg.Done()

			insights, err := g.getInsightsForSingleNode(logs, configurationsByNode[nodeName])
			if err != nil {
				g.logger.Error("Failed to get insights for an node", zap.Error(err), zap.String("node", nodeName))
				return
			}

			mapper := array.Map(func(insight NodeLogsInsight) NodeInsightsWithMetadata {
				return g.addMetadataToNodeInsight(insight, logs)
			})

			insightsChannel <- mapper(insights)
		}()
	}

	wg.Wait()

	close(insightsChannel)

	for insights := range insightsChannel {
		allInsights = append(allInsights, insights...)
	}

	return allInsights, nil
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
		openai.CreateJsonReponseFormat("node_insights", nodeInsightsResponseDto{}))

	if err != nil {
		g.logger.Error("Failed to get node logs insights from openai client", zap.Error(err))
		return nil, err
	}

	var insights nodeInsightsResponseDto

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
	logs []*repositories.NodeLogsDocument,
	configuration []*NodeInsightConfiguration,
	scheduledTime time.Time,
	clusterId string,
	sinceMs int64,
	toMs int64,
) (*ScheduledNodeInsights, error) {

	groupedLogs := GroupNodeLogsByName(logs)
	configurationsByNode := MapNodeNameToConfiguration(configuration)
	completionRequests := make([]*openai.CompletionRequest, 0, len(groupedLogs))

	// In place filter based on node configuration
	FilterByNodesAccuracy(groupedLogs, configurationsByNode)

	// Generate insights for each application separately.
	for nodeName, logs := range groupedLogs {

		// Split completion request to packets not greater than OpenAi model's context
		logPackets := repositories.SplitLogsIntoPackets(logs, g.client.ContextSizeBytes)

		for _, logPacket := range logPackets {
			messages, err := g.createMessagesFromNodeLogs(
				logPacket,
				configurationsByNode[nodeName],
			)
			if err != nil {
				g.logger.Error("Failed to messages from logs", zap.Error(err), zap.String("node", nodeName))
			}

			completionRequests = append(completionRequests,
				&openai.CompletionRequest{
					Messages:       messages,
					Temperature:    g.client.Temperature,
					ResponseFormat: openai.CreateJsonReponseFormat("insights", nodeInsightsResponseDto{}),
					Model:          g.client.Model(),
				},
			)
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

	return &ScheduledNodeInsights{
		ScheduledJobIds:   ids,
		ClusterId:         clusterId,
		SinceMs:           sinceMs,
		ToMs:              toMs,
		NodeConfiguration: configuration,
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
			Content: fmt.Sprintf(`You are a kubernetes cluster system administrator. 
			Given a list of logs from a Kubernetes cluster
			find logs which might suggest any kind of errors or issues. Try to give a possible reason, 
			category of an issue, urgency and possible resolution.   
			Source is an fragment of a the provided log that you are referencing in summary and recommendation. 
			Always declare a unmodified source log with every insight you give.  Title is a few word summary of the insight.
			Summary itself might be longer (max 50 words).
			Always give a recommendation on how to resolve the issue. Always give a source. Never repeat insights, ie. 
			if you once use the source do not create an insight for it again. One insight per source. Do not duplicate insights, 
			only mention the same issue once. For each incident assign urgency as an integer number between 0 and 2. 
			Ignore logs which do not explicitly suggest an issue. Ignore logs which are describing usual actions.
			If there are no errors or warnings don't even mention an insight. Here is the additional configuration 
			that you should consider while generating insights %s`, customPrompt),
		},
		{
			Role: "user",
			Content: fmt.Sprintf(`These are logs from my cluster. 
			Please tell me if they might suggest any kind of issues:
			%s`, encodedLogs),
		},
	}

	return messages, nil
}

// TODO: Remove once the AwaitScheduledNodeInsights is working
func (g *OpenAiInsightsGenerator) GetScheduledNodeInsights(
	sheduledInsights *ScheduledNodeInsights,
) ([]NodeInsightsWithMetadata, error) {

	batches, err := g.client.Batches(sheduledInsights.ScheduledJobIds)
	if err != nil {
		g.logger.Error("Failed to get batch from id", zap.Error(err))
		return nil, err
	}

	completionResponses, err := g.client.CompletionResponseEntriesFromBatches(batches)
	if err != nil {
		g.logger.Error("Failed to get node batch completion responses", zap.Error(err))
		return nil, err
	}

	insightLogs, err := g.nodeLogsRepository.
		GetLogs(context.TODO(), sheduledInsights.ClusterId,
			time.UnixMilli(sheduledInsights.SinceMs),
			time.UnixMilli(sheduledInsights.ToMs))

	if err != nil {
		g.logger.Error("Failed to get application logs for scheduled insight")
		return nil, err
	}

	insights, err := g.getNodeInsightsFromBatchEntries(completionResponses, insightLogs)
	if err != nil {
		g.logger.Error("Failed to transform batch entries into node insights")
		return nil, err
	}

	return insights, nil
}

func (g *OpenAiInsightsGenerator) AwaitScheduledNodeInsights(
	sheduledInsights *ScheduledNodeInsights,
) ([]NodeInsightsWithMetadata, error) {

	jobs, failedBatches, err := g.batchPoller.AwaitPendingJobs(sheduledInsights.ScheduledJobIds)
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

	insightLogs, err := g.nodeLogsRepository.
		GetLogs(context.TODO(), sheduledInsights.ClusterId,
			time.UnixMilli(sheduledInsights.SinceMs),
			time.UnixMilli(sheduledInsights.ToMs))

	if err != nil {
		g.logger.Error("Failed to get application logs for scheduled insight")
		return nil, err
	}

	insights, err := g.getNodeInsightsFromBatchEntries(completionResponses, insightLogs)
	if err != nil {
		g.logger.Error("Failed to transform batch entries into node insights")
		return nil, err
	}

	return insights, nil
}

func (g OpenAiInsightsGenerator) getNodeInsightsFromBatchEntries(batchEntries []*openai.BatchFileCompletionResponseEntry, logs []*repositories.NodeLogsDocument) ([]NodeInsightsWithMetadata, error) {

	// Each jsonl entry contains insights for a single node
	insights := []NodeInsightsWithMetadata{}
	for _, response := range batchEntries {
		var nodeInsights nodeInsightsResponseDto
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

		insightsWithMetadata := array.Map(func(insight NodeLogsInsight) NodeInsightsWithMetadata {
			return g.addMetadataToNodeInsight(insight, logs)
		})(nodeInsights.Insights)

		insights = append(insights, insightsWithMetadata...)
	}

	return insights, nil
}

func (g *OpenAiInsightsGenerator) addMetadataToNodeInsight(
	insight NodeLogsInsight,
	logs []*repositories.NodeLogsDocument) NodeInsightsWithMetadata {

	nodeInsightsMetadata := make([]NodeInsightMetadata, 0, len(insight.SourceLogIds))
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
			Source:        log.Content,
			ClusterId:     log.ClusterId,
			Filename:      log.Filename,
		},
		)
	}

	return NodeInsightsWithMetadata{
		Insight:  &insight,
		Metadata: nodeInsightsMetadata,
	}
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

func GroupNodeLogsByName(logs []*repositories.NodeLogsDocument) map[string][]*repositories.NodeLogsDocument {
	groupedLogs := make(map[string][]*repositories.NodeLogsDocument)
	for _, log := range logs {
		groupedLogs[log.Name] = append(groupedLogs[log.Name], log)
	}

	return groupedLogs
}

var _ NodeInsightsGenerator = &OpenAiInsightsGenerator{}
