package insights

import (
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
