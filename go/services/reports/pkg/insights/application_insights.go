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

type ApplicationInsightConfiguration struct {
	ApplicationName string   `json:"applicationName"`
	Accuracy        Accuracy `json:"accuracy"`
	CustomPrompt    string   `json:"customPrompt"`
}

type ScheduledApplicationInsights struct {
	ScheduledJobIds          []string                           `bson:"scheduledJobIds" json:"scheduledJobIds"`
	SinceMs                  int64                              `bson:"sinceMs" json:"sinceMs"`
	ToMs                     int64                              `bson:"toMs" json:"toMs"`
	ClusterId                string                             `bson:"clusterId" json:"clusterId"`
	ApplicationConfiguration []*ApplicationInsightConfiguration `bson:"applicationConfiguration" json:"applicationConfiguration"`
}

type ApplicationLogsInsight struct {
	ApplicationName string   `json:"applicationName"`
	Title           string   `json:"title"`
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
	ScheduleApplicationInsights(
		logsByApplication map[string][]*repositories.ApplicationLogsDocument,
		configurationByApplication map[string]*ApplicationInsightConfiguration,
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

type ApplicationInsightsResponseDto struct {
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

func GroupApplicationConfigurationByName(applicationConfigurations []*ApplicationInsightConfiguration) map[string]*ApplicationInsightConfiguration {
	applicationConfigurationsByName := make(map[string]*ApplicationInsightConfiguration)

	for _, configuration := range applicationConfigurations {
		applicationName := configuration.ApplicationName
		applicationConfigurationsByName[applicationName] = configuration
	}

	return applicationConfigurationsByName
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
	scheduledInsight *ScheduledApplicationInsights) (*ApplicationInsightsWithMetadata, error) {

	applicationInsightsMetadata := make([]ApplicationInsightMetadata, 0, len(insight.SourceLogIds))

	logs, err := g.applicationLogsRepository.GetLogsByIds(context.Background(),
		scheduledInsight.ClusterId,
		time.UnixMilli(scheduledInsight.SinceMs),
		time.UnixMilli(scheduledInsight.ToMs),
		insight.SourceLogIds)
	if err != nil {
		g.logger.Error("Failed to get logs by ids to insight metadata")
		return nil, err
	}
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

	return &ApplicationInsightsWithMetadata{
		Insight:  &insight,
		Metadata: applicationInsightsMetadata,
	}, nil
}

// Get Application insights grouped by application name
func (g *OpenAiInsightsGenerator) GetScheduledApplicationInsights(
	scheduledInsights *ScheduledApplicationInsights,
) ([]ApplicationInsightsWithMetadata, error) {
	batches, err := g.client.Batches(scheduledInsights.ScheduledJobIds)

	if err != nil {
		g.logger.Error("Failed to get batch from id", zap.Error(err))
		return nil, err
	}

	completionResponses, err := g.client.CompletionResponseEntriesFromBatches(batches)
	if err != nil {
		g.logger.Error("Failed to get application completion responses from batches", zap.Error(err))
		return nil, err
	}

	if err != nil {
		g.logger.Error("Failed to get application logs for scheduled insight")
		return nil, err
	}

	insights, err := g.GetApplicationInsightsFromBatchEntries(completionResponses, scheduledInsights)
	if err != nil {
		g.logger.Error("Failed to transform batch entries into application insights")
		return nil, err
	}

	return insights, nil
}

func (g *OpenAiInsightsGenerator) AwaitScheduledApplicationInsights(
	scheduledInsights *ScheduledApplicationInsights,
) ([]ApplicationInsightsWithMetadata, error) {
	jobs, failedBatches, err := g.batchPoller.AwaitPendingJobs(scheduledInsights.ScheduledJobIds)

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

	if err != nil {
		g.logger.Error("Failed to get application logs for scheduled insight")
		return nil, err
	}

	insights, err := g.GetApplicationInsightsFromBatchEntries(completionResponses, scheduledInsights)
	if err != nil {
		g.logger.Error("Failed to transform batch entries into application insights")
		return nil, err
	}

	return insights, nil
}

func (g *OpenAiInsightsGenerator) GetApplicationInsightsFromBatchEntries(
	batchEntries []*openai.BatchFileCompletionResponseEntry,
	scheduledInsights *ScheduledApplicationInsights) ([]ApplicationInsightsWithMetadata, error) {

	// Each jsonl entry contains insights for a single application
	res := []ApplicationInsightsWithMetadata{}
	for _, response := range batchEntries {
		var applicationInsights ApplicationInsightsResponseDto
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

		insightsWithMetadata := make([]ApplicationInsightsWithMetadata, 0, 0)
		for _, insight := range applicationInsights.Insights {

			insightWithMetadata, err := g.addMetadataToApplicationInsight(insight, scheduledInsights)
			if err != nil {
				g.logger.Error("Failed to add metadata to node insight, skipping")
				continue
			}
			insightsWithMetadata = append(insightsWithMetadata, *insightWithMetadata)
		}

		// We ignore insights without valid source
		insightsWithMetadata = array.Filter(func(insight ApplicationInsightsWithMetadata) bool {
			return len(insight.Metadata) > 0
		})(insightsWithMetadata)

		res = append(res, insightsWithMetadata...)
	}

	return res, nil
}

func (g *OpenAiInsightsGenerator) ScheduleApplicationInsights(
	groupedLogs map[string][]*repositories.ApplicationLogsDocument,
	configurationByApplication map[string]*ApplicationInsightConfiguration,
	clusterId string,
	sinceMs int64,
	toMs int64,
) (*ScheduledApplicationInsights, error) {

	completionRequests := make(map[string]*openai.CompletionRequest, len(groupedLogs))

	// Generate insights for each application separately.
	for applicationName, logs := range groupedLogs {

		logPackets := repositories.SplitLogsIntoPackets(logs, g.client.ContextSizeBytes)

		for idx, logPacket := range logPackets {
			messages, err := g.createMessagesFromApplicationLogs(
				logPacket,
				configurationByApplication[applicationName],
			)
			if err != nil {
				g.logger.Error("Failed to create messages from application logs", zap.Error(err), zap.String("app", applicationName))
				return nil, err
			}

			// Assign each request a customId with applicationName and a packetId
			completionRequests[fmt.Sprintf("%s-%d", applicationName, idx)] =
				&openai.CompletionRequest{
					Messages:       messages,
					Temperature:    g.client.Temperature,
					ResponseFormat: openai.CreateJsonReponseFormat("insights", ApplicationInsightsResponseDto{}),
					Model:          g.client.Model(),
				}
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

	applicationConfigurationsList := maps.Values(configurationByApplication)

	return &ScheduledApplicationInsights{
		ScheduledJobIds:          ids,
		ClusterId:                clusterId,
		SinceMs:                  sinceMs,
		ToMs:                     toMs,
		ApplicationConfiguration: applicationConfigurationsList,
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
			Content: fmt.Sprintf(`You are a Kubernetes cluster system administrator. 
        Analyze the provided Kubernetes cluster logs and identify only meaningful insights that indicate potential errors or issues. 
        For each identified issue, provide:
        - A title as a brief summary (few words).
        - A detailed summary (max 50 words) explaining the issue, likely cause, and category of the issue.
        - An urgency level (LOW, MEDIUM, HIGH).
        - A recommended action to resolve the issue.

        Each insight should:
        - Include the relevant unmodified source log entries in the sourceIds array, referencing logs across containers/pods if they indicate the same issue. 
        Use _id field of provided log.
        - Avoid redundant insights; if a source is already used in an insight, do not generate another one for it.
        - Ignore logs that do not explicitly indicate issues or represent typical operational activities.
        - Exclude insights entirely if no errors or warnings are present.

        Here is the additional configuration to consider when generating insights: %s`, customPrompt),
		},
		{
			Role: "user",
			Content: fmt.Sprintf(`These are logs from an application named: %s.
        Please identify and describe any issues or anomalies indicated by the logs:
        %s`, configuration.ApplicationName, encodedLogs),
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
		openai.CreateJsonReponseFormat("application_insights", ApplicationInsightsResponseDto{}),
	)

	if err != nil {
		g.logger.Error("Failed to get application logs insights from openai client", zap.Error(err))
		return nil, err
	}

	var insights ApplicationInsightsResponseDto

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

var _ ApplicationInsightsGenerator = &OpenAiInsightsGenerator{}
