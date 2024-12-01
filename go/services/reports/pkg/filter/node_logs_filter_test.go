package filter_test

import (
	"github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/filter"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/insights"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFilterByNodesAccuracy(
	t *testing.T,
) {

	testCases := []struct {
		description         string
		logsByNode          map[string][]*repositories.NodeLogsDocument
		configurationByNode map[string]*insights.NodeInsightConfiguration
		expectedLogs        map[string][]*repositories.NodeLogsDocument
	}{
		{
			description: "Filter all logs without key words withing High accuracy",
			logsByNode: map[string][]*repositories.NodeLogsDocument{
				"test-node-1": []*repositories.NodeLogsDocument{
					{
						Id:      "test-log-1",
						Content: "This logs will be filtered",
					},
					{
						Id:      "test-log-2",
						Content: "This is a log with High Disk Usage",
					},
				},

				"test-node-2": []*repositories.NodeLogsDocument{
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
			configurationByNode: map[string]*insights.NodeInsightConfiguration{
				"test-node-1": {
					NodeName: "test-node-1",
					Accuracy: insights.Accuracy__High,
				},
				"test-node-2": {
					NodeName: "test-node-2",
					Accuracy: insights.Accuracy__Medium,
				},
			},

			expectedLogs: map[string][]*repositories.NodeLogsDocument{
				"test-node-1": []*repositories.NodeLogsDocument{
					{
						Id:      "test-log-2",
						Content: "This is a log with High Disk Usage",
					},
				},

				"test-node-2": []*repositories.NodeLogsDocument{
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
			filter.FilterByNodesAccuracy(tc.logsByNode, tc.configurationByNode)
			assert.Equal(t, tc.expectedLogs, tc.logsByNode)
		})
	}
}
