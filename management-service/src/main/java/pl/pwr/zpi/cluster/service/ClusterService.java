package pl.pwr.zpi.cluster.service;

import com.mongodb.connection.ClusterId;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.cluster.dto.ClusterConfigurationDTO;
import pl.pwr.zpi.cluster.dto.ClusterConfigurationRequest;
import pl.pwr.zpi.cluster.dto.ClusterIdResponse;
import pl.pwr.zpi.cluster.entity.Cluster;
import pl.pwr.zpi.cluster.repository.ClusterRepository;
import pl.pwr.zpi.metadata.service.MetadataService;
import pl.pwr.zpi.notifications.ReceiverService;
import pl.pwr.zpi.notifications.slack.entity.SlackReceiver;

import java.util.List;
import java.util.Optional;

@Service
@RequiredArgsConstructor
public class ClusterService {

    // TODO - decouple notification from receivers
    private final ClusterRepository clusterRepository;
    private final ReceiverService receiverService;
    private final MetadataService metadataService;

    // TODO - check if sent receivers exist
    public ClusterIdResponse updateCluster(ClusterConfigurationRequest configurationRequest) {
//        validateReceiverIds(configurationRequest);

        Cluster cluster = Cluster.ofClusterConfigurationRequest(configurationRequest);
        cluster.setSlackReceivers(getSlackReceiversByIds(configurationRequest.slackReceiverIds()));
//        cluster.setDiscordReceivers();
//        cluster.setMailReceivers();

        clusterRepository.save(cluster);

        return new ClusterIdResponse(cluster.getId());
    }

    private List<SlackReceiver> getSlackReceiversByIds(List<Long> receiverIds) {
        return receiverIds.stream()
                .map(receiverService::getReceiverById)
                .toList();
    }

//    private void setClusterReceivers(ClusterConfigurationRequest configurationRequest, Cluster cluster) {
////        cluster.setSlackReceivers();
//    }

//    private void validateReceiverIds(ClusterConfigurationRequest configurationRequest) {
//        validateSlackReceiverIds(configurationRequest.slackReceiverIds());
//    }
//
//    private void validateSlackReceiverIds(List<Long> receiverIds) {
//        receiverIds.forEach(id -> {
//            if (receiverService.slackReceiverExists(id)) {
//                throw new RuntimeException("Slack receiver with an id " + id + " does not exist");
//            }
//        });
//    }

    public Optional<ClusterConfigurationDTO> getClusterById(String clusterId) {
        return clusterRepository.findById(clusterId).map(cluster -> {
            Optional<Boolean> isRunning = isClusterRunning(clusterId);
            return isRunning
                    .map(running -> ClusterConfigurationDTO.ofCluster(cluster, running))
                    .orElse(ClusterConfigurationDTO.ofCluster(cluster, false));

        });
    }

    public Optional<Boolean> isClusterRunning(String clusterId) {
        return metadataService.getClusterById(clusterId).map(pl.pwr.zpi.metadata.dto.cluster.Cluster::running);
    }
}
