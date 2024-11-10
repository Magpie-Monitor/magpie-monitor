package pl.pwr.zpi.cluster.dto;

import pl.pwr.zpi.cluster.entity.ApplicationConfiguration;
import pl.pwr.zpi.cluster.entity.NodeConfiguration;
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
        List<ApplicationConfiguration> applicationConfigurations,
        List<NodeConfiguration> nodeConfigurations
) {

}
