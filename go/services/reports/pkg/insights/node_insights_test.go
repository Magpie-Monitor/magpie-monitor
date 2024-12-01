package insights_test

import (
	"encoding/json"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/tests"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/config"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/insights"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/openai"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"testing"
)

func TestOpenAiInsightGenerator_GetNodeInsightsFromBatchEntries(t *testing.T) {

	type TestDependencies struct {
		fx.In
		Logger                  *zap.Logger
		OpenAiInsightsGenerator *insights.OpenAiInsightsGenerator
		NodeLogsRepository      repositories.NodeLogsRepository
	}

	test := func(dependencies TestDependencies) {

		insertedLogIds := repositories.PrefillNodeLogs(
			t,
			dependencies.Logger,
			dependencies.NodeLogsRepository,
			repositories.GetDefaultNodeTestLogsFromCluster(repositories.DEFAULT_TEST_CLUSTER_ID),
		)

		encodedNodeInsights, err := json.Marshal(insights.NodeInsightsResponseDto{
			Insights: []insights.NodeLogsInsight{
				{
					NodeName:       "test-node",
					Title:          "test-title",
					Category:       "test-category",
					Summary:        "test-summary",
					Recommendation: "test-recommendation",
					Urgency:        insights.Urgency_Medium,
					SourceLogIds:   insertedLogIds,
				},
			},
		})
		assert.NoError(t, err, "Failed to prefill application logs")

		testCases := []struct {
			description       string
			batchEntries      []*openai.BatchFileCompletionResponseEntry
			scheduledInsights *insights.ScheduledNodeInsights
			expectedInsights  []insights.NodeInsightsWithMetadata
		}{
			{
				batchEntries: []*openai.BatchFileCompletionResponseEntry{
					{
						CustomId: "node-1",
						Response: struct {
							StatusCode int                        `json:"status_code"`
							RequestId  string                     `json:"request_id"`
							Body       *openai.CompletionResponse `json:"body"`
						}{
							Body: &openai.CompletionResponse{
								Choices: []openai.Choice{
									{
										Message: &openai.Message{
											Content: string(encodedNodeInsights),
										},
									},
								},
							},
						},
					},
				},
				scheduledInsights: &insights.ScheduledNodeInsights{
					SinceMs:           int64(repositories.DEFAULT_TEST_NODE_LOGS_COLLECTED_AT_DATE) - 1,
					ToMs:              int64(repositories.DEFAULT_TEST_NODE_LOGS_COLLECTED_AT_DATE) + 1,
					ClusterId:         repositories.DEFAULT_TEST_CLUSTER_ID,
					NodeConfiguration: []*insights.NodeInsightConfiguration{},
				},
				expectedInsights: []insights.NodeInsightsWithMetadata{
					{
						Insight: &insights.NodeLogsInsight{
							NodeName:       repositories.DEFAULT_TEST_NODE_NAME,
							Title:          "test-title",
							Category:       "test-category",
							Summary:        "test-summary",
							Recommendation: "test-recommendation",
							Urgency:        insights.Urgency_Medium,
							SourceLogIds:   insertedLogIds,
						},
						Metadata: []insights.NodeInsightMetadata{
							{
								NodeName:      repositories.DEFAULT_TEST_NODE_NAME,
								Filename:      repositories.DEFAULT_TEST_NODE_LOGS_FILENAME,
								CollectedAtMs: repositories.DEFAULT_NODE_LOG.CollectedAtMs,
								ClusterId:     repositories.DEFAULT_TEST_CLUSTER_ID,
								Source:        repositories.DEFAULT_TEST_NODE_LOGS_CONTENT,
							},
						},
					},
				},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.description, func(t *testing.T) {
				insights, err := dependencies.OpenAiInsightsGenerator.GetNodeInsightsFromBatchEntries(tc.batchEntries, tc.scheduledInsights)
				assert.NoError(t, err, "Failed to get insights")
				assert.Equal(t, tc.expectedInsights, insights)
			})
		}
	}

	tests.RunTest(test, t, config.AppModule)
}
