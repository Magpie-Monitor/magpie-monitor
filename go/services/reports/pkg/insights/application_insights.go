package insights

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IBM/fp-go/array"
	"github.com/IBM/fp-go/option"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/jsonl"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/openai"
	reportrepositories "github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/repositories"
	"go.uber.org/zap"
	"sync"
	"time"
)

type ApplicationLogsInsight struct {
	ApplicationName string                     `json:"applicationName"`
	IncidentName    string                     `json:"name"`
	Category        string                     `json:"category"`
	Summary         string                     `json:"summary"`
	Recommendation  string                     `json:"recommendation"`
	Urgency         reportrepositories.Urgency `json:"urgency"`
	SourceLogIds    []string                   `json:"sourceLogIds"`
}

type ApplicationInsightMetadata struct {
	Timestamp     int64  `json:"timestamp"`
	ContainerName string `json:"containerName"`
	PodName       string `json:"podName"`
	Image         string `json:"image"`
	Source        string `json:"source"`
}

type ApplicationInsightsWithMetadata struct {
	Insight  *ApplicationLogsInsight      `json:"insight"`
	Metadata []ApplicationInsightMetadata `json:"metadata"`
}

type ApplicationInsightsGenerator interface {
	OnDemandApplicationInsights(
		logs []*repositories.ApplicationLogsDocument,
		configuration []*reportrepositories.ApplicationInsightConfiguration) ([]ApplicationInsightsWithMetadata, error)

	ScheduleApplicationInsights(
		logs []*repositories.ApplicationLogsDocument,
		configuration []*reportrepositories.ApplicationInsightConfiguration,
		scheduledTime time.Time,
		cluster string,
		fromDate int64,
		toDate int64,
	) (*reportrepositories.ScheduledApplicationInsights, error)

	GetScheduledApplicationInsights(
		sheduledInsights *reportrepositories.ScheduledApplicationInsights,
	) ([]ApplicationInsightsWithMetadata, error)
}

type applicationInsightsResponseDto struct {
	Insights []ApplicationLogsInsight
}

func GroupInsightsByApplication(applicationInsights []ApplicationInsightsWithMetadata) map[string][]ApplicationInsightsWithMetadata {
	insightsByApplication := make(map[string][]ApplicationInsightsWithMetadata)

	for _, insight := range applicationInsights {
		applicationName := insight.Insight.ApplicationName
		insightsByApplication[applicationName] = append(insightsByApplication[applicationName], insight)
	}
	return insightsByApplication
}

func (g *OpenAiInsightsGenerator) getApplicationLogById(logId string, logs []*repositories.ApplicationLogsDocument) (*repositories.ApplicationLogsDocument, error) {

	if len(logs) == 0 {
		return nil, errors.New("Failed to find application log by id in an empty logs array")
	}

	firstById := array.FindFirst(func(log *repositories.ApplicationLogsDocument) bool {
		return log.Id == logId
	})

	first, isSome := option.Unwrap(firstById(logs))
	if !isSome {
		g.logger.Error("Failed to find log by id", zap.String("id", logId))
		return nil, errors.New(fmt.Sprintf("Failed to find log by id %s", logId))
	}

	return first, nil
}

func (g *OpenAiInsightsGenerator) addMetadataToApplicationInsight(
	insight ApplicationLogsInsight,
	logs []*repositories.ApplicationLogsDocument) ApplicationInsightsWithMetadata {

	applicationInsightsMetadata := make([]ApplicationInsightMetadata, 0, len(insight.SourceLogIds))
	for _, sourceLogId := range insight.SourceLogIds {
		log, err := g.getApplicationLogById(sourceLogId, logs)
		if err != nil {
			g.logger.Error("Failed to source application insights", zap.Error(err))

			// Skipping source in case of a fake id
			continue
		}

		applicationInsightsMetadata = append(applicationInsightsMetadata, ApplicationInsightMetadata{
			Timestamp:     log.Timestamp,
			ContainerName: log.ContainerName,
			PodName:       log.PodName,
			Source:        log.Content,
			Image:         log.Image,
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
	configurations []*reportrepositories.ApplicationInsightConfiguration) ([]ApplicationInsightsWithMetadata, error) {

	groupedLogs := GroupApplicationLogsByName(logs)

	// Map report configuration for an app (precision/customPrompt) to a app name.
	configurationsByApplication := reportrepositories.MapApplicationNameToConfiguration(configurations)

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
				return g.addMetadataToApplicationInsight(insights, logs)
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

// Get Application insights by grouped by application name
func (g *OpenAiInsightsGenerator) GetScheduledApplicationInsights(
	sheduledInsights *reportrepositories.ScheduledApplicationInsights,
) ([]ApplicationInsightsWithMetadata, error) {
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

	insightLogs, err := g.applicationLogsRepository.
		GetLogs(context.TODO(), sheduledInsights.Cluster,
			time.Unix(0, sheduledInsights.FromDateNs),
			time.Unix(0, sheduledInsights.ToDateNs))

	if err != nil {
		g.logger.Error("Failed to get application logs for scheduled insight")
		return nil, err
	}

	// Each jsonl entry contains insights for a single application
	res := []ApplicationInsightsWithMetadata{}
	for _, response := range responses {
		var applicationInsights applicationInsightsResponseDto
		if len(response.Response.Body.Choices) == 0 {
			return nil, errors.New("Failed to get insights from batch completion choices")
		}

		messageContent := response.Response.Body.Choices[0].Message.Content
		err := json.Unmarshal([]byte(messageContent), &applicationInsights)
		if err != nil {
			g.logger.Error("Failed to decode application insight", zap.Error(err))
			return nil, err
		}

		if len(response.Response.Body.Choices) == 0 {
			return nil, errors.New("Failed to get insights from batch completion message content")
		}

		insightsWithMetadata := array.Map(func(insight ApplicationLogsInsight) ApplicationInsightsWithMetadata {
			return g.addMetadataToApplicationInsight(insight, insightLogs)
		})(applicationInsights.Insights)

		res = append(res, insightsWithMetadata...)
	}

	return res, nil
}

func (g *OpenAiInsightsGenerator) ScheduleApplicationInsights(
	logs []*repositories.ApplicationLogsDocument,
	configuration []*reportrepositories.ApplicationInsightConfiguration,
	scheduledTime time.Time,
	cluster string,
	fromDateNs int64,
	toDateNs int64,
) (*reportrepositories.ScheduledApplicationInsights, error) {

	groupedLogs := GroupApplicationLogsByName(logs)
	configurationsByApplication := reportrepositories.MapApplicationNameToConfiguration(configuration)
	completionRequests := make([]*openai.CompletionRequest, 0, len(groupedLogs))

	// Generate insights for each application separately.
	for applicationName, logs := range groupedLogs {
		messages, err := g.createMessagesFromApplicationLogs(
			logs,
			configurationsByApplication[applicationName],
		)
		if err != nil {
			g.logger.Error("Failed to create messages from application logs", zap.Error(err), zap.String("app", applicationName))
			return nil, err
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

	return &reportrepositories.ScheduledApplicationInsights{
		Id:                       resp.Id,
		CreatedAt:                resp.CreatedAt,
		ExpiresAt:                resp.ExpiresAt,
		CompletedAt:              resp.CompletedAt,
		Cluster:                  cluster,
		FromDateNs:               fromDateNs,
		ToDateNs:                 toDateNs,
		ApplicationConfiguration: configuration,
	}, nil
}

func (g *OpenAiInsightsGenerator) createMessagesFromApplicationLogs(
	logs []*repositories.ApplicationLogsDocument,
	configuration *reportrepositories.ApplicationInsightConfiguration) ([]*openai.Message, error) {

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
			Always declare a unmodified sources with every insight you give.  
			Always give a recommendation on how to resolve the issue. Always give a source. Never repeat insights, ie. 
			if you once use the source do not create an insight for it again. One insight per source. If you recognize the 
			same events on different containers/pods. For each incident assign urgency as an integer number between 1 and 3.
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
	}

	return messages, nil
}

func (g *OpenAiInsightsGenerator) getInsightsForSingleApplication(
	logs []*repositories.ApplicationLogsDocument,
	configuration *reportrepositories.ApplicationInsightConfiguration) ([]ApplicationLogsInsight, error) {

	messages, err := g.createMessagesFromApplicationLogs(logs, configuration)
	if err != nil {
		g.logger.Error("Failed to create messages", zap.Error(err))
		return nil, err
	}

	openAiResponse, err := g.client.Complete(messages,
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

func GroupApplicationLogsByName(logs []*repositories.ApplicationLogsDocument) map[string][]*repositories.ApplicationLogsDocument {
	groupedLogs := make(map[string][]*repositories.ApplicationLogsDocument)
	for _, log := range logs {
		groupedLogs[log.ApplicationName] = append(groupedLogs[log.ApplicationName], log)
	}

	return groupedLogs
}

var _ ApplicationInsightsGenerator = &OpenAiInsightsGenerator{}
