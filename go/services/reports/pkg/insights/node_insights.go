package insights

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/IBM/fp-go/array"
	"github.com/IBM/fp-go/option"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/openai"
	"go.uber.org/zap"
)

type NodeLogsInsight struct {
	Name           string `json:"name"`
	Category       string `json:"category"`
	Summary        string `json:"summary"`
	Recommendation string `json:"recommendation"`
	SourceLogId    string `json:"sourceLogId"`
}

type NodeInsightMetadata struct {
	NodeName  string `json:"nodeName"`
	Timestamp int64  `json:"timestamp"`
	Source    string `json:"source"`
}

type NodeInsightsWithMetadata struct {
	Insight  *NodeLogsInsight     `json:"insight"`
	Metadata *NodeInsightMetadata `json:"metadata"`
}

type NodeInsightsGenerator interface {
	OnDemandNodeInsights(logs []*repositories.NodeLogsDocument, configuration []*NodeInsightConfiguration) ([]NodeInsightsWithMetadata, error)
	ScheduledNodeInsights(logs []*repositories.NodeLogsDocument, scheduledTime time.Time) ([]NodeLogsInsight, error)
}

type nodeInsightsResponseDto struct {
	Insights []NodeLogsInsight
}

type NodeInsightConfiguration struct {
	NodeName     string `json:"nodeName"`
	Precision    string `json:"precision"`
	CustomPrompt string `json:"customPrompt"`
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
	configurations []*NodeInsightConfiguration) ([]NodeInsightsWithMetadata, error) {

	groupedLogs := make(map[string][]*repositories.NodeLogsDocument)
	for _, log := range logs {
		groupedLogs[log.Name] = append(groupedLogs[log.Name], log)
	}

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
			}
			mapper := array.Map(func(insight NodeLogsInsight) NodeInsightsWithMetadata {

				log, err := g.getNodeLogById(insight.SourceLogId, logs)
				if err != nil {
					g.logger.Error("Failed to source node insights", zap.Error(err))
				}
				return NodeInsightsWithMetadata{
					Insight: &insight,
					Metadata: &NodeInsightMetadata{
						NodeName:  nodeName,
						Timestamp: log.Timestamp,
						Source:    log.Content,
					},
				}
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

	encodedLogs, err := json.Marshal(logs)
	if err != nil {
		g.logger.Error("Failed to encode node logs for application insights", zap.Error(err))
		return nil, err
	}

	customPrompt := ""
	if configuration != nil {
		customPrompt = configuration.CustomPrompt
	}

	openAiResponse, err := g.client.Complete([]*openai.Message{
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
			only mention the same issue once.
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
	},
		openai.CreateJsonReponseFormat("node_insights", applicationInsightsResponseDto{}),
	)

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

func MapNodeNameToConfiguration(configurations []*NodeInsightConfiguration) map[string]*NodeInsightConfiguration {
	groupedConfigurations := make(map[string]*NodeInsightConfiguration)
	for _, conf := range configurations {
		groupedConfigurations[conf.NodeName] = conf
	}

	return groupedConfigurations
}

func (g *OpenAiInsightsGenerator) ScheduledNodeInsights(logs []*repositories.NodeLogsDocument, scheduledTime time.Time) ([]NodeLogsInsight, error) {
	return nil, nil
}

var _ NodeInsightsGenerator = &OpenAiInsightsGenerator{}
