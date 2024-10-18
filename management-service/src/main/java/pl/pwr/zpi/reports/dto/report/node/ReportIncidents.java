package pl.pwr.zpi.reports.dto.report.node;

import pl.pwr.zpi.reports.dto.report.application.ApplicationIncident;

import java.util.List;

public record ReportIncidents(
        List<ApplicationIncident> applicationIncidents,
        List<NodeIncident> nodeIncidents
) {
}
