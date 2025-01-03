package pl.pwr.zpi.cluster.service;

import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.cluster.dto.ClusterConfigurationDTO;
import pl.pwr.zpi.cluster.dto.UpdateClusterConfigurationRequest;
import pl.pwr.zpi.cluster.dto.UpdateClusterConfigurationResponse;
import pl.pwr.zpi.cluster.entity.ClusterConfiguration;
import pl.pwr.zpi.cluster.repository.ClusterRepository;
import pl.pwr.zpi.metadata.dto.application.ApplicationMetadataDTO;
import pl.pwr.zpi.metadata.dto.cluster.ClusterMetadataDTO;
import pl.pwr.zpi.metadata.dto.node.NodeMetadataDTO;
import pl.pwr.zpi.metadata.service.MetadataService;
import pl.pwr.zpi.notifications.ReceiverService;
import pl.pwr.zpi.notifications.discord.entity.DiscordReceiver;
import pl.pwr.zpi.notifications.email.entity.EmailReceiver;
import pl.pwr.zpi.notifications.slack.entity.SlackReceiver;
import pl.pwr.zpi.reports.enums.Accuracy;

import java.util.Collections;
import java.util.List;
import java.util.Optional;

@Service
@RequiredArgsConstructor
public class ClusterService {

    private final ClusterRepository clusterRepository;
    private final ReceiverService receiverService;
    private final MetadataService metadataService;

    public List<ClusterMetadataDTO> getAllClusters() {
        return metadataService.getAllClusters().stream()
                .peek(this::setClusterConfigurationForMetadata)
                .toList();
    }

    public void setClusterConfigurationForMetadata(ClusterMetadataDTO metadata) {
        clusterRepository.findById(metadata.getClusterId()).ifPresentOrElse(configuration -> {
            metadata.setAccuracy(configuration.getAccuracy());
            metadata.setUpdatedAtMillis(configuration.getUpdatedAtMillis());
            metadata.setSlackReceivers(configuration.getSlackReceivers());
            metadata.setDiscordReceivers(configuration.getDiscordReceivers());
            metadata.setEmailReceivers(configuration.getEmailReceivers());
        }, () -> setDefaultClusterConfigurationForMetadata(metadata));
    }

    private void setDefaultClusterConfigurationForMetadata(ClusterMetadataDTO metadata) {
        metadata.setAccuracy(Accuracy.HIGH);
        metadata.setUpdatedAtMillis(0L);
        metadata.setSlackReceivers(Collections.emptyList());
        metadata.setDiscordReceivers(Collections.emptyList());
        metadata.setEmailReceivers(Collections.emptyList());
    }

    public List<NodeMetadataDTO> getClusterNodes(String clusterId) {
        return metadataService.getClusterNodes(clusterId);
    }

    public List<ApplicationMetadataDTO> getClusterApplications(String clusterId) {
        return metadataService.getClusterApplications(clusterId);
    }

    public UpdateClusterConfigurationResponse updateClusterConfiguration(UpdateClusterConfigurationRequest configurationRequest) {
        ClusterConfiguration clusterConfiguration = ClusterConfiguration.ofClusterConfigurationRequest(configurationRequest);
        setClusterNotificationReceivers(clusterConfiguration, configurationRequest);
        clusterConfiguration.setUpdatedAtMillis(System.currentTimeMillis());
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

    public ClusterConfigurationDTO getClusterById(String clusterId) {
        return clusterRepository.findById(clusterId).map(clusterConfiguration -> {
            Optional<Boolean> isRunning = isClusterRunning(clusterId);
            return isRunning
                    .map(running -> ClusterConfigurationDTO.ofCluster(clusterConfiguration, running))
                    .orElse(ClusterConfigurationDTO.ofCluster(clusterConfiguration, false));
        }).orElse(ClusterConfigurationDTO.defaultConfiguration());
    }

    private Optional<Boolean> isClusterRunning(String clusterId) {
        return metadataService.getClusterById(clusterId).map(ClusterMetadataDTO::isRunning);
    }
}
