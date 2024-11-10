package pl.pwr.zpi.reports.dto.report.application;

import lombok.Builder;
import pl.pwr.zpi.cluster.entity.ApplicationConfiguration;
import pl.pwr.zpi.reports.enums.Accuracy;

@Builder
public record ApplicationConfigurationDTO(
        String applicationName,
        String customPrompt,
        Accuracy accuracy
) {
    public static ApplicationConfigurationDTO ofApplicationConfiguration(ApplicationConfiguration applicationConfiguration) {
        return ApplicationConfigurationDTO.builder()
                .applicationName(applicationConfiguration.getName())
                .customPrompt(applicationConfiguration.getCustomPrompt())
                .accuracy(applicationConfiguration.getAccuracy())
                .build();
    }
}
