package incidentcorrelation_test

import (
	incidentcorrelation "github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/incident_correlation"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/insights"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/repositories"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMergeApplicationIncidentsByGroups(t *testing.T) {

	testCases := []struct {
		description          string
		groups               []incidentcorrelation.IncidentMergeGroup
		applicationIncidents *repositories.ApplicationReport
		expectedNewIncidents []*repositories.ApplicationIncident
	}{
		{
			description: "Join incidents by merging their sources",
			groups: []incidentcorrelation.IncidentMergeGroup{
				{
					Title:          "group-1-title",
					Summary:        "group-1-summary",
					Recommendation: "group-1-recommendation",
					Category:       "group-1-category",
					Urgency:        insights.Urgency_Medium,
					IncidentIds:    []string{"incident-1", "incident-3"},
				},
			},
			applicationIncidents: &repositories.ApplicationReport{
				Incidents: []*repositories.ApplicationIncident{
					{
						Id:             "incident-1",
						ClusterId:      "test-cluster",
						Title:          "group-1-title",
						Summary:        "group-1-summary",
						Recommendation: "group-1-recommendation",
						Category:       "group-1-category",
						Urgency:        insights.Urgency_Medium,
						Sources: []repositories.ApplicationIncidentSource{
							{
								Timestamp:     0,
								PodName:       "test-pod",
								ContainerName: "test-container",
								Image:         "test-image",
								SourceLog:     "Log 1",
							},
							{
								Timestamp:     1,
								PodName:       "test-pod",
								ContainerName: "test-container",
								Image:         "test-image",
								SourceLog:     "Log 2",
							},
						},
					},
					{
						Id:             "incident-2",
						ClusterId:      "test-cluster",
						Title:          "group-1-title",
						Summary:        "group-1-summary",
						Recommendation: "group-1-recommendation",
						Category:       "group-1-category",
						Urgency:        insights.Urgency_Medium,
						Sources: []repositories.ApplicationIncidentSource{
							{
								Timestamp:     2,
								PodName:       "test-pod",
								ContainerName: "test-container",
								Image:         "test-image",
								SourceLog:     "Log 3",
							},
							{
								Timestamp:     3,
								PodName:       "test-pod",
								ContainerName: "test-container",
								Image:         "test-image",
								SourceLog:     "Log 4",
							},
						},
					},
					{
						Id:             "incident-3",
						ClusterId:      "test-cluster",
						Title:          "group-1-title",
						Summary:        "group-1-summary",
						Recommendation: "group-1-recommendation",
						Category:       "group-1-category",
						Urgency:        insights.Urgency_Medium,
						Sources: []repositories.ApplicationIncidentSource{
							{
								Timestamp:     4,
								PodName:       "test-pod",
								ContainerName: "test-container",
								Image:         "test-image",
								SourceLog:     "Log 5",
							},
							{
								Timestamp:     5,
								PodName:       "test-pod",
								ContainerName: "test-container",
								Image:         "test-image",
								SourceLog:     "Log 6",
							},
						},
					},
				},
			},
			expectedNewIncidents: []*repositories.ApplicationIncident{
				{
					Title:          "group-1-title",
					ClusterId:      "test-cluster",
					Summary:        "group-1-summary",
					Recommendation: "group-1-recommendation",
					Category:       "group-1-category",
					Urgency:        insights.Urgency_Medium,
					Sources: []repositories.ApplicationIncidentSource{
						{
							Timestamp:     0,
							PodName:       "test-pod",
							ContainerName: "test-container",
							Image:         "test-image",
							SourceLog:     "Log 1",
						},
						{
							Timestamp:     1,
							PodName:       "test-pod",
							ContainerName: "test-container",
							Image:         "test-image",
							SourceLog:     "Log 2",
						},

						{
							Timestamp:     4,
							PodName:       "test-pod",
							ContainerName: "test-container",
							Image:         "test-image",
							SourceLog:     "Log 5",
						},
						{
							Timestamp:     5,
							PodName:       "test-pod",
							ContainerName: "test-container",
							Image:         "test-image",
							SourceLog:     "Log 6",
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(*testing.T) {

			newIncidents := incidentcorrelation.MergeApplicationIncidentsByGroups(
				tc.groups, tc.applicationIncidents)

			assert.ElementsMatch(t, tc.expectedNewIncidents, newIncidents)

		})
	}

}
func TestMergeNodeIncidentsByGroups(t *testing.T) {

	testCases := []struct {
		description          string
		groups               []incidentcorrelation.IncidentMergeGroup
		nodeIncidents        *repositories.NodeReport
		expectedNewIncidents []*repositories.NodeIncident
	}{
		{
			description: "Join incidents by merging their sources",
			groups: []incidentcorrelation.IncidentMergeGroup{
				{
					Title:          "group-1-title",
					Summary:        "group-1-summary",
					Recommendation: "group-1-recommendation",
					Category:       "group-1-category",
					Urgency:        insights.Urgency_Medium,
					IncidentIds:    []string{"incident-1", "incident-3"},
				},
			},
			nodeIncidents: &repositories.NodeReport{
				Incidents: []*repositories.NodeIncident{
					{
						Id:             "incident-1",
						ClusterId:      "test-cluster",
						Title:          "group-1-title",
						Summary:        "group-1-summary",
						Recommendation: "group-1-recommendation",
						Category:       "group-1-category",
						Urgency:        insights.Urgency_Medium,
						Sources: []repositories.NodeIncidentSource{
							{
								Timestamp: 0,
								Filename:  "test-filename",
								SourceLog: "Log 1",
							},
							{
								Timestamp: 1,
								Filename:  "test-filename",
								SourceLog: "Log 2",
							},
						},
					},
					{
						Id:             "incident-2",
						ClusterId:      "test-cluster",
						Title:          "group-1-title",
						Summary:        "group-1-summary",
						Recommendation: "group-1-recommendation",
						Category:       "group-1-category",
						Urgency:        insights.Urgency_Medium,
						Sources: []repositories.NodeIncidentSource{
							{
								Timestamp: 2,
								Filename:  "test-filename",
								SourceLog: "Log 3",
							},
							{
								Timestamp: 3,
								Filename:  "test-filename",
								SourceLog: "Log 4",
							},
						},
					},
					{
						Id:             "incident-3",
						ClusterId:      "test-cluster",
						Title:          "group-1-title",
						Summary:        "group-1-summary",
						Recommendation: "group-1-recommendation",
						Category:       "group-1-category",
						Urgency:        insights.Urgency_Medium,
						Sources: []repositories.NodeIncidentSource{
							{
								Timestamp: 4,
								Filename:  "test-filename",
								SourceLog: "Log 5",
							},
							{
								Timestamp: 5,
								Filename:  "test-filename",
								SourceLog: "Log 6",
							},
						},
					},
				},
			},
			expectedNewIncidents: []*repositories.NodeIncident{
				{
					Title:          "group-1-title",
					ClusterId:      "test-cluster",
					Summary:        "group-1-summary",
					Recommendation: "group-1-recommendation",
					Category:       "group-1-category",
					Urgency:        insights.Urgency_Medium,
					Sources: []repositories.NodeIncidentSource{
						{
							Timestamp: 0,
							Filename:  "test-filename",
							SourceLog: "Log 1",
						},
						{
							Timestamp: 1,
							Filename:  "test-filename",
							SourceLog: "Log 2",
						},

						{
							Timestamp: 4,
							Filename:  "test-filename",
							SourceLog: "Log 5",
						},
						{
							Timestamp: 5,
							Filename:  "test-filename",
							SourceLog: "Log 6",
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(*testing.T) {

			newIncidents := incidentcorrelation.MergeNodeIncidentsByGroups(
				tc.groups, tc.nodeIncidents)

			assert.ElementsMatch(t, tc.expectedNewIncidents, newIncidents)

		})
	}

}
