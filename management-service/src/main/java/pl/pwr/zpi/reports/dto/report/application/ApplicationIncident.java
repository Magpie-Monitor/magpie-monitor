package pl.pwr.zpi.reports.dto.report.application;

import pl.pwr.zpi.reports.Urgency;

import java.util.List;

public record ApplicationIncident(
        String id,
        String clusterId,
        String applicationName,
        String category,
        String summary,
        String recommendation,
        Urgency urgency,
        List<ApplicationIncidentSource> sources
) {
}
