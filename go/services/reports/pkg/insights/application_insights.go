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

type ApplicationLogsInsight struct {
	ApplicationName string   `json:"applicationName"`
	IncidentName    string   `json:"name"`
	Category        string   `json:"category"`
	Summary         string   `json:"summary"`
	Recommendation  string   `json:"recommendation"`
	SourceLogIds    []string `json:"sourceLogIds"`
}

type ApplicationInsightMetadata struct {
	Timestamp     int64  `json:"timestamp"`
	ContainerName string `json:"containerName"`
	PodName       string `json:"podName"`
	Source        string `json:"source"`
}

type ApplicationInsightsWithMetadata struct {
	Insight  *ApplicationLogsInsight      `json:"insight"`
	Metadata []ApplicationInsightMetadata `json:"metadata"`
}

type ApplicationInsightsGenerator interface {
	OnDemandApplicationInsights(logs []*repositories.ApplicationLogsDocument, configuration []*ApplicationInsightConfiguration) ([]ApplicationInsightsWithMetadata, error)
	ScheduledApplicationInsights(logs []*repositories.ApplicationLogsDocument, scheduledTime time.Time) ([]ApplicationLogsInsight, error)
}

type applicationInsightsResponseDto struct {
	Insights []ApplicationLogsInsight
}

type ApplicationInsightConfiguration struct {
	ApplicationName string `json:"applicationName"`
	Precision       string `json:"precision"`
	CustomPrompt    string `json:"customPrompt"`
}

func (g *OpenAiInsightsGenerator) getApplicationLogById(logId string, logs []*repositories.ApplicationLogsDocument) (*repositories.ApplicationLogsDocument, error) {

	firstById := array.FindFirst(func(log *repositories.ApplicationLogsDocument) bool {
		return log.Id == logId
	})

	first, isSome := option.Unwrap(firstById(logs))
	if !isSome {
		g.logger.Error("Failed to find log by id", zap.String("id", logId))
	}

	return first, nil
}

func (g *OpenAiInsightsGenerator) addMetadataToInsight(insight ApplicationLogsInsight, logs []*repositories.ApplicationLogsDocument) ApplicationInsightsWithMetadata {

	applicationInsightsMetadata := make([]ApplicationInsightMetadata, 0, len(insight.SourceLogIds))
	for _, sourceLogId := range insight.SourceLogIds {
		log, err := g.getApplicationLogById(sourceLogId, logs)
		if err != nil {
			g.logger.Error("Failed to source application insights", zap.Error(err))
		}

		applicationInsightsMetadata = append(applicationInsightsMetadata, ApplicationInsightMetadata{
			Timestamp:     log.Timestamp,
			ContainerName: log.ContainerName,
			PodName:       log.PodName,
			Source:        log.Content,
		},
		)
	}

	return ApplicationInsightsWithMetadata{
		Insight:  &insight,
		Metadata: applicationInsightsMetadata,
	}
}

func (g *OpenAiInsightsGenerator) OnDemandApplicationInsights(
	logs []*repositories.ApplicationLogsDocument,
	configurations []*ApplicationInsightConfiguration) ([]ApplicationInsightsWithMetadata, error) {

	groupedLogs := GroupLogsByName(logs)

	// Map report configuration for an app (precision/customPrompt) to a app name.
	configurationsByApplication := MapApplicationNameToConfiguration(configurations)

	insightsChannel := make(chan []ApplicationInsightsWithMetadata, len(groupedLogs))
	allInsights := make([]ApplicationInsightsWithMetadata, 0, len(groupedLogs))

	var wg sync.WaitGroup

	// Generate insights for each application separately.
	for applicationName, logs := range groupedLogs {
		wg.Add(1)
		go func() {
			defer wg.Done()

			insights, err := g.getInsightsForSingleApplication(
				logs,
				configurationsByApplication[applicationName],
			)

			if err != nil {
				g.logger.Error("Failed to get insights for an application", zap.Error(err), zap.String("app", applicationName))
			}

			// Add metadata about insights (container/pod)
			mapper := array.Map(func(insights ApplicationLogsInsight) ApplicationInsightsWithMetadata {
				return g.addMetadataToInsight(insights, logs)
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

func (g *OpenAiInsightsGenerator) getInsightsForSingleApplication(
	logs []*repositories.ApplicationLogsDocument,
	configuration *ApplicationInsightConfiguration) ([]ApplicationLogsInsight, error) {

	encodedLogs, err := json.Marshal(logs)
	if err != nil {
		g.logger.Error("Failed to encode application logs for application insights", zap.Error(err))
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
			Always declare a unmodified sources with every insight you give.  
			Always give a recommendation on how to resolve the issue. Always give a source. Never repeat insights, ie. 
			if you once use the source do not create an insight for it again. One insight per source. If you recognize the 
			same events on different containers/pods. 
			Add all logs (from all pods/containers) which belong to the same insight to the sourceIds array of a single insight.
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
		openai.CreateJsonReponseFormat("application_insights", applicationInsightsResponseDto{}),
	)

	if err != nil {
		g.logger.Error("Failed to get application logs insights from openai client", zap.Error(err))
		return nil, err
	}

	var insights applicationInsightsResponseDto

	err = json.Unmarshal([]byte(openAiResponse.Choices[0].Message.Content), &insights)
	if err != nil {
		g.logger.Error("Failed to decode application insights from openai client", zap.Error(err))
		return nil, err
	}

	return insights.Insights, nil

}

func MapApplicationNameToConfiguration(configurations []*ApplicationInsightConfiguration) map[string]*ApplicationInsightConfiguration {
	groupedConfigurations := make(map[string]*ApplicationInsightConfiguration)
	for _, conf := range configurations {
		groupedConfigurations[conf.ApplicationName] = conf
	}

	return groupedConfigurations
}

func (g *OpenAiInsightsGenerator) ScheduledApplicationInsights(logs []*repositories.ApplicationLogsDocument, scheduledTime time.Time) ([]ApplicationLogsInsight, error) {
	return nil, nil
}

var _ ApplicationInsightsGenerator = &OpenAiInsightsGenerator{}
