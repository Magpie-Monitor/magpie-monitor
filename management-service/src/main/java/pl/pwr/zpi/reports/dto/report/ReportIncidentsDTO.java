package pl.pwr.zpi.reports.dto.report;

import pl.pwr.zpi.reports.dto.report.application.ApplicationIncidentDTO;
import pl.pwr.zpi.reports.dto.report.node.NodeIncidentDTO;
import pl.pwr.zpi.reports.entity.report.application.ApplicationIncident;
import pl.pwr.zpi.reports.entity.report.node.NodeIncident;

import java.util.List;

public record ReportIncidentsDTO(
        List<ApplicationIncident> applicationIncidents,
        List<NodeIncident> nodeIncidents
) {
}
