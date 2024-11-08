package pl.pwr.zpi.reports.dto.report.node;

import pl.pwr.zpi.reports.enums.Accuracy;

public record NodeConfiguration(
        String nodeName,
        String customPrompt,
        Accuracy accuracy
) {
}
