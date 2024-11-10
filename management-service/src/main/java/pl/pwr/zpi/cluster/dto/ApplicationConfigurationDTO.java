package pl.pwr.zpi.cluster.dto;

import pl.pwr.zpi.reports.enums.Accuracy;

public record ApplicationConfigurationDTO(
        String name,
        String kind,
        Accuracy accuracy,
        String customPrompt
) {
}
