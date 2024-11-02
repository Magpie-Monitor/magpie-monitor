package pl.pwr.zpi.reports.dto.report.application;

public record ScheduledApplicationInsightConfigurationDTO(
        String applicationName,
        String precision,
        String customPrompt
) {
}
