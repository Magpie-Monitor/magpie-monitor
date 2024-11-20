package pl.pwr.zpi.cluster.service;

import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.cluster.dto.ClusterConfigurationDTO;
import pl.pwr.zpi.cluster.dto.UpdateClusterConfigurationRequest;
import pl.pwr.zpi.cluster.dto.UpdateClusterConfigurationResponse;
import pl.pwr.zpi.cluster.entity.ClusterConfiguration;
import pl.pwr.zpi.cluster.repository.ClusterRepository;
import pl.pwr.zpi.metadata.dto.cluster.ClusterMetadataDTO;
import pl.pwr.zpi.metadata.service.MetadataService;
import pl.pwr.zpi.notifications.ReceiverService;
import pl.pwr.zpi.notifications.discord.entity.DiscordReceiver;
import pl.pwr.zpi.notifications.email.entity.EmailReceiver;
import pl.pwr.zpi.notifications.slack.entity.SlackReceiver;
import pl.pwr.zpi.reports.dto.request.CreateReportRequest;
import pl.pwr.zpi.reports.service.ReportGenerationService;

import java.util.List;
import java.util.Optional;

@Service
@RequiredArgsConstructor
public class ClusterService {

    private final ClusterRepository clusterRepository;
    private final ReceiverService receiverService;
    private final MetadataService metadataService;
    private final ReportGenerationService reportGenerationService;

    public void generateReportForCluster(String clusterId, Long sinceMs, Long toMs) {
        clusterRepository.findById(clusterId).ifPresentOrElse(
                clusterConfiguration -> generateReportForClusterConfiguration(clusterConfiguration, sinceMs, toMs),
                () -> {
                    throw new RuntimeException(String.format("Report configuration not found for cluster of an id: %s", clusterId));
                }
        );
    }

    private void generateReportForClusterConfiguration(ClusterConfiguration clusterConfiguration, Long sinceMs, Long toMs) {
        CreateReportRequest createReportRequest =
                CreateReportRequest.fromClusterConfiguration(clusterConfiguration, sinceMs, toMs);
        reportGenerationService.createReport(createReportRequest);
    }

    public UpdateClusterConfigurationResponse updateClusterConfiguration(UpdateClusterConfigurationRequest configurationRequest) {
        ClusterConfiguration clusterConfiguration = ClusterConfiguration.ofClusterConfigurationRequest(configurationRequest);
        setClusterNotificationReceivers(clusterConfiguration, configurationRequest);
        clusterRepository.save(clusterConfiguration);
        return new UpdateClusterConfigurationResponse(clusterConfiguration.getId());
    }

    private void setClusterNotificationReceivers(ClusterConfiguration clusterConfiguration, UpdateClusterConfigurationRequest configurationRequest) {
        clusterConfiguration.setSlackReceivers(getSlackReceiversByIds(configurationRequest.slackReceiverIds()));
        clusterConfiguration.setDiscordReceivers(getDiscordReceiversByIds(configurationRequest.discordReceiverIds()));
        clusterConfiguration.setEmailReceivers(getEmailReceiversByIds(configurationRequest.emailReceiverIds()));
    }

    private List<SlackReceiver> getSlackReceiversByIds(List<Long> receiverIds) {
        return receiverIds.stream()
                .map(receiverService::getSlackReceiverById)
                .toList();
    }

    private List<DiscordReceiver> getDiscordReceiversByIds(List<Long> receiverIds) {
        return receiverIds.stream()
                .map(receiverService::getDiscordReceiverById)
                .toList();
    }

    private List<EmailReceiver> getEmailReceiversByIds(List<Long> receiverIds) {
        return receiverIds.stream()
                .map(receiverService::getEmailReceiverById)
                .toList();
    }

    public Optional<ClusterConfigurationDTO> getClusterById(String clusterId) {
        return clusterRepository.findById(clusterId).map(clusterConfiguration -> {
            Optional<Boolean> isRunning = isClusterRunning(clusterId);
            return isRunning
                    .map(running -> ClusterConfigurationDTO.ofCluster(clusterConfiguration, running))
                    .orElse(ClusterConfigurationDTO.ofCluster(clusterConfiguration, false));

        });
    }

    private Optional<Boolean> isClusterRunning(String clusterId) {
        return metadataService.getClusterById(clusterId).map(ClusterMetadataDTO::running);
    }
}
