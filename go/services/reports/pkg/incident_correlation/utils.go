package incidentcorrelation

import (
	"slices"

	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/repositories"
)

func ConvertConcreteIncidentArrayIntoIncidents[T repositories.Incident](incidents []T) []repositories.Incident {
	convertedIncidents := make([]repositories.Incident, 0, len(incidents))

	for _, incident := range incidents {
		convertedIncidents = append(convertedIncidents, incident)
	}

	return convertedIncidents
}

// Removes merged incidents from report and returns new incident as a result
func mergeApplicationIncidents(group IncidentMergeGroup, applicationIncidents *repositories.ApplicationReport) *repositories.ApplicationIncident {
	sourcesToBeMerged := make([]string, 0)
	filteredIncidents := make([]*repositories.ApplicationIncident, 0)

	if len(applicationIncidents.Incidents) == 0 {
		panic("Empty list of incidents from a application report was passed to an incident merger")
	}

	for _, incident := range applicationIncidents.Incidents {
		if !slices.Contains(group.IncidentIds, incident.Id) {
			filteredIncidents = append(filteredIncidents, incident)
			continue
		}

		sourcesToBeMerged = append(sourcesToBeMerged, incident.SourceIds...)
	}

	newIncident := &repositories.ApplicationIncident{
		Title:           group.Title,
		Urgency:         group.Urgency,
		Summary:         group.Summary,
		Recommendation:  group.Recommendation,
		Category:        group.Category,
		CustomPrompt:    applicationIncidents.CustomPrompt,
		Accuracy:        applicationIncidents.Accuracy,
		ApplicationName: applicationIncidents.ApplicationName,
		ClusterId:       applicationIncidents.Incidents[0].ClusterId,
		// Sources:         sourcesToBeMerged,
		SourceIds: sourcesToBeMerged,
	}

	applicationIncidents.Incidents = filteredIncidents

	return newIncident
}
func MergeApplicationIncidentsByGroups(groups []IncidentMergeGroup, applicationIncidents *repositories.ApplicationReport) []*repositories.ApplicationIncident {
	var newIncidents []*repositories.ApplicationIncident
	for _, group := range groups {
		newIncident := mergeApplicationIncidents(group, applicationIncidents)
		newIncidents = append(newIncidents, newIncident)
	}

	return newIncidents
}

// Removes merged incidents from report and returns new incident as a result
func mergeNodeIncidents(group IncidentMergeGroup, nodeIncidents *repositories.NodeReport) *repositories.NodeIncident {
	sourcesToBeMerged := make([]string, 0)
	filteredIncidents := make([]*repositories.NodeIncident, 0)

	if len(nodeIncidents.Incidents) == 0 {
		panic("Empty list of incidents from a node report was passed to an incident merger")
	}

	for _, incident := range nodeIncidents.Incidents {
		if !slices.Contains(group.IncidentIds, incident.Id) {
			filteredIncidents = append(filteredIncidents, incident)
			continue
		}

		sourcesToBeMerged = append(sourcesToBeMerged, incident.SourceIds...)
	}

	newIncident := &repositories.NodeIncident{
		Title:          group.Title,
		Urgency:        group.Urgency,
		Summary:        group.Summary,
		Recommendation: group.Recommendation,
		Category:       group.Category,
		CustomPrompt:   nodeIncidents.CustomPrompt,
		Accuracy:       nodeIncidents.Accuracy,
		NodeName:       nodeIncidents.Node,
		ClusterId:      nodeIncidents.Incidents[0].ClusterId,
		SourceIds:      sourcesToBeMerged,
	}

	nodeIncidents.Incidents = filteredIncidents

	return newIncident
}

func MergeNodeIncidentsByGroups(groups []IncidentMergeGroup, nodeIncident *repositories.NodeReport) []*repositories.NodeIncident {
	var newIncidents []*repositories.NodeIncident
	for _, group := range groups {
		newIncident := mergeNodeIncidents(group, nodeIncident)
		newIncidents = append(newIncidents, newIncident)
	}

	return newIncidents
}
