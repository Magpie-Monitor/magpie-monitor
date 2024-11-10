package pl.pwr.zpi.cluster.dto;

import pl.pwr.zpi.reports.enums.Accuracy;

import java.util.List;

public record UpdateClusterConfigurationRequest(
        String id,
        Accuracy accuracy,
        boolean isEnabled,
        Long generatedEveryMillis,
        List<Long> slackReceiverIds,
        List<Long> discordReceiverIds,
        List<Long> emailReceiverIds,
        List<ApplicationConfigurationDTO> applicationConfigurations,
        List<NodeConfigurationDTO> nodeConfigurations
) {

}
