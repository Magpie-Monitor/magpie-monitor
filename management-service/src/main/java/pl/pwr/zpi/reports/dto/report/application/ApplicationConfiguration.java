package pl.pwr.zpi.reports.dto.report.application;

import pl.pwr.zpi.reports.enums.Accuracy;

public record ApplicationConfiguration(
        String applicationName,
        String customPrompt,
        Accuracy accuracy
) {
}
