package pl.pwr.zpi.cluster.dto;

import lombok.Builder;
import pl.pwr.zpi.cluster.entity.ApplicationConfiguration;
import pl.pwr.zpi.cluster.entity.ClusterConfiguration;
import pl.pwr.zpi.cluster.entity.NodeConfiguration;
import pl.pwr.zpi.notifications.discord.entity.DiscordReceiver;
import pl.pwr.zpi.notifications.email.entity.EmailReceiver;
import pl.pwr.zpi.notifications.slack.entity.SlackReceiver;
import pl.pwr.zpi.reports.enums.Accuracy;

import java.util.List;

@Builder
public record ClusterConfigurationDTO(
        String id,
        Accuracy accuracy,
        boolean isEnabled,
        boolean running,
        Long sinceMs,
        Long toMs,
        List<SlackReceiver> slackReceivers,
        List<DiscordReceiver> discordReceivers,
        List<EmailReceiver> emailReceivers,
        List<ApplicationConfiguration> applicationConfigurations,
        List<NodeConfiguration> nodeConfigurations
) {

    public static ClusterConfigurationDTO ofCluster(ClusterConfiguration clusterConfiguration, boolean running) {
        return ClusterConfigurationDTO.builder()
                .id(clusterConfiguration.getId())
                .accuracy(clusterConfiguration.getAccuracy())
                .isEnabled(clusterConfiguration.isEnabled())
                .running(running)
                .sinceMs(clusterConfiguration.getSinceMs())
                .toMs(clusterConfiguration.getToMs())
                .slackReceivers(clusterConfiguration.getSlackReceivers())
                .discordReceivers(clusterConfiguration.getDiscordReceivers())
                .emailReceivers(clusterConfiguration.getEmailReceivers())
                .applicationConfigurations(clusterConfiguration.getApplicationConfigurations())
                .nodeConfigurations(clusterConfiguration.getNodeConfigurations())
                .build();
    }
}
