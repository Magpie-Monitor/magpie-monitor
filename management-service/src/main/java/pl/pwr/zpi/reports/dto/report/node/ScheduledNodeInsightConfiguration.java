package pl.pwr.zpi.reports.dto.report.node;

public record ScheduledNodeInsightConfiguration(
        String nodeName,
        String precision,
        String customPrompt
) {
}
