package pl.pwr.zpi.metadata.service;

import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.metadata.dto.application.ApplicationMetadataDTO;
import pl.pwr.zpi.metadata.dto.node.NodeMetadataDTO;
import pl.pwr.zpi.metadata.entity.ClusterHistory;
import pl.pwr.zpi.metadata.broker.dto.application.ApplicationMetadata;
import pl.pwr.zpi.metadata.broker.dto.cluster.ClusterMetadata;
import pl.pwr.zpi.metadata.broker.dto.node.NodeMetadata;
import pl.pwr.zpi.metadata.repository.ClusterHistoryRepository;

import java.util.HashSet;
import java.util.List;
import java.util.Set;

@Service
@RequiredArgsConstructor
public class MetadataHistoryService {

    private final ClusterHistoryRepository clusterHistoryRepository;

    public List<ClusterHistory> getClustersHistory() {
        return clusterHistoryRepository.findAll();
    }

    public Set<NodeMetadataDTO> getNodeHistory(String clusterId) {
        ClusterHistory history = clusterHistoryRepository.findById(clusterId)
                .orElse(new ClusterHistory(clusterId, new HashSet<>(), new HashSet<>()));
        return history.nodeMetadataDTOS();
    }

    public Set<ApplicationMetadataDTO> getApplicationHistory(String clusterId) {
        ClusterHistory history = clusterHistoryRepository.findById(clusterId)
                .orElse(new ClusterHistory(clusterId, new HashSet<>(), new HashSet<>()));
        return history.applicationMetadataDTOS();
    }

    public void updateClustersHistory(List<ClusterMetadata> metadata) {
        metadata.stream()
                .filter(clusterMetadata -> !clusterHistoryRepository.existsById(clusterMetadata.clusterId()))
                .forEach(clusterMetadata -> clusterHistoryRepository.save(ClusterHistory.of(clusterMetadata)));
    }

    public void updateNodeHistory(String clusterId, List<NodeMetadata> metadata) {
        ClusterHistory history = clusterHistoryRepository.findById(clusterId)
                .orElse(new ClusterHistory(clusterId, new HashSet<>(), new HashSet<>()));

        metadata.forEach(m -> history.nodeMetadataDTOS().add(new NodeMetadataDTO(m.name(), false)));

        clusterHistoryRepository.save(history);
    }

    public void updateApplicationHistory(String clusterId, List<ApplicationMetadata> metadata) {
        ClusterHistory history = clusterHistoryRepository.findById(clusterId)
                .orElse(new ClusterHistory(clusterId, new HashSet<>(), new HashSet<>()));

        metadata.forEach(m -> history.applicationMetadataDTOS().add(new ApplicationMetadataDTO(m.name(), m.kind(), false)));

        clusterHistoryRepository.save(history);
    }

}
