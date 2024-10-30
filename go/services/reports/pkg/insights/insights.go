package insights

import (
	"github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/openai"
	scheduledjobs "github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/scheduled_jobs"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Urgency int

const (
	_ Urgency = iota
	Urgency_Low
	Urgency_Medium
	Urgency_High
)

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
