package pl.pwr.zpi.reports.dto.report.node;

import pl.pwr.zpi.reports.dto.report.Accuracy;

public record ScheduledNodeInsightConfiguration(
        String nodeName,
        Accuracy accuracy,
        String customPrompt
) {
}
