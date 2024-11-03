package pl.pwr.zpi.reports.dto.report.application;

public record ScheduledApplicationInsightConfiguration(
        String applicationName,
        String accuracy,
        String customPrompt
) {
}
