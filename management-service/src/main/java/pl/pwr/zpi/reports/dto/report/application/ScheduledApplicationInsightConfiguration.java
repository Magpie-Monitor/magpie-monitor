package pl.pwr.zpi.reports.dto.report.application;

public record ScheduledApplicationInsightConfiguration(
        String applicationName,
        String precision,
        String customPrompt
) {
}
