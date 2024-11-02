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

type ApplicationInsightConfiguration struct {
	ApplicationName string   `json:"applicationName"`
	Accuracy        Accuracy `json:"accuracy"`
	CustomPrompt    string   `json:"customPrompt"`
}

type ScheduledApplicationInsights struct {
	ScheduledJobIds          []string                           `json:"scheduledJobIds"`
	SinceMs                  int64                              `bson:"sinceMs" json:"sinceMs"`
	ToMs                     int64                              `bson:"toMs" json:"toMs"`
	ClusterId                string                             `bson:"clusterId" json:"clusterId"`
	ApplicationConfiguration []*ApplicationInsightConfiguration `json:"applicationConfiguration"`
}

type ApplicationLogsInsight struct {
	ApplicationName string   `json:"applicationName"`
	IncidentName    string   `json:"name"`
	Category        string   `json:"category"`
	Summary         string   `json:"summary"`
	Recommendation  string   `json:"recommendation"`
	Urgency         Urgency  `json:"urgency"`
	SourceLogIds    []string `json:"sourceLogIds"`
}

type ApplicationInsightMetadata struct {
	CollectedAtMs   int64  `json:"collectedAtMs"`
	ApplicationName string `json:"applicationName"`
	ClusterId       string `json:"clusterId"`
	ContainerName   string `json:"containerName"`
	PodName         string `json:"podName"`
	Image           string `json:"image"`
	Source          string `json:"source"`
}

type ApplicationInsightsWithMetadata struct {
	Insight  *ApplicationLogsInsight      `json:"insight"`
	Metadata []ApplicationInsightMetadata `json:"metadata"`
}

type ApplicationInsightsGenerator interface {
	OnDemandApplicationInsights(
		logs []*repositories.ApplicationLogsDocument,
		configuration []*ApplicationInsightConfiguration) ([]ApplicationInsightsWithMetadata, error)

	ScheduleApplicationInsights(
		logs []*repositories.ApplicationLogsDocument,
		configuration []*ApplicationInsightConfiguration,
		scheduledTime time.Time,
		cluster string,
		fromDate int64,
		toDate int64,
	) (*ScheduledApplicationInsights, error)

	GetScheduledApplicationInsights(
		sheduledInsights *ScheduledApplicationInsights,
	) ([]ApplicationInsightsWithMetadata, error)

	AwaitScheduledApplicationInsights(
		sheduledInsights *ScheduledApplicationInsights,
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

func FilterByApplicationsAccuracy(logsByApplication map[string][]*repositories.ApplicationLogsDocument, configurationByApplication map[string]*ApplicationInsightConfiguration) {
	for application, logs := range logsByApplication {
		// Check if application configuration is in params
		config, ok := configurationByApplication[application]
		var accuracy Accuracy
		if !ok {
			// By default the app has low accuracy
			accuracy = Accuracy__Low
		} else {
			accuracy = config.Accuracy
		}

		filter := NewAccuracyFilter[*repositories.ApplicationLogsDocument](accuracy)
		logsByApplication[application] = filter.Filter(logs)
	}
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

func MapApplicationNameToConfiguration(configurations []*ApplicationInsightConfiguration) map[string]*ApplicationInsightConfiguration {
	groupedConfigurations := make(map[string]*ApplicationInsightConfiguration)
	for _, conf := range configurations {
		groupedConfigurations[conf.ApplicationName] = conf
	}

	return groupedConfigurations
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
			ApplicationName: log.ApplicationName,
			ClusterId:       log.ClusterId,
			CollectedAtMs:   log.CollectedAtMs,
			ContainerName:   log.ContainerName,
			PodName:         log.PodName,
			Source:          log.Content,
			Image:           log.Image,
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
	groupedLogs := GroupApplicationLogsByName(logs)

	// Map report configuration for an app (accuracy/customPrompt) to a app name.
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
				return
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

// Get Application insights grouped by application name
func (g *OpenAiInsightsGenerator) GetScheduledApplicationInsights(
	sheduledInsights *ScheduledApplicationInsights,
) ([]ApplicationInsightsWithMetadata, error) {
	batches, err := g.client.Batches(sheduledInsights.ScheduledJobIds)

	if err != nil {
		g.logger.Error("Failed to get batch from id", zap.Error(err))
		return nil, err
	}

	completionResponses, err := g.client.CompletionResponseEntriesFromBatches(batches)
	if err != nil {
		g.logger.Error("Failed to get application completion responses from batches", zap.Error(err))
		return nil, err
	}

	insightLogs, err := g.applicationLogsRepository.
		GetLogs(context.TODO(), sheduledInsights.ClusterId,
			time.UnixMilli(sheduledInsights.SinceMs),
			time.UnixMilli(sheduledInsights.ToMs))

	if err != nil {
		g.logger.Error("Failed to get application logs for scheduled insight")
		return nil, err
	}

	insights, err := g.getApplicationInsightsFromBatchEntries(completionResponses, insightLogs)
	if err != nil {
		g.logger.Error("Failed to transform batch entries into application insights")
		return nil, err
	}

	return insights, nil
}

func (g *OpenAiInsightsGenerator) AwaitScheduledApplicationInsights(
	sheduledInsights *ScheduledApplicationInsights,
) ([]ApplicationInsightsWithMetadata, error) {
	jobs, failedBatches, err := g.batchPoller.AwaitPendingJobs(sheduledInsights.ScheduledJobIds)

	if err != nil {
		g.logger.Error("Failed to wait for pending application batches", zap.Error(err))
		return nil, err
	}

	// Ignoring failed batches
	if len(failedBatches) > 0 {
		g.logger.Error("Some of the application batches have failed", zap.Any("failedBatches", failedBatches))
	}

	batches, err := g.batchPoller.BatchesFromJobs(jobs)
	if err != nil {
		g.logger.Error("Failed to fetch node batches from scheduled jobs", zap.Error(err))
		return nil, err
	}

	completionResponses, err := g.client.CompletionResponseEntriesFromBatches(batches)
	if err != nil {
		g.logger.Error("Failed to get application completion responses from batches", zap.Error(err))
		return nil, err
	}

	insightLogs, err := g.applicationLogsRepository.
		GetLogs(context.TODO(), sheduledInsights.ClusterId,
			time.UnixMilli(sheduledInsights.SinceMs),
			time.UnixMilli(sheduledInsights.ToMs))

	if err != nil {
		g.logger.Error("Failed to get application logs for scheduled insight")
		return nil, err
	}

	insights, err := g.getApplicationInsightsFromBatchEntries(completionResponses, insightLogs)
	if err != nil {
		g.logger.Error("Failed to transform batch entries into application insights")
		return nil, err
	}

	return insights, nil
}

func (g *OpenAiInsightsGenerator) getApplicationInsightsFromBatchEntries(
	batchEntries []*openai.BatchFileCompletionResponseEntry,
	logs []*repositories.ApplicationLogsDocument) ([]ApplicationInsightsWithMetadata, error) {

	// Each jsonl entry contains insights for a single application
	res := []ApplicationInsightsWithMetadata{}
	for _, response := range batchEntries {
		var applicationInsights applicationInsightsResponseDto
		if len(response.Response.Body.Choices) == 0 {
			return nil, errors.New("Failed to get insights from batch completion choices")
		}

		messageContent := response.Response.Body.Choices[0].Message.Content
		err := json.Unmarshal([]byte(messageContent), &applicationInsights)
		if err != nil {

			//OpenAI returned incorrectly structured output (skipping to minimize impact)
			g.logger.Error("OpenAI returned incorrectly formatted application insights output", zap.Error(err), zap.Any("content", messageContent))
			continue
		}

		insightsWithMetadata := array.Map(func(insight ApplicationLogsInsight) ApplicationInsightsWithMetadata {
			return g.addMetadataToApplicationInsight(insight, logs)
		})(applicationInsights.Insights)

		// We ignore insights without valid source
		insightsWithMetadata = array.Filter(func(insight ApplicationInsightsWithMetadata) bool {
			return len(insight.Metadata) > 0
		})(insightsWithMetadata)

		res = append(res, insightsWithMetadata...)
	}

	return res, nil
}

// func (g *OpenAiInsightsGenerator)
// func FilterByApplication()

func (g *OpenAiInsightsGenerator) ScheduleApplicationInsights(
	logs []*repositories.ApplicationLogsDocument,
	configuration []*ApplicationInsightConfiguration,
	scheduledTime time.Time,
	clusterId string,
	sinceMs int64,
	toMs int64,
) (*ScheduledApplicationInsights, error) {

	groupedLogs := GroupApplicationLogsByName(logs)
	configurationsByApplication := MapApplicationNameToConfiguration(configuration)
	completionRequests := make([]*openai.CompletionRequest, 0, len(groupedLogs))

	// In place filtering based on configured accuracy
	FilterByApplicationsAccuracy(groupedLogs, configurationsByApplication)

	// Generate insights for each application separately.
	for applicationName, logs := range groupedLogs {

		logPackets := repositories.SplitLogsIntoPackets(logs, g.client.ContextSizeBytes)
		// logPackets

		for _, logPacket := range logPackets {
			messages, err := g.createMessagesFromApplicationLogs(
				logPacket,
				configurationsByApplication[applicationName],
			)
			if err != nil {
				g.logger.Error("Failed to create messages from application logs", zap.Error(err), zap.String("app", applicationName))
				return nil, err
			}

			completionRequests = append(completionRequests,
				&openai.CompletionRequest{
					Messages:       messages,
					Temperature:    g.client.Temperature,
					ResponseFormat: openai.CreateJsonReponseFormat("insights", applicationInsightsResponseDto{}),
					Model:          g.client.Model(),
				},
			)
		}
	}

	completionRequestsPerJob, err := g.client.SplitCompletionReqestsByBatchSize(completionRequests)
	if err != nil {
		g.logger.Error("Failed to split application insights completion requests", zap.Error(err))
		return nil, err
	}

	scheduledJobs := array.Map(openai.NewOpenAiJob)(completionRequestsPerJob)

	ids, repositoryErr := g.scheduledJobsRepository.InsertScheduledJobs(context.Background(), scheduledJobs)
	if repositoryErr != nil {
		g.logger.Error("Failed to split application insights completion requests", zap.Error(err))
		return nil, err
	}

	return &ScheduledApplicationInsights{
		ScheduledJobIds:          ids,
		ClusterId:                clusterId,
		SinceMs:                  sinceMs,
		ToMs:                     toMs,
		ApplicationConfiguration: configuration,
	}, nil
}

func (g *OpenAiInsightsGenerator) createMessagesFromApplicationLogs(
	logs []*repositories.ApplicationLogsDocument,
	configuration *ApplicationInsightConfiguration) ([]*openai.Message, error) {

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
			same events on different containers/pods. For each incident assign urgency as an integer number between 0 and 2.
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
	configuration *ApplicationInsightConfiguration) ([]ApplicationLogsInsight, error) {

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

	if len(openAiResponse.Choices) == 0 {
		g.logger.Error("No insight choices were returned for", zap.Any("application", configuration.ApplicationName))
		return nil, err
	}

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
