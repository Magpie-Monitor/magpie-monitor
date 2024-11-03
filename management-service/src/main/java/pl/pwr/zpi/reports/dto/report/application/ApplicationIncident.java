package pl.pwr.zpi.reports.dto.report.application;

import pl.pwr.zpi.reports.dto.report.Accuracy;
import pl.pwr.zpi.reports.dto.report.Urgency;

import java.util.List;

public record ApplicationIncident(
        String id,
        String title,
        String clusterId,
        Accuracy accuracy,
        String customPrompt,
        String applicationName,
        String category,
        String summary,
        String recommendation,
        Urgency urgency,
        List<ApplicationIncidentSource> sources
) {
}
