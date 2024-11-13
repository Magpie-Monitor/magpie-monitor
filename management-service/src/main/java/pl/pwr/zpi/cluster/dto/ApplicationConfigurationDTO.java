package pl.pwr.zpi.cluster.dto;

import lombok.Builder;
import pl.pwr.zpi.cluster.entity.ApplicationConfiguration;
import pl.pwr.zpi.reports.enums.Accuracy;

@Builder
public record ApplicationConfigurationDTO(
        String name,
        String kind,
        Accuracy accuracy,
        String customPrompt
) {
    public static ApplicationConfigurationDTO ofApplicationConfiguration(ApplicationConfiguration applicationConfiguration) {
        return ApplicationConfigurationDTO.builder()
                .name(applicationConfiguration.getName())
                .kind(applicationConfiguration.getKind())
                .accuracy(applicationConfiguration.getAccuracy())
                .customPrompt(applicationConfiguration.getCustomPrompt())
                .build();
    }
}
