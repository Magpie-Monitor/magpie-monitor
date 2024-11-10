package pl.pwr.zpi.reports.dto.report.node;

import lombok.Builder;
import pl.pwr.zpi.cluster.entity.NodeConfiguration;
import pl.pwr.zpi.reports.enums.Accuracy;

@Builder
public record NodeConfigurationDTO(
        String nodeName,
        Accuracy accuracy,
        String customPrompt
) {
    public static NodeConfigurationDTO fromNodeConfiguration(NodeConfiguration nodeConfiguration) {
        return NodeConfigurationDTO.builder()
                .nodeName(nodeConfiguration.getName())
                .accuracy(nodeConfiguration.getAccuracy())
                .customPrompt(nodeConfiguration.getCustomPrompt())
                .build();
    }
}
