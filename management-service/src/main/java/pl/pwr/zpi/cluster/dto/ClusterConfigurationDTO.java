package pl.pwr.zpi.cluster.dto;

import lombok.Builder;
import pl.pwr.zpi.cluster.entity.ClusterConfiguration;
import pl.pwr.zpi.notifications.discord.entity.DiscordReceiver;
import pl.pwr.zpi.notifications.email.entity.EmailReceiver;
import pl.pwr.zpi.notifications.slack.entity.SlackReceiver;
import pl.pwr.zpi.reports.enums.Accuracy;

import java.util.Collections;
import java.util.List;

@Builder
public record ClusterConfigurationDTO(
        String id,
        Accuracy accuracy,
        boolean isEnabled,
        boolean running,
        Long generatedEveryMillis,
        List<SlackReceiver> slackReceivers,
        List<DiscordReceiver> discordReceivers,
        List<EmailReceiver> emailReceivers,
        List<ApplicationConfigurationDTO> applicationConfigurations,
        List<NodeConfigurationDTO> nodeConfigurations,
        Long updatedAtMillis
) {

    public static ClusterConfigurationDTO ofCluster(ClusterConfiguration clusterConfiguration, boolean running) {
        return ClusterConfigurationDTO.builder()
                .id(clusterConfiguration.getId())
                .accuracy(clusterConfiguration.getAccuracy())
                .isEnabled(clusterConfiguration.isEnabled())
                .running(running)
                .generatedEveryMillis(clusterConfiguration.getGeneratedEveryMillis())
                .slackReceivers(clusterConfiguration.getSlackReceivers())
                .discordReceivers(clusterConfiguration.getDiscordReceivers())
                .emailReceivers(clusterConfiguration.getEmailReceivers())
                .applicationConfigurations(
                        clusterConfiguration.getApplicationConfigurations()
                                .stream()
                                .map(ApplicationConfigurationDTO::ofApplicationConfiguration)
                                .toList()
                )
                .nodeConfigurations(
                        clusterConfiguration.getNodeConfigurations()
                                .stream()
                                .map(NodeConfigurationDTO::ofNodeConfiguration)
                                .toList()
                )
                .updatedAtMillis(clusterConfiguration.getUpdatedAtMillis())
                .build();
    }

    public static ClusterConfigurationDTO defaultConfiguration() {
        return ClusterConfigurationDTO.builder()
                .accuracy(Accuracy.HIGH)
                .isEnabled(true)
                .running(true)
                .generatedEveryMillis(0L)
                .slackReceivers(Collections.emptyList())
                .discordReceivers(Collections.emptyList())
                .emailReceivers(Collections.emptyList())
                .applicationConfigurations(Collections.emptyList())
                .nodeConfigurations(Collections.emptyList())
                .updatedAtMillis(0L)
                .build();
    }
}
