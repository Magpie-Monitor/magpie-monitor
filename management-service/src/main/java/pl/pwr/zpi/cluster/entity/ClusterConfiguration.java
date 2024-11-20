package pl.pwr.zpi.cluster.entity;

import jakarta.persistence.*;
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
    private Long generatedEveryMillis;

    @OneToMany(fetch = FetchType.EAGER)
    private List<SlackReceiver> slackReceivers;
    @OneToMany(fetch = FetchType.EAGER)
    private List<DiscordReceiver> discordReceivers;
    @OneToMany(fetch = FetchType.EAGER)
    private List<EmailReceiver> emailReceivers;

    @OneToMany(cascade = CascadeType.ALL, fetch = FetchType.EAGER)
    private List<ApplicationConfiguration> applicationConfigurations;
    @OneToMany(cascade = CascadeType.ALL, fetch = FetchType.EAGER)
    private List<NodeConfiguration> nodeConfigurations;

    public List<Long> getSlackReceiverIds() {
        return slackReceivers.stream()
                .map(SlackReceiver::getId)
                .toList();
    }

    public List<Long> getDiscordReceiverIds() {
        return discordReceivers.stream()
                .map(DiscordReceiver::getId)
                .toList();
    }

    public List<Long> getEmailReceiverIds() {
        return emailReceivers.stream()
                .map(EmailReceiver::getId)
                .toList();
    }

    public static ClusterConfiguration ofClusterConfigurationRequest(UpdateClusterConfigurationRequest configurationRequest) {
        return ClusterConfiguration.builder()
                .id(configurationRequest.id())
                .accuracy(configurationRequest.accuracy())
                .isEnabled(configurationRequest.isEnabled())
                .generatedEveryMillis(configurationRequest.generatedEveryMillis())
                .applicationConfigurations(
                        configurationRequest.applicationConfigurations().stream()
                                .map(ApplicationConfiguration::fromApplicationConfigurationDTO)
                                .toList()
                )
                .nodeConfigurations(
                        configurationRequest.nodeConfigurations().stream()
                                .map(NodeConfiguration::fromNodeConfigurationDTO)
                                .toList()
                )
                .build();
    }
}