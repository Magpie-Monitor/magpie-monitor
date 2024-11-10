package pl.pwr.zpi.cluster.dto;

import lombok.Builder;
import pl.pwr.zpi.cluster.entity.ApplicationConfiguration;
import pl.pwr.zpi.cluster.entity.Cluster;
import pl.pwr.zpi.cluster.entity.NodeConfiguration;
import pl.pwr.zpi.notifications.discord.controller.DiscordReceiver;
import pl.pwr.zpi.notifications.email.controller.EmailReceiver;
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
    public static ClusterConfigurationDTO ofCluster(Cluster cluster, boolean running) {
        return ClusterConfigurationDTO.builder()
                .id(cluster.getId())
                .accuracy(cluster.getAccuracy())
                .isEnabled(cluster.isEnabled())
                .running(running)
                .sinceMs(cluster.getSinceMs())
                .toMs(cluster.getToMs())
                .slackReceivers(cluster.getSlackReceivers())
                .discordReceivers(cluster.getDiscordReceivers())
                .emailReceivers(cluster.getEmailReceivers())
                .applicationConfigurations(cluster.getApplicationConfigurations())
                .nodeConfigurations(cluster.getNodeConfigurations())
                .build();
    }
}
