package insights

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/IBM/fp-go/array"
	"github.com/IBM/fp-go/option"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/jsonl"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/openai"
	reportrepo "github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/repositories"
	"go.uber.org/zap"
)

type NodeLogsInsight struct {
	Name           string             `json:"name"`
	Category       string             `json:"category"`
	Summary        string             `json:"summary"`
	Recommendation string             `json:"recommendation"`
	Urgency        reportrepo.Urgency `json:"urgency"`
	SourceLogIds   []string           `json:"sourceLogIds"`
}

type NodeInsightMetadata struct {
	ClusterId string `json:"clusterId"`
	NodeName  string `json:"nodeName"`
	Timestamp int64  `json:"timestamp"`
	Source    string `json:"source"`
}

type NodeInsightsWithMetadata struct {
	Insight  *NodeLogsInsight      `json:"insight"`
	Metadata []NodeInsightMetadata `json:"metadata"`
}

type NodeInsightsGenerator interface {
	OnDemandNodeInsights(
		logs []*repositories.NodeLogsDocument,
		configuration []*reportrepo.NodeInsightConfiguration) ([]NodeInsightsWithMetadata, error)

	ScheduleNodeInsights(
		logs []*repositories.NodeLogsDocument,
		configuration []*reportrepo.NodeInsightConfiguration,
		scheduledTime time.Time,
		cluster string,
		fromDate int64,
		toDate int64,
	) (*reportrepo.ScheduledNodeInsights, error)

	GetScheduledNodeInsights(
		sheduledInsights *reportrepo.ScheduledNodeInsights,
	) ([]NodeInsightsWithMetadata, error)
}

type nodeInsightsResponseDto struct {
	Insights []NodeLogsInsight
}

func (g *OpenAiInsightsGenerator) getNodeLogById(logId string, logs []*repositories.NodeLogsDocument) (*repositories.NodeLogsDocument, error) {

	firstById := array.FindFirst(func(log *repositories.NodeLogsDocument) bool {
		return log.Id == logId
	})

	first, isSome := option.Unwrap(firstById(logs))
	if !isSome {
		g.logger.Error("Failed to find log by id", zap.String("id", logId))
	}

	return first, nil
}

func (g *OpenAiInsightsGenerator) OnDemandNodeInsights(
	logs []*repositories.NodeLogsDocument,
	configurations []*reportrepo.NodeInsightConfiguration) ([]NodeInsightsWithMetadata, error) {

	groupedLogs := GroupNodeLogsByName(logs)
	configurationsByNode := reportrepo.MapNodeNameToConfiguration(configurations)

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
	configuration *reportrepo.NodeInsightConfiguration) ([]NodeLogsInsight, error) {

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

	err = json.Unmarshal([]byte(openAiResponse.Choices[0].Message.Content), &insights)
	if err != nil {
		g.logger.Error("Failed to decode node insights from openai client", zap.Error(err))
		return nil, err
	}

	return insights.Insights, nil
}

func (g *OpenAiInsightsGenerator) ScheduleNodeInsights(
	logs []*repositories.NodeLogsDocument,
	configuration []*reportrepo.NodeInsightConfiguration,
	scheduledTime time.Time,
	clusterId string,
	sinceNano int64,
	toNano int64,
) (*reportrepo.ScheduledNodeInsights, error) {

	groupedLogs := GroupNodeLogsByName(logs)
	configurationsByApplication := reportrepo.MapNodeNameToConfiguration(configuration)
	completionRequests := make([]*openai.CompletionRequest, 0, len(groupedLogs))

	// Generate insights for each application separately.
	for nodeName, logs := range groupedLogs {
		messages, err := g.createMessagesFromNodeLogs(
			logs,
			configurationsByApplication[nodeName],
		)
		if err != nil {
			g.logger.Error("Failed to messages from logs", zap.Error(err), zap.String("node", nodeName))
		}

		completionRequests = append(completionRequests,
			&openai.CompletionRequest{
				Messages:       messages,
				Temperature:    0.6,
				ResponseFormat: openai.CreateJsonReponseFormat("insigts", applicationInsightsResponseDto{}),
				Model:          g.client.Model(),
			},
		)
	}

	resp, err := g.client.UploadAndCreateBatch(completionRequests)
	if err != nil {
		g.logger.Error("Failed to create a batch", zap.Error(err))
		return nil, err
	}

	return &reportrepo.ScheduledNodeInsights{
		Id:                resp.Id,
		ClusterId:         clusterId,
		SinceNano:         sinceNano,
		ToNano:            toNano,
		NodeConfiguration: configuration,
	}, nil
}

func (g *OpenAiInsightsGenerator) createMessagesFromNodeLogs(
	logs []*repositories.NodeLogsDocument,
	configuration *reportrepo.NodeInsightConfiguration) ([]*openai.Message, error) {

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
			Always declare a unmodified source log with every insight you give.  
			Always give a recommendation on how to resolve the issue. Always give a source. Never repeat insights, ie. 
			if you once use the source do not create an insight for it again. One insight per source. Do not duplicate insights, 
			only mention the same issue once. For each incident assign urgency as an integer number between 1 and 3. 
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

func (g *OpenAiInsightsGenerator) GetScheduledNodeInsights(
	sheduledInsights *reportrepo.ScheduledNodeInsights,
) ([]NodeInsightsWithMetadata, error) {

	batch, err := g.client.Batch(sheduledInsights.Id)
	if err != nil {
		g.logger.Error("Failed to get batch from id", zap.Error(err))
		return nil, err
	}

	outputFile, err := g.client.File(batch.OutputFileId)
	if err != nil {
		g.logger.Error("Failed to get output file from batch")
		return nil, err
	}

	var responses []openai.BatchFileCompletionResponseEntry
	g.logger.Sugar().Debugf("RESPONSE %s", string(outputFile))
	err = jsonl.NewJsonLinesDecoder(bytes.NewReader(outputFile)).Decode(&responses)
	if err != nil {
		g.logger.Error("Failed to decode application response", zap.Error(err))
		return nil, err
	}

	insightLogs, err := g.nodeLogsRepository.
		GetLogs(context.TODO(), sheduledInsights.ClusterId,
			time.Unix(0, sheduledInsights.SinceNano),
			time.Unix(0, sheduledInsights.ToNano))

	if err != nil {
		g.logger.Error("Failed to get application logs for scheduled insight")
		return nil, err
	}

	insights, err := g.getNodeInsightsFromBatchEntries(responses, insightLogs)
	if err != nil {
		g.logger.Error("Failed to transform batch entries into node insights")
		return nil, err
	}

	return insights, nil
}

func (g OpenAiInsightsGenerator) getNodeInsightsFromBatchEntries(batchEntries []openai.BatchFileCompletionResponseEntry, logs []*repositories.NodeLogsDocument) ([]NodeInsightsWithMetadata, error) {

	// Each jsonl entry contains insights for a single application
	insights := []NodeInsightsWithMetadata{}
	for _, response := range batchEntries {
		var nodeInsights nodeInsightsResponseDto
		if len(response.Response.Body.Choices) == 0 {
			return nil, errors.New("Failed to get insights from batch completion choices")
		}
		messageContent := response.Response.Body.Choices[0].Message.Content
		err := json.Unmarshal([]byte(messageContent), &nodeInsights)
		if err != nil {
			g.logger.Error("Failed to decode node insight", zap.Error(err))
			return nil, err
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
			Timestamp: log.Timestamp,
			NodeName:  log.Name,
			Source:    log.Content,
			ClusterId: log.Cluster,
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
		nodeName := insight.Metadata[0].NodeName
		insightsByNode[nodeName] = append(insightsByNode[nodeName], insight)
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
