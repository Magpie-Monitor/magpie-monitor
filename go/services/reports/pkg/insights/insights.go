package insights

import (
	"github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/openai"
	scheduledjobs "github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/scheduled_jobs"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Urgency string

const (
	Urgency_Low    Urgency = "LOW"
	Urgency_Medium Urgency = "MEDIUM"
	Urgency_High   Urgency = "HIGH"
)

func (u Urgency) Level() int {
	switch u {
	case Urgency_Low:
		return 0
	case Urgency_Medium:
		return 1
	case Urgency_High:
		return 2
	}

	return 0
}

func MaxUrgency(urgencies []Urgency) Urgency {
	maxUrgency := Urgency_Low

	for _, urgency := range urgencies {
		if urgency.Level() > maxUrgency.Level() {
			maxUrgency = urgency
		}
	}

	return maxUrgency
}

type OpenAiInsightsGenerator struct {
	client                    *openai.Client
	logger                    *zap.Logger
	applicationLogsRepository repositories.ApplicationLogsRepository
	nodeLogsRepository        repositories.NodeLogsRepository
	batchPoller               *openai.BatchPoller
	scheduledJobsRepository   scheduledjobs.ScheduledJobRepository[*openai.OpenAiJob]
}

type OpenAiInsightsGeneratorParams struct {
	fx.In
	Client                    *openai.Client
	Logger                    *zap.Logger
	ApplicationLogsRepository repositories.ApplicationLogsRepository
	NodeLogsRepository        repositories.NodeLogsRepository
	BatchPoller               *openai.BatchPoller
	ScheduledJobsRepository   scheduledjobs.ScheduledJobRepository[*openai.OpenAiJob]
}

func NewOpenAiInsightsGenerator(params OpenAiInsightsGeneratorParams) *OpenAiInsightsGenerator {
	return &OpenAiInsightsGenerator{
		client:                    params.Client,
		logger:                    params.Logger,
		applicationLogsRepository: params.ApplicationLogsRepository,
		nodeLogsRepository:        params.NodeLogsRepository,
		batchPoller:               params.BatchPoller,
		scheduledJobsRepository:   params.ScheduledJobsRepository,
	}
}

type Accuracy = string

const (
	Accuracy__High   Accuracy = "HIGH"
	Accuracy__Medium Accuracy = "MEDIUM"
	Accuracy__Low    Accuracy = "LOW"
)
