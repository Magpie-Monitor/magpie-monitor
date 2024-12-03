package services_test

import (
	sharedrepositories "github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/tests"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/internal/services"
	config "github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/config"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/insights"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/repositories"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"testing"
)

func TestReportsService_getReportUrgencyFromApplicationAndNodeReports(t *testing.T) {

	type TestDependencies struct {
		fx.In
		Logger                    *zap.Logger
		ReportsService            *services.ReportsService
		ApplicationLogsRepository sharedrepositories.ApplicationLogsRepository
		NodeLogsRepository        sharedrepositories.NodeLogsRepository
	}

	testCases := []struct {
		applicationReports []*repositories.ApplicationReport
		nodeReports        []*repositories.NodeReport
		description        string
		expectedUrgency    insights.Urgency
	}{
		{
			description: "Get highest urgency of all incidents",
			applicationReports: []*repositories.ApplicationReport{
				{
					Incidents: []*repositories.ApplicationIncident{
						{
							Urgency: insights.Urgency_Low,
						},
						{
							Urgency: insights.Urgency_Low,
						},
						{
							Urgency: insights.Urgency_Medium,
						},
					},
				},
			},

			nodeReports: []*repositories.NodeReport{
				{
					Incidents: []*repositories.NodeIncident{
						{
							Urgency: insights.Urgency_Low,
						},
						{
							Urgency: insights.Urgency_High,
						},
						{
							Urgency: insights.Urgency_Medium,
						},
					},
				},
			},
			expectedUrgency: insights.Urgency_High,
		},
		{
			description: "Given an empty list of incidents return LOW urgency",
			applicationReports: []*repositories.ApplicationReport{
				{
					Incidents: []*repositories.ApplicationIncident{},
				},
			},

			nodeReports: []*repositories.NodeReport{
				{
					Incidents: []*repositories.NodeIncident{},
				},
			},
			expectedUrgency: insights.Urgency_Low,
		},
	}

	test := func(dependencies TestDependencies) {

		for _, tc := range testCases {
			actualUrgency := dependencies.ReportsService.GetReportUrgencyFromApplicationAndNodeReports(tc.applicationReports, tc.nodeReports)
			assert.Equal(t, tc.expectedUrgency, actualUrgency, tc.description)
		}

	}

	tests.RunTest(test, t, config.AppModule)
}

// func TestReportsService_GetNodeReportsFromInsights(t *testing.T) {
//
// 	type TestDependencies struct {
// 		fx.In
// 		Logger                    *zap.Logger
// 		ReportsService            *services.ReportsService
// 		ApplicationLogsRepository sharedrepositories.ApplicationLogsRepository
// 		NodeLogsRepository        sharedrepositories.NodeLogsRepository
// 	}
//
// 	testCases := []struct {
// 		nodeInsights       []insights.NodeInsightsWithMetadata
// 		nodeConfigurations []*insights.NodeInsightConfiguration
// 		description        string
// 		expectedReports    []*repositories.NodeReport
// 	}{
// 		{
// 			description: "Generate report by grouping incidents by nodes",
// 			nodeInsights: []insights.NodeInsightsWithMetadata{
// 				{
// 					Insight: &insights.NodeLogsInsight{
// 						NodeName:       "test-node",
// 						Title:          "test-title",
// 						Category:       "test-category",
// 						Summary:        "test-summary",
// 						Recommendation: "test-recommendation",
// 						Urgency:        insights.Urgency_Medium,
// 						SourceLogIds:   []string{},
// 					},
// 					Metadata: []insights.NodeInsightMetadata{
// 						{
//
// 							ClusterId:     "test-cluster",
// 							NodeName:      "test-node",
// 							CollectedAtMs: 0,
// 							Filename:      "test-filename",
// 							Source:        "Log content",
// 						},
// 					},
// 				},
// 				{
// 					Insight: &insights.NodeLogsInsight{
// 						NodeName:       "test-node",
// 						Title:          "test-title-2",
// 						Category:       "test-category-2",
// 						Summary:        "test-summary-2",
// 						Recommendation: "test-recommendation-2",
// 						Urgency:        insights.Urgency_Medium,
// 						SourceLogIds:   []string{},
// 					},
// 					Metadata: []insights.NodeInsightMetadata{
// 						{
//
// 							ClusterId:     "test-cluster",
// 							NodeName:      "test-node",
// 							CollectedAtMs: 2,
// 							Filename:      "test-filename-2",
// 							Source:        "Log content 2",
// 						},
// 					},
// 				},
// 			},
// 			nodeConfigurations: []*insights.NodeInsightConfiguration{
// 				{
// 					NodeName:     "test-node",
// 					CustomPrompt: "Custom prompt for test-node",
// 					Accuracy:     insights.Accuracy__High,
// 				},
// 				{
// 					NodeName:     "test-node-2",
// 					CustomPrompt: "Custom prompt for test-node-2",
// 					Accuracy:     insights.Accuracy__Medium,
// 				},
// 			},
// 			expectedReports: []*repositories.NodeReport{
// 				{
// 					Node:         "test-node",
// 					CustomPrompt: "Custom prompt for test-node",
// 					Accuracy:     insights.Accuracy__High,
// 					Incidents: []*repositories.NodeIncident{
// 						{
// 							NodeName:       "test-node",
// 							Title:          "test-title",
// 							Category:       "test-category",
// 							Summary:        "test-summary",
// 							Recommendation: "test-recommendation",
// 							Urgency:        insights.Urgency_Medium,
// 							ClusterId:      "test-cluster",
// 							CustomPrompt:   "Custom prompt for test-node",
// 							Accuracy:       insights.Accuracy__High,
// 							Sources: []repositories.NodeIncidentSource{
// 								{
// 									Timestamp: 0,
// 									Content:   "Log content",
// 									Filename:  "test-filename",
// 								},
// 							},
// 						},
// 						{
// 							NodeName:       "test-node",
// 							Title:          "test-title-2",
// 							Category:       "test-category-2",
// 							Summary:        "test-summary-2",
// 							Recommendation: "test-recommendation-2",
// 							Urgency:        insights.Urgency_Medium,
// 							ClusterId:      "test-cluster",
// 							CustomPrompt:   "Custom prompt for test-node",
// 							Accuracy:       insights.Accuracy__High,
// 							Sources: []repositories.NodeIncidentSource{
// 								{
// 									Timestamp: 2,
// 									Content:   "Log content 2",
// 									Filename:  "test-filename-2",
// 								},
// 							},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}
//
// 	test := func(dependencies TestDependencies) {
//
// 		for _, tc := range testCases {
// 			t.Run(tc.description, func(t *testing.T) {
// 				reports, err := dependencies.ReportsService.GetNodeReportsFromInsights(
// 					tc.nodeInsights, tc.nodeConfigurations,
// 				)
//
// 				assert.NoError(t, err)
// 				assert.Equal(t, len(tc.expectedReports), len(reports), tc.description)
// 				for idx, report := range reports {
// 					expectedReport := tc.expectedReports[idx]
// 					assert.Equal(t, expectedReport.Accuracy, report.Accuracy)
// 					assert.Equal(t, expectedReport.Node, report.Node)
// 					assert.Equal(t, expectedReport.CustomPrompt, report.CustomPrompt)
//
// 					assert.Equal(t, len(expectedReport.Incidents), len(report.Incidents))
// 					for incidentId, incident := range report.Incidents {
// 						expectedIncident := expectedReport.Incidents[incidentId]
// 						expectedIncident.Id = incident.Id
// 						assert.Equal(t, expectedIncident, incident)
// 					}
// 				}
// 			})
// 		}
// 	}
//
// 	tests.RunTest(test, t, config.AppModule)
// }
//
//
// func TestReportsService_GetApplicationReportsFromInsights(t *testing.T) {
//
// 	type TestDependencies struct {
// 		fx.In
// 		Logger                    *zap.Logger
// 		ReportsService            *services.ReportsService
// 		ApplicationLogsRepository sharedrepositories.ApplicationLogsRepository
// 		NodeLogsRepository        sharedrepositories.NodeLogsRepository
// 	}
//
// 	testCases := []struct {
// 		applicationInsights       []insights.ApplicationInsightsWithMetadata
// 		applicationConfigurations []*insights.ApplicationInsightConfiguration
// 		description               string
// 		expectedReports           []*repositories.ApplicationReport
// 	}{
// 		{
// 			description: "Generate report by grouping incidents by applications",
// 			applicationInsights: []insights.ApplicationInsightsWithMetadata{
// 				{
// 					Insight: &insights.ApplicationLogsInsight{
// 						ApplicationName: "test-app",
// 						Title:           "test-title",
// 						Category:        "test-category",
// 						Summary:         "test-summary",
// 						Recommendation:  "test-recommendation",
// 						Urgency:         insights.Urgency_Medium,
// 						SourceLogIds:    []string{},
// 					},
// 					Metadata: []insights.ApplicationInsightMetadata{
// 						{
//
// 							ClusterId:       "test-cluster",
// 							ApplicationName: "test-app",
// 							ContainerName:   "test-container",
// 							PodName:         "test-pod",
// 							Image:           "test-image",
// 							CollectedAtMs:   0,
// 							Source:          "Log content",
// 						},
// 					},
// 				},
// 				{
// 					Insight: &insights.ApplicationLogsInsight{
// 						ApplicationName: "test-app",
// 						Title:           "test-title-2",
// 						Category:        "test-category-2",
// 						Summary:         "test-summary-2",
// 						Recommendation:  "test-recommendation-2",
// 						Urgency:         insights.Urgency_Medium,
// 						SourceLogIds:    []string{},
// 					},
// 					Metadata: []insights.ApplicationInsightMetadata{
// 						{
// 							ApplicationName: "test-app",
// 							ClusterId:       "test-cluster",
// 							ContainerName:   "test-container",
// 							PodName:         "test-pod",
// 							Image:           "test-image",
// 							CollectedAtMs:   2,
// 							Source:          "Log content 2",
// 						},
// 					},
// 				},
// 			},
//
// 			applicationConfigurations: []*insights.ApplicationInsightConfiguration{
// 				{
// 					ApplicationName: "test-app",
// 					CustomPrompt:    "Custom prompt for test-app",
// 					Accuracy:        insights.Accuracy__High,
// 				},
// 				{
// 					ApplicationName: "test-app-2",
// 					CustomPrompt:    "Custom prompt for test-app-2",
// 					Accuracy:        insights.Accuracy__Medium,
// 				},
// 			},
// 			expectedReports: []*repositories.ApplicationReport{
// 				{
// 					ApplicationName: "test-app",
// 					CustomPrompt:    "Custom prompt for test-app",
// 					Accuracy:        insights.Accuracy__High,
// 					Incidents: []*repositories.ApplicationIncident{
// 						{
// 							ApplicationName: "test-app",
// 							Title:           "test-title",
// 							Category:        "test-category",
// 							Summary:         "test-summary",
// 							Recommendation:  "test-recommendation",
// 							Urgency:         insights.Urgency_Medium,
// 							ClusterId:       "test-cluster",
// 							CustomPrompt:    "Custom prompt for test-app",
// 							Accuracy:        insights.Accuracy__High,
// 							Sources: []repositories.ApplicationIncidentSource{
// 								{
// 									Timestamp:     0,
// 									Content:       "Log content",
// 									PodName:       "test-pod",
// 									ContainerName: "test-container",
// 									Image:         "test-image",
// 								},
// 							},
// 						},
// 						{
// 							ApplicationName: "test-app",
// 							Title:           "test-title-2",
// 							Category:        "test-category-2",
// 							Summary:         "test-summary-2",
// 							Recommendation:  "test-recommendation-2",
// 							Urgency:         insights.Urgency_Medium,
// 							ClusterId:       "test-cluster",
// 							CustomPrompt:    "Custom prompt for test-app",
// 							Accuracy:        insights.Accuracy__High,
// 							Sources: []repositories.ApplicationIncidentSource{
// 								{
// 									Timestamp:     2,
// 									Content:       "Log content 2",
// 									PodName:       "test-pod",
// 									ContainerName: "test-container",
// 									Image:         "test-image",
// 								},
// 							},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}
//
// 	test := func(dependencies TestDependencies) {
//
// 		for _, tc := range testCases {
// 			t.Run(tc.description, func(t *testing.T) {
// 				reports, err := dependencies.ReportsService.GetApplicationReportsFromInsights(
// 					tc.applicationInsights, tc.applicationConfigurations,
// 				)
//
// 				assert.NoError(t, err)
// 				assert.Equal(t, len(tc.expectedReports), len(reports), tc.description)
// 				for idx, report := range reports {
// 					expectedReport := tc.expectedReports[idx]
// 					assert.Equal(t, expectedReport.Accuracy, report.Accuracy)
// 					assert.Equal(t, expectedReport.ApplicationName, report.ApplicationName)
// 					assert.Equal(t, expectedReport.CustomPrompt, report.CustomPrompt)
//
// 					assert.Equal(t, len(expectedReport.Incidents), len(report.Incidents))
// 					for incidentId, incident := range report.Incidents {
// 						expectedIncident := expectedReport.Incidents[incidentId]
// 						expectedIncident.Id = incident.Id
// 						assert.Equal(t, expectedIncident, incident)
// 					}
// 				}
// 			})
// 		}
// 	}
//
// 	tests.RunTest(test, t, config.AppModule)
// }
