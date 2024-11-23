package pl.pwr.zpi.metadata.dto.cluster;

import lombok.Builder;
import lombok.Data;
import pl.pwr.zpi.notifications.discord.entity.DiscordReceiver;
import pl.pwr.zpi.notifications.email.entity.EmailReceiver;
import pl.pwr.zpi.notifications.slack.entity.SlackReceiver;
import pl.pwr.zpi.reports.enums.Accuracy;

import java.util.List;

@Data
@Builder
public class ClusterMetadataDTO {
    String clusterId;
    Long updatedAtMillis;
    Accuracy accuracy;
    boolean running;
    List<SlackReceiver> slackReceivers;
    List<DiscordReceiver> discordReceivers;
    List<EmailReceiver> emailReceivers;

    public static ClusterMetadataDTO of(String clusterId, boolean running) {
        return ClusterMetadataDTO.builder()
                .clusterId(clusterId)
                .running(running)
                .build();
    }
}
