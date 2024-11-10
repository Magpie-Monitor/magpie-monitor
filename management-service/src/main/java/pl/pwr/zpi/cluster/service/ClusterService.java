package pl.pwr.zpi.cluster.service;

import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.cluster.dto.ClusterConfigurationDTO;
import pl.pwr.zpi.cluster.dto.ClusterConfigurationRequest;
import pl.pwr.zpi.cluster.dto.ClusterIdResponse;
import pl.pwr.zpi.cluster.entity.Cluster;
import pl.pwr.zpi.cluster.repository.ClusterRepository;
import pl.pwr.zpi.metadata.service.MetadataService;
import pl.pwr.zpi.notifications.ReceiverService;
import pl.pwr.zpi.notifications.discord.entity.DiscordReceiver;
import pl.pwr.zpi.notifications.email.entity.EmailReceiver;
import pl.pwr.zpi.notifications.slack.entity.SlackReceiver;

import java.util.List;
import java.util.Optional;

@Service
@RequiredArgsConstructor
public class ClusterService {

    private final ClusterRepository clusterRepository;
    private final ReceiverService receiverService;
    private final MetadataService metadataService;

    public ClusterIdResponse updateClusterConfiguration(ClusterConfigurationRequest configurationRequest) {
        Cluster cluster = Cluster.ofClusterConfigurationRequest(configurationRequest);
        setClusterNotificationReceivers(cluster, configurationRequest);
        clusterRepository.save(cluster);
        return new ClusterIdResponse(cluster.getId());
    }

    private void setClusterNotificationReceivers(Cluster cluster, ClusterConfigurationRequest configurationRequest) {
        cluster.setSlackReceivers(getSlackReceiversByIds(configurationRequest.slackReceiverIds()));
        cluster.setDiscordReceivers(getDiscordReceiversByIds(configurationRequest.discordReceiverIds()));
        cluster.setEmailReceivers(getEmailReceiversByIds(configurationRequest.emailReceiverIds()));
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
        return clusterRepository.findById(clusterId).map(cluster -> {
            Optional<Boolean> isRunning = isClusterRunning(clusterId);
            return isRunning
                    .map(running -> ClusterConfigurationDTO.ofCluster(cluster, running))
                    .orElse(ClusterConfigurationDTO.ofCluster(cluster, false));

        });
    }

    private Optional<Boolean> isClusterRunning(String clusterId) {
        return metadataService.getClusterById(clusterId).map(pl.pwr.zpi.metadata.dto.cluster.Cluster::running);
    }
}
