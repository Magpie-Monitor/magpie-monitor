package pl.pwr.zpi.reports.dto.report.application;

import pl.pwr.zpi.reports.dto.report.Accuracy;

public record ScheduledApplicationInsightConfiguration(
        String applicationName,
        Accuracy accuracy,
        String customPrompt
) {
}
