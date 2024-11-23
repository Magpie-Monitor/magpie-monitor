package handlers_test

import (
	"context"
	"testing"
	"time"

	"github.com/Magpie-Monitor/magpie-monitor/pkg/envs"
	sharedrepositories "github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/tests"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/internal/brokers"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/internal/handlers"
	config "github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/config"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/insights"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/repositories"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func TestReportHandler_ScheduleReport(t *testing.T) {

	type TestDependencies struct {
		fx.In
		Logger                    *zap.Logger
		ReportHandler             *handlers.ReportsHandler
		ApplicationLogsRepository sharedrepositories.ApplicationLogsRepository
		NodeLogsRepository        sharedrepositories.NodeLogsRepository
	}

	var (
		testSinceMs        int64 = int64(sharedrepositories.DEFAULT_TEST_APPLICATION_LOGS_COLLECTED_AT_DATE) - 1
		testEndMs          int64 = int64(sharedrepositories.DEFAULT_TEST_APPLICATION_LOGS_COLLECTED_AT_DATE) + 1
		testClusterName          = "test-report-handler-cluster"
		testCorreleationId       = "23123-43242"
	)

	testCases := []struct {
		request        *brokers.ReportRequest
		expectedErr    *brokers.ReportRequestFailed
		expectedReport *repositories.Report
	}{
		{
			request: &brokers.ReportRequest{
				ClusterId: &testClusterName,
				SinceMs:   &testSinceMs,
				ToMs:      &testEndMs,
				ApplicationConfiguration: []*insights.ApplicationInsightConfiguration{
					{
						ApplicationName: sharedrepositories.DEFAULT_TEST_APPLICATION_NAME,
						CustomPrompt:    "Test application custom prompt",
						Accuracy:        insights.Accuracy__Medium,
					},
				},
				NodeConfiguration: []*insights.NodeInsightConfiguration{
					{
						NodeName:     sharedrepositories.DEFAULT_TEST_NODE_NAME,
						Accuracy:     insights.Accuracy__High,
						CustomPrompt: "Test node custom prompt",
					},
				},
			},
			expectedErr: nil,
			expectedReport: &repositories.Report{
				Status:    repositories.ReportState_AwaitingGeneration,
				ClusterId: testClusterName,
				SinceMs:   testSinceMs,
				ToMs:      testEndMs,
				ScheduledNodeInsights: &insights.ScheduledNodeInsights{
					SinceMs:   testSinceMs,
					ToMs:      testEndMs,
					ClusterId: testClusterName,
					NodeConfiguration: []*insights.NodeInsightConfiguration{
						{
							NodeName:     sharedrepositories.DEFAULT_TEST_NODE_NAME,
							Accuracy:     insights.Accuracy__High,
							CustomPrompt: "Test node custom prompt",
						},
					},
				},
				ScheduledApplicationInsights: &insights.ScheduledApplicationInsights{
					SinceMs:   testSinceMs,
					ToMs:      testEndMs,
					ClusterId: testClusterName,
					ApplicationConfiguration: []*insights.ApplicationInsightConfiguration{
						{
							ApplicationName: sharedrepositories.DEFAULT_TEST_APPLICATION_NAME,
							CustomPrompt:    "Test application custom prompt",
							Accuracy:        insights.Accuracy__Medium,
						},
					},
				},
			},
		},
		{
			request: &brokers.ReportRequest{
				ClusterId:                nil,
				SinceMs:                  &testSinceMs,
				ToMs:                     &testEndMs,
				ApplicationConfiguration: []*insights.ApplicationInsightConfiguration{},
				NodeConfiguration:        []*insights.NodeInsightConfiguration{},
			},
			expectedErr: brokers.NewReportRequestFailedValidation(
				testCorreleationId, "Missing clusterId parameter"),
			expectedReport: &repositories.Report{},
		},
		{
			request: &brokers.ReportRequest{
				ClusterId:                &testClusterName,
				SinceMs:                  nil,
				ToMs:                     &testEndMs,
				ApplicationConfiguration: []*insights.ApplicationInsightConfiguration{},
				NodeConfiguration:        []*insights.NodeInsightConfiguration{},
			},
			expectedErr: brokers.NewReportRequestFailedValidation(
				testCorreleationId, "Missing sinceMs parameter"),
			expectedReport: &repositories.Report{},
		},
		{
			request: &brokers.ReportRequest{
				ClusterId:                &testClusterName,
				SinceMs:                  &testSinceMs,
				ToMs:                     nil,
				ApplicationConfiguration: []*insights.ApplicationInsightConfiguration{},
				NodeConfiguration:        []*insights.NodeInsightConfiguration{},
			},
			expectedErr: brokers.NewReportRequestFailedValidation(
				testCorreleationId, "Missing toMs parameter"),
			expectedReport: &repositories.Report{},
		},
	}

	integrationTestWaitModifier := envs.ConvertToInt(tests.INTEGRATION_TEST_WAIT_MODIFIER_KEY)

	test := func(dependencies TestDependencies) {

		if dependencies.ReportHandler == nil {
			t.Fatal("Failed to load report handler")
		}

		sharedrepositories.PrefillApplicationLogs(t,
			dependencies.Logger,
			dependencies.ApplicationLogsRepository,
			sharedrepositories.GetDefaultApplicationTestLogsFromCluster(testClusterName),
		)

		sharedrepositories.PrefillNodeLogs(t,
			dependencies.Logger,
			dependencies.NodeLogsRepository,
			sharedrepositories.GetDefaultNodeTestLogsFromCluster(testClusterName),
		)

		time.Sleep(time.Second * 10 * time.Duration(integrationTestWaitModifier))

		for _, tc := range testCases {
			actualReport, err := dependencies.ReportHandler.ScheduleReport(context.Background(), testCorreleationId, tc.request)
			if tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr.CorrelationId, err.CorrelationId)
				assert.Equal(t, tc.expectedErr.ErrorType, err.ErrorType)
				assert.Equal(t, tc.expectedErr.ErrorMessage, err.ErrorMessage)
				continue
			}

			assert.Equal(t, tc.expectedReport.Status, actualReport.Status)
			assert.Equal(t, testCorreleationId, actualReport.CorrelationId)
			assert.Equal(t, tc.expectedReport.ClusterId, actualReport.ClusterId)
			assert.Equal(t, tc.expectedReport.SinceMs, actualReport.SinceMs)
			assert.Equal(t, tc.expectedReport.ToMs, actualReport.ToMs)
			assert.Equal(t,
				tc.expectedReport.ScheduledApplicationInsights.ApplicationConfiguration,
				actualReport.ScheduledApplicationInsights.ApplicationConfiguration)
			assert.Equal(t,
				tc.expectedReport.ScheduledNodeInsights.NodeConfiguration,
				actualReport.ScheduledNodeInsights.NodeConfiguration)
		}

	}

	tests.RunTest(test, t, config.AppModule)
}
