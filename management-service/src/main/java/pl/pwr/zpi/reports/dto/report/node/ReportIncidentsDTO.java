package pl.pwr.zpi.reports.dto.report.node;

import pl.pwr.zpi.reports.dto.report.application.ApplicationIncidentDTO;

import java.util.List;

public record ReportIncidentsDTO(
        List<ApplicationIncidentDTO> applicationIncidents,
        List<NodeIncidentDTO> nodeIncidents
) {
}
