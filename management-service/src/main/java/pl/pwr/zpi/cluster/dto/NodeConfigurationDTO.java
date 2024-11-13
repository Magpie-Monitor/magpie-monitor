package pl.pwr.zpi.cluster.dto;

import lombok.Builder;
import pl.pwr.zpi.cluster.entity.NodeConfiguration;
import pl.pwr.zpi.reports.enums.Accuracy;

@Builder
public record NodeConfigurationDTO(
        String name,
        Accuracy accuracy,
        String customPrompt
) {

    public static NodeConfigurationDTO ofNodeConfiguration(NodeConfiguration nodeConfiguration) {
        return NodeConfigurationDTO.builder()
                .name(nodeConfiguration.getName())
                .accuracy(nodeConfiguration.getAccuracy())
                .customPrompt(nodeConfiguration.getCustomPrompt())
                .build();
    }
}
