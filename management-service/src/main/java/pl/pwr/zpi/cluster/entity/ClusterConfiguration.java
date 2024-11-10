package pl.pwr.zpi.cluster.entity;

import jakarta.persistence.CascadeType;
import jakarta.persistence.Entity;
import jakarta.persistence.Id;
import jakarta.persistence.OneToMany;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;
import pl.pwr.zpi.cluster.dto.UpdateClusterConfigurationRequest;
import pl.pwr.zpi.notifications.discord.entity.DiscordReceiver;
import pl.pwr.zpi.notifications.email.entity.EmailReceiver;
import pl.pwr.zpi.notifications.slack.entity.SlackReceiver;
import pl.pwr.zpi.reports.enums.Accuracy;

import java.util.List;

@Data
@Entity
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class ClusterConfiguration {
    @Id
    private String id;
    private Accuracy accuracy;
    private boolean isEnabled;
    private Long sinceMs;
    private Long toMs;

    @OneToMany
    private List<SlackReceiver> slackReceivers;
    @OneToMany
    private List<DiscordReceiver> discordReceivers;
    @OneToMany
    private List<EmailReceiver> emailReceivers;

    @OneToMany(cascade = CascadeType.ALL)
    private List<ApplicationConfiguration> applicationConfigurations;
    @OneToMany(cascade = CascadeType.ALL)
    private List<NodeConfiguration> nodeConfigurations;

    public static ClusterConfiguration ofClusterConfigurationRequest(UpdateClusterConfigurationRequest configurationRequest) {
        return ClusterConfiguration.builder()
                .id(configurationRequest.id())
                .accuracy(configurationRequest.accuracy())
                .isEnabled(configurationRequest.isEnabled())
                .sinceMs(configurationRequest.sinceMs())
                .toMs(configurationRequest.toMs())
                .applicationConfigurations(configurationRequest.applicationConfigurations())
                .nodeConfigurations(configurationRequest.nodeConfigurations())
                .build();
    }
}