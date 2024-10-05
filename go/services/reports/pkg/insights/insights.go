package insights

import (
	"github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	batchcache "github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/batch_cache"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/openai"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type OpenAiInsightsGenerator struct {
	client     *openai.Client
	logger     *zap.Logger
	batchCache batchcache.BatchCache
}

type OpenAiInsightsGeneratorParams struct {
	fx.In
	Client     *openai.Client
	Logger     *zap.Logger
	BatchCache batchcache.BatchCache
}

func NewOpenAiInsightsGenerator(params OpenAiInsightsGeneratorParams) *OpenAiInsightsGenerator {
	return &OpenAiInsightsGenerator{
		client:     params.Client,
		logger:     params.Logger,
		batchCache: params.BatchCache,
	}
}

func GroupLogsByName(logs []*repositories.ApplicationLogsDocument) map[string][]*repositories.ApplicationLogsDocument {
	groupedLogs := make(map[string][]*repositories.ApplicationLogsDocument)
	for _, log := range logs {
		groupedLogs[log.ApplicationName] = append(groupedLogs[log.ApplicationName], log)
	}

	return groupedLogs
}

func GroupInsightsByApplication(applicationInsights []ApplicationInsightsWithMetadata) map[string][]ApplicationInsightsWithMetadata {
	insightsByApplication := make(map[string][]ApplicationInsightsWithMetadata)

	for _, insight := range applicationInsights {
		applicationName := insight.Insight.ApplicationName
		insightsByApplication[applicationName] = append(insightsByApplication[applicationName], insight)
	}
	return insightsByApplication
}
