package filter

import (
	"github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/insights"
)

func FilterByApplicationsAccuracy(
	logsByApplication map[string][]*repositories.ApplicationLogsDocument,
	configurationByApplication map[string]*insights.ApplicationInsightConfiguration,
) {

	for application, logs := range logsByApplication {
		// Check if application configuration is in params
		config, ok := configurationByApplication[application]
		var accuracy insights.Accuracy
		if !ok {
			// By default the app is not included
			delete(logsByApplication, application)
		} else {
			accuracy = config.Accuracy
			filter := NewAccuracyFilter[*repositories.ApplicationLogsDocument](accuracy)
			logsByApplication[application] = filter.Filter(logs)
		}
	}
}
