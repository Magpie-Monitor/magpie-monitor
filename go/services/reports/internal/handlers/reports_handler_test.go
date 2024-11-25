package handlers_test

import (
	"context"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/envs"
	messagebroker "github.com/Magpie-Monitor/magpie-monitor/pkg/message-broker"
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
	"testing"
	"time"
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

func TestReportHandler_ListenForReportRequests(t *testing.T) {

	type TestDependencies struct {
		fx.In
		Logger                    *zap.Logger
		ReportHandler             *handlers.ReportsHandler
		ApplicationLogsRepository sharedrepositories.ApplicationLogsRepository
		NodeLogsRepository        sharedrepositories.NodeLogsRepository
		ReportRequestedBroker     messagebroker.MessageBroker[brokers.ReportRequested]
		ReportFailedBroker        messagebroker.MessageBroker[brokers.ReportRequestFailed]
		ReportsRepository         repositories.ReportRepository
	}

	var (
		testSinceMs        int64 = int64(sharedrepositories.DEFAULT_TEST_APPLICATION_LOGS_COLLECTED_AT_DATE) - 1
		testEndMs          int64 = int64(sharedrepositories.DEFAULT_TEST_APPLICATION_LOGS_COLLECTED_AT_DATE) + 1
		testClusterName          = "test-report-handler-cluster"
		testCorreleationId       = "23123-43242"
	)

	testCases := []struct {
		description string
		request     *brokers.ReportRequest
		expectedErr *brokers.ReportRequestFailed
	}{
		{
			description: "Properly generete scheduled report",
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
		},
		{
			description: "Fail to generate report, becase of missing clusterId parameter",
			request: &brokers.ReportRequest{
				ClusterId:                nil,
				SinceMs:                  &testSinceMs,
				ToMs:                     &testEndMs,
				ApplicationConfiguration: []*insights.ApplicationInsightConfiguration{},
				NodeConfiguration:        []*insights.NodeInsightConfiguration{},
			},
			expectedErr: brokers.NewReportRequestFailedValidation(
				testCorreleationId, "Missing clusterId parameter"),
		},
		{
			description: "Fail to generate report, becase of missing sinceMs parameter",
			request: &brokers.ReportRequest{
				ClusterId:                &testClusterName,
				SinceMs:                  nil,
				ToMs:                     &testEndMs,
				ApplicationConfiguration: []*insights.ApplicationInsightConfiguration{},
				NodeConfiguration:        []*insights.NodeInsightConfiguration{},
			},
			expectedErr: brokers.NewReportRequestFailedValidation(
				testCorreleationId, "Missing sinceMs parameter"),
		},
		{
			description: "Fail to generate report, becase of missing toMs parameter",
			request: &brokers.ReportRequest{
				ClusterId:                &testClusterName,
				SinceMs:                  &testSinceMs,
				ToMs:                     nil,
				ApplicationConfiguration: []*insights.ApplicationInsightConfiguration{},
				NodeConfiguration:        []*insights.NodeInsightConfiguration{},
			},
			expectedErr: brokers.NewReportRequestFailedValidation(
				testCorreleationId, "Missing toMs parameter"),
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

		for _, tc := range testCases {
			t.Run(tc.description, func(t *testing.T) {

				err := dependencies.ReportsRepository.DeleteAll(context.TODO())
				assert.NoError(t, err)

				dependencies.ReportRequestedBroker.Publish("test-message", brokers.ReportRequested{
					CorrelationId: testCorreleationId,
					ReportRequest: *tc.request,
				})

				go dependencies.ReportHandler.ListenForReportRequests()

				if tc.expectedErr != nil {
					messages := make(chan brokers.ReportRequestFailed)
					go dependencies.ReportFailedBroker.Subscribe(messages, make(chan error))

					ctx, _ := context.WithTimeout(context.Background(), time.Second*20*time.Duration(integrationTestWaitModifier))
					done := ctx.Done()

					select {
					case message := <-messages:
						message.TimestampMs = tc.expectedErr.TimestampMs
						assert.Equal(t, *tc.expectedErr, message)
						break
					case <-done:
						break
					}

				} else {
					time.Sleep(time.Duration(integrationTestWaitModifier) * 20 * time.Second)

					reports, err := dependencies.ReportsRepository.GetPendingGenerationReports(context.Background())

					assert.NoError(t, err)
					assert.Equal(t, 1, len(reports))
					assert.Equal(t, testCorreleationId, reports[0].CorrelationId)
				}

			})
		}
	}

	tests.RunTest(test, t, config.AppModule)
}
