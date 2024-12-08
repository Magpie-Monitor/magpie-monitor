package pl.pwr.zpi.reports.dto.report;

import lombok.Builder;
import pl.pwr.zpi.reports.dto.report.application.ApplicationIncidentSimplifiedDTO;
import pl.pwr.zpi.reports.dto.report.node.NodeIncidentSimplifiedDTO;
import pl.pwr.zpi.reports.entity.report.application.ApplicationIncident;
import pl.pwr.zpi.reports.entity.report.node.NodeIncident;

import java.util.List;

@Builder
public record ReportIncidentsDTO(
        List<ApplicationIncidentSimplifiedDTO> applicationIncidents,
        List<NodeIncidentSimplifiedDTO> nodeIncidents
) {
}
