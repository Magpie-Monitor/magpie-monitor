package filter_test

import (
	"testing"

	"github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/filter"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/insights"
	"github.com/stretchr/testify/assert"
)

func TestFilterByApplicationsAccuracy(
	t *testing.T,
) {

	testCases := []struct {
		description                string
		logsByApplication          map[string][]*repositories.ApplicationLogsDocument
		configurationByApplication map[string]*insights.ApplicationInsightConfiguration
		expectedLogs               map[string][]*repositories.ApplicationLogsDocument
	}{
		{
			description: "Filter all logs without key words withing High accuracy",
			logsByApplication: map[string][]*repositories.ApplicationLogsDocument{
				"test-app-1": []*repositories.ApplicationLogsDocument{
					{
						Id:      "test-log-1",
						Content: "This logs will be filtered",
					},
					{
						Id:      "test-log-2",
						Content: "This is a log with High Disk Usage",
					},
				},

				"test-app-2": []*repositories.ApplicationLogsDocument{
					{
						Id:      "test-log-3",
						Content: "This is a medium log. CPU Usage",
					},
					{
						Id:      "test-log-4",
						Content: "This logs will be filtered",
					},
				},
			},
			configurationByApplication: map[string]*insights.ApplicationInsightConfiguration{
				"test-app-1": {
					ApplicationName: "test-app-1",
					Accuracy:        insights.Accuracy__High,
				},
				"test-app-2": {
					ApplicationName: "test-app-2",
					Accuracy:        insights.Accuracy__Medium,
				},
			},

			expectedLogs: map[string][]*repositories.ApplicationLogsDocument{
				"test-app-1": []*repositories.ApplicationLogsDocument{
					{
						Id:      "test-log-2",
						Content: "This is a log with High Disk Usage",
					},
				},

				"test-app-2": []*repositories.ApplicationLogsDocument{
					{
						Id:      "test-log-3",
						Content: "This is a medium log. CPU Usage",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			filter.FilterByApplicationsAccuracy(tc.logsByApplication, tc.configurationByApplication)
			assert.Equal(t, tc.expectedLogs, tc.logsByApplication)
		})
	}
}
