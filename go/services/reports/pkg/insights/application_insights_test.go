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

func TestOpenAiInsightGenerator_GetApplicationInsightsFromBatchEntries(t *testing.T) {

	type TestDependencies struct {
		fx.In
		Logger                    *zap.Logger
		OpenAiInsightsGenerator   *insights.OpenAiInsightsGenerator
		ApplicationLogsRepository repositories.ApplicationLogsRepository
	}

	test := func(dependencies TestDependencies) {

		insertedLogIds := repositories.PrefillApplicationLogs(
			t,
			dependencies.Logger,
			dependencies.ApplicationLogsRepository,
			repositories.GetDefaultApplicationTestLogsFromCluster(repositories.DEFAULT_TEST_CLUSTER_ID),
		)

		encodedApp1Insights, err := json.Marshal(insights.ApplicationInsightsResponseDto{
			Insights: []insights.ApplicationLogsInsight{
				{
					ApplicationName: "test-app",
					Title:           "test-title",
					Category:        "test-category",
					Summary:         "test-summary",
					Recommendation:  "test-recommendation",
					Urgency:         insights.Urgency_Medium,
					SourceLogIds:    insertedLogIds,
				},
			},
		})
		assert.NoError(t, err, "Failed to prefill application logs")

		testCases := []struct {
			description       string
			batchEntries      []*openai.BatchFileCompletionResponseEntry
			scheduledInsights *insights.ScheduledApplicationInsights
			expectedInsights  []insights.ApplicationInsightsWithMetadata
		}{
			{
				batchEntries: []*openai.BatchFileCompletionResponseEntry{
					{
						CustomId: "app-1",
						Response: struct {
							StatusCode int                        `json:"status_code"`
							RequestId  string                     `json:"request_id"`
							Body       *openai.CompletionResponse `json:"body"`
						}{
							Body: &openai.CompletionResponse{
								Choices: []openai.Choice{
									{
										Message: &openai.Message{
											Content: string(encodedApp1Insights),
										},
									},
								},
							},
						},
					},
				},
				scheduledInsights: &insights.ScheduledApplicationInsights{
					SinceMs:                  int64(repositories.DEFAULT_TEST_APPLICATION_LOGS_COLLECTED_AT_DATE) - 1,
					ToMs:                     int64(repositories.DEFAULT_TEST_APPLICATION_LOGS_COLLECTED_AT_DATE) + 1,
					ClusterId:                repositories.DEFAULT_TEST_CLUSTER_ID,
					ApplicationConfiguration: []*insights.ApplicationInsightConfiguration{},
				},
				expectedInsights: []insights.ApplicationInsightsWithMetadata{
					{
						Insight: &insights.ApplicationLogsInsight{
							ApplicationName: "test-app",
							Title:           "test-title",
							Category:        "test-category",
							Summary:         "test-summary",
							Recommendation:  "test-recommendation",
							Urgency:         insights.Urgency_Medium,
							SourceLogIds:    insertedLogIds,
						},
						Metadata: []insights.ApplicationInsightMetadata{
							{
								CollectedAtMs:   repositories.DEFAULT_NODE_LOG.CollectedAtMs,
								ApplicationName: repositories.DEFAULT_TEST_APPLICATION_NAME,
								ClusterId:       repositories.DEFAULT_TEST_CLUSTER_ID,
								PodName:         repositories.DEFAULT_TEST_APPLICATION_POD_NAME,
								Image:           repositories.DEFAULT_TEST_APPLICATION_IMAGE_NAME,
								Source:          repositories.DEFAULT_TEST_APPLICATION_LOGS_CONTENT,
								ContainerName:   repositories.DEFAULT_TEST_APPLICATION_CONTAINER_NAME,
							},
						},
					},
				},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.description, func(t *testing.T) {
				insights, err := dependencies.OpenAiInsightsGenerator.GetApplicationInsightsFromBatchEntries(tc.batchEntries, tc.scheduledInsights)
				assert.NoError(t, err, "Failed to get insights")
				assert.Equal(t, tc.expectedInsights, insights)
			})
		}
	}

	tests.RunTest(test, t, config.AppModule)
}
