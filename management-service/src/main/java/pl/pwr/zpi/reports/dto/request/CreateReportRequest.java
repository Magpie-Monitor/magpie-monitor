package pl.pwr.zpi.reports.dto.request;

import lombok.Builder;
import lombok.NonNull;
import pl.pwr.zpi.cluster.entity.ClusterConfiguration;
import pl.pwr.zpi.reports.dto.report.application.ApplicationConfigurationDTO;
import pl.pwr.zpi.reports.dto.report.node.NodeConfigurationDTO;
import pl.pwr.zpi.reports.enums.Accuracy;

import java.util.List;

@Builder
public record CreateReportRequest(
        @NonNull
        String clusterId,
        @NonNull
        Accuracy accuracy,
        @NonNull
        Long sinceMs,
        @NonNull
        Long toMs,
        @NonNull
        List<Long> slackReceiverIds,
        @NonNull
        List<Long> discordReceiverIds,
        @NonNull
        List<Long> emailReceiverIds,
        @NonNull
        List<ApplicationConfigurationDTO> applicationConfigurations,
        @NonNull
        List<NodeConfigurationDTO> nodeConfigurations
) {

    public static CreateReportRequest fromClusterConfiguration(
            ClusterConfiguration clusterConfiguration, Long sinceMs, Long toMs) {
        return CreateReportRequest.builder()
                .clusterId(clusterConfiguration.getId())
                .accuracy(clusterConfiguration.getAccuracy())
                .sinceMs(sinceMs)
                .toMs(toMs)
                .slackReceiverIds(clusterConfiguration.getSlackReceiverIds())
                .discordReceiverIds(clusterConfiguration.getDiscordReceiverIds())
                .emailReceiverIds(clusterConfiguration.getEmailReceiverIds())
                .applicationConfigurations(
                        clusterConfiguration.getApplicationConfigurations().stream()
                                .map(ApplicationConfigurationDTO::ofApplicationConfiguration)
                                .toList()
                )
                .nodeConfigurations(
                        clusterConfiguration.getNodeConfigurations().stream()
                                .map(NodeConfigurationDTO::fromNodeConfiguration)
                                .toList()
                )
                .build();
    }
}
