package insights

import (
	"github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/openai"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type OpenAiInsightsGenerator struct {
	client                    *openai.Client
	logger                    *zap.Logger
	applicationLogsRepository repositories.ApplicationLogsRepository
	nodeLogsRepository        repositories.NodeLogsRepository
}

type OpenAiInsightsGeneratorParams struct {
	fx.In
	Client                    *openai.Client
	Logger                    *zap.Logger
	ApplicationLogsRepository repositories.ApplicationLogsRepository
	NodeLogsRepository        repositories.NodeLogsRepository
}

func NewOpenAiInsightsGenerator(params OpenAiInsightsGeneratorParams) *OpenAiInsightsGenerator {
	return &OpenAiInsightsGenerator{
		client:                    params.Client,
		logger:                    params.Logger,
		applicationLogsRepository: params.ApplicationLogsRepository,
		nodeLogsRepository:        params.NodeLogsRepository,
	}
}

func GroupApplicationLogsByName(logs []*repositories.ApplicationLogsDocument) map[string][]*repositories.ApplicationLogsDocument {
	groupedLogs := make(map[string][]*repositories.ApplicationLogsDocument)
	for _, log := range logs {
		groupedLogs[log.ApplicationName] = append(groupedLogs[log.ApplicationName], log)
	}

	return groupedLogs
}

func GroupNodeLogsByName(logs []*repositories.NodeLogsDocument) map[string][]*repositories.NodeLogsDocument {
	groupedLogs := make(map[string][]*repositories.NodeLogsDocument)
	for _, log := range logs {
		groupedLogs[log.Name] = append(groupedLogs[log.Name], log)
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

func GroupInsightsByNode(nodeInsights []NodeInsightsWithMetadata) map[string][]NodeInsightsWithMetadata {
	insightsByNode := make(map[string][]NodeInsightsWithMetadata)

	for _, insight := range nodeInsights {
		nodeName := insight.Metadata[0].NodeName
		insightsByNode[nodeName] = append(insightsByNode[nodeName], insight)
	}
	return insightsByNode
}
