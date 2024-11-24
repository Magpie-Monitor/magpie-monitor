package incidentcorrelation

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/IBM/fp-go/array"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/splitting"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/insights"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/openai"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/repositories"
	scheduledjobs "github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/scheduled_jobs"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type IncidentMergeCriteria struct {
	Id             string           `json:"id"`
	Title          string           `json:"title"`
	Summary        string           `json:"summary"`
	Recommendation string           `json:"recommendation"`
	Category       string           `json:"category"`
	Urgency        insights.Urgency `json:"urgency"`
}

type IncidentMerger interface {
	ScheduleIncidentsMerge(incidents map[string][]repositories.Incident) ([]*repositories.ScheduledIncidentMergerJob, error)
	IsIncidentMergerJobFinished(job *repositories.ScheduledIncidentMergerJob) (bool, error)
	AreAllJobsFinished(jobs []*repositories.ScheduledIncidentMergerJob) (bool, error)
	TryGettingIncidentMergerJobIfFinished(job *repositories.ScheduledIncidentMergerJob) (map[string][]IncidentMergeGroup, error)
}

type OpenAiIncidentMerger struct {
	client                  *openai.Client
	logger                  *zap.Logger
	batchPoller             *openai.BatchPoller
	scheduledJobsRepository scheduledjobs.ScheduledJobRepository[*openai.OpenAiJob]
}

type OpenAiIncidentMergerParams struct {
	fx.In
	Client                  *openai.Client
	Logger                  *zap.Logger
	BatchPoller             *openai.BatchPoller
	ScheduledJobsRepository scheduledjobs.ScheduledJobRepository[*openai.OpenAiJob]
}

func NewOpenAiIncidentMerger(params OpenAiIncidentMergerParams) *OpenAiIncidentMerger {
	return &OpenAiIncidentMerger{
		client:                  params.Client,
		logger:                  params.Logger,
		batchPoller:             params.BatchPoller,
		scheduledJobsRepository: params.ScheduledJobsRepository,
	}
}

func (m *OpenAiIncidentMerger) getIncidentMergeCriteria(incident repositories.Incident) *IncidentMergeCriteria {
	return &IncidentMergeCriteria{
		Id:             incident.GetId(),
		Title:          incident.GetTitle(),
		Summary:        incident.GetSummary(),
		Recommendation: incident.GetRecommendation(),
		Urgency:        incident.GetUrgency(),
		Category:       incident.GetCategory(),
	}
}

func (m *OpenAiIncidentMerger) encodeMergeCriteria(incidents []*IncidentMergeCriteria) ([]string, error) {

	var encodedSummaries = make([]string, 0, len(incidents))

	for _, summary := range incidents {
		encodedSummary, err := json.Marshal(summary)
		if err != nil {
			m.logger.Error("Failed to encode incident summary", zap.Error(err), zap.Any("incident", summary))
			return nil, err
		}

		encodedSummaries = append(encodedSummaries, string(encodedSummary))
	}

	return encodedSummaries, nil
}

func (m *OpenAiIncidentMerger) createSummaryRequestMessage(content []string) []*openai.Message {

	return []*openai.Message{
		{
			Role:    "system",
			Content: "You are a Kubernetes cluster system administrator.",
		},
		{
			Role: "user",
			Content: fmt.Sprintf(`
    Given a list of incident reports, each with a summary, title, and recommendation for resolution, identify and group incidents 
    likely related to the same underlying event. For each incident:
    Title and Summary Similarity: Analyze the titles and summaries to identify incidents describing similar issues. Incidents with overlapping language, 
    key terms, or repeated phrases should be flagged as potentially related.
    Recommendation Consistency: Compare recommended resolutions. Incidents with highly similar or identical recommendations are 
    likely related, as they suggest a common resolution approach.
    Semantic Understanding: Use natural language processing to understand the contextual and semantic meaning of each summary and title. 
    Incidents that describe similar issues and solutions in different wording should still be grouped together based on meaning.
    Event Grouping: For each group of related incidents, assign a unique event identifier and 
    list all included incidents with their title, summary, and recommendation, noting the factors contributing to their grouping.
        %s`, content),
		},
	}
}

type IncidentMergeGroup struct {
	Title          string           `json:"title"`
	Summary        string           `json:"summary"`
	Recommendation string           `json:"recommendation"`
	Category       string           `json:"category"`
	Urgency        insights.Urgency `json:"urgency"`
	IncidentIds    []string         `json:"incidentIds"`
}

type incidentMergerResponseDto struct {
	IncidentMergeGroups []IncidentMergeGroup `json:"incidentMergeGroup"`
}

// Accepts a map of incidents groups, where a key is a group identifier
func (m *OpenAiIncidentMerger) ScheduleIncidentsMerge(incidentGroups map[string][]repositories.Incident) ([]*repositories.ScheduledIncidentMergerJob, error) {

	completionRequests := make(map[string]*openai.CompletionRequest, len(incidentGroups))
	if len(incidentGroups) == 0 {
		return make([]*repositories.ScheduledIncidentMergerJob, 0), nil
	}

	for groupId, report := range incidentGroups {
		criterias := array.Map(m.getIncidentMergeCriteria)(report)
		encodedCriterias, err := m.encodeMergeCriteria(criterias)

		encodedCriteriasPerPacket := splitting.SplitStringsIntoPackets(encodedCriterias, m.client.ContextSizeBytes)

		if err != nil {
			return nil, err
		}
		for idx, packet := range encodedCriteriasPerPacket {
			messages := m.createSummaryRequestMessage(packet)

			completionRequests[getGroupPacketId(groupId, idx)] = &openai.CompletionRequest{
				Messages:       messages,
				Temperature:    m.client.Temperature,
				ResponseFormat: openai.CreateJsonReponseFormat("incidentMerger", incidentMergerResponseDto{}),
				Model:          m.client.Model(),
			}
		}
	}

	completionReuqestsPerBatch, err := m.client.SplitCompletionReqestsByBatchSize(completionRequests)
	if err != nil {
		m.logger.Error("Failed to split merger requests by batch", zap.Error(err))
		return nil, err
	}

	jobs := make([]*repositories.ScheduledIncidentMergerJob, 0, len(completionReuqestsPerBatch))
	for _, batch := range completionReuqestsPerBatch {

		jobId, err := m.scheduledJobsRepository.InsertScheduledJob(context.Background(), &openai.OpenAiJob{
			ScheduledAt:        time.Now().UnixMilli(),
			CompletionRequests: batch,
			Status:             openai.OpenAiJobStatus__Enqueued,
		})

		if err != nil {
			m.logger.Error("Failed to insert scheduled jobs", zap.Error(err))
			return nil, err
		}

		jobs = append(jobs, &repositories.ScheduledIncidentMergerJob{
			Id: jobId,
		})

	}

	return jobs, nil
}

func getGroupPacketId(groupId string, packetId int) string {
	return fmt.Sprintf("%s-%d", groupId, packetId)
}

func getGroupFromEncodedGroupPacketId(groupPacketId string) string {
	splittedPacketId := strings.Split(groupPacketId, "-")
	return strings.Join(splittedPacketId[:len(splittedPacketId)-1], "-")
}

func (m *OpenAiIncidentMerger) TryGettingIncidentMergerJobIfFinished(job *repositories.ScheduledIncidentMergerJob) (map[string][]IncidentMergeGroup, error) {

	openAiJob, err := m.batchPoller.TryGettingJobIfFinished(job.Id)
	if err != nil {
		return nil, err
	}

	if openAiJob.Status == openai.OpenAiJobStatus__Failed {
		return nil, errors.New(fmt.Sprintf("Job %s has failed", job.Id))
	}

	if openAiJob.Status != openai.OpenAiJobStatus__Completed || *openAiJob.BatchId == "" {
		return nil, errors.New(fmt.Sprintf("Job %s is in an invalid state", job.Id))
	}

	completionEntries, err := m.client.CompletionResponseEntriesFromBatchById(*openAiJob.BatchId)
	if err != nil {
		m.logger.Error("Failed to get completion entries", zap.Error(err))
		return nil, err
	}

	incidentMergeGroups := make(map[string][]IncidentMergeGroup, 0)

	for _, completionEntry := range completionEntries {
		var responseDto incidentMergerResponseDto
		message := completionEntry.Response.Body.Choices[0].Message.Content
		err := json.Unmarshal([]byte(message), &responseDto)
		if err != nil {
			m.logger.Error("Failed to decode incidentMergeGroup from message", zap.Error(err))
			return nil, err
		}

		group := getGroupFromEncodedGroupPacketId(completionEntry.CustomId)

		incidentMergeGroups[group] = append(incidentMergeGroups[group], responseDto.IncidentMergeGroups...)
	}

	return incidentMergeGroups, nil
}

func (m *OpenAiIncidentMerger) IsIncidentMergerJobFinished(job *repositories.ScheduledIncidentMergerJob) (bool, error) {

	isFinished, err := m.batchPoller.IsJobFinished(job.Id)
	if err != nil {
		return false, err
	}

	return isFinished, nil
}

func (m *OpenAiIncidentMerger) AreAllJobsFinished(jobs []*repositories.ScheduledIncidentMergerJob) (bool, error) {
	for _, job := range jobs {
		if finished, err := m.IsIncidentMergerJobFinished(job); finished != true {
			return false, err
		}
	}

	return true, nil
}

var _ IncidentMerger = &OpenAiIncidentMerger{}
