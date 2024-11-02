package pl.pwr.zpi.reports.dto.report.node;

public record ScheduledNodeInsightConfigurationDTO(
        String nodeName,
        String precision,
        String customPrompt
) {
}
