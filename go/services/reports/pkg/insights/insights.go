package insights

import (
	"github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/openai"
	"go.uber.org/zap"
)

type OpenAiInsightsGenerator struct {
	client *openai.Client
	logger *zap.Logger
}

func NewOpenAiInsightsGenerator(client *openai.Client, logger *zap.Logger) *OpenAiInsightsGenerator {
	return &OpenAiInsightsGenerator{
		client: client,
		logger: logger,
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
