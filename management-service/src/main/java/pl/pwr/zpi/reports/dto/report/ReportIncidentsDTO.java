package pl.pwr.zpi.reports.dto.report;

import lombok.Builder;
import pl.pwr.zpi.reports.entity.report.application.ApplicationIncident;
import pl.pwr.zpi.reports.entity.report.node.NodeIncident;

import java.util.List;

@Builder
public record ReportIncidentsDTO(
        List<ApplicationIncident> applicationIncidents,
        List<NodeIncident> nodeIncidents
) {
}
