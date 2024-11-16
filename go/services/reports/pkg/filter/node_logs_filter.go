package filter

import (
	"github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/insights"
)

func FilterByNodesAccuracy(
	logsByNode map[string][]*repositories.NodeLogsDocument,
	configurationByNode map[string]*insights.NodeInsightConfiguration,
) {
	for node, logs := range logsByNode {
		config, ok := configurationByNode[node]
		var accuracy insights.Accuracy
		if !ok {
			// By default the node is not included
			delete(logsByNode, node)
		} else {
			accuracy = config.Accuracy
			filter := NewAccuracyFilter[*repositories.NodeLogsDocument](accuracy)
			logsByNode[node] = filter.Filter(logs)
		}
	}
}
