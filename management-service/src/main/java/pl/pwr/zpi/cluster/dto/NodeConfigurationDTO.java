package pl.pwr.zpi.cluster.dto;

import pl.pwr.zpi.reports.enums.Accuracy;

public record NodeConfigurationDTO(
        String name,
        Accuracy accuracy,
        String customPrompt
) {
}
