package pl.pwr.zpi.reports.dto.report.application;

import pl.pwr.zpi.reports.enums.Urgency;

import java.util.List;

public record ApplicationIncidentDTO(
        String id,
        String clusterId,
        String applicationName,
        String category,
        String summary,
        String recommendation,
        Urgency urgency,
        List<ApplicationIncidentSourceDTO> sources
) {
}
