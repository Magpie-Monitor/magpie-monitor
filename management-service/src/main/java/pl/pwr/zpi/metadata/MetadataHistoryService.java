package pl.pwr.zpi.metadata;

import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.metadata.dto.application.ApplicationMetadata;
import pl.pwr.zpi.metadata.dto.cluster.ClusterHistory;
import pl.pwr.zpi.metadata.dto.cluster.ClusterMetadata;
import pl.pwr.zpi.metadata.dto.node.Node;
import pl.pwr.zpi.metadata.dto.node.NodeMetadata;
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

    public Set<Node> getNodeHistory(String clusterId) {
        ClusterHistory history = clusterHistoryRepository.findById(clusterId).orElse(new ClusterHistory(clusterId, new HashSet<>(), new HashSet<>()));
        return history.nodes();
    }

    public Set<ApplicationMetadata> getApplicationHistory(String clusterId) {
        ClusterHistory history = clusterHistoryRepository.findById(clusterId).orElse(new ClusterHistory(clusterId, new HashSet<>(), new HashSet<>()));
        return history.applications();
    }

    public void updateClustersHistory(List<ClusterMetadata> metadata) {
        for (ClusterMetadata m : metadata) {
            if (!clusterHistoryRepository.existsById(m.name())) {
                clusterHistoryRepository.save(new ClusterHistory(m.name(), new HashSet<>(), new HashSet<>()));
            }
        }
    }

    public void updateNodeHistory(String clusterId, List<NodeMetadata> metadata) {
        ClusterHistory history = clusterHistoryRepository.findById(clusterId)
                .orElse(new ClusterHistory(clusterId, new HashSet<>(), new HashSet<>()));

        metadata.forEach(m -> history.nodes().add(new Node(m.name(), false)));

        clusterHistoryRepository.save(history);
    }

    public void updateApplicationHistory(String clusterId, List<ApplicationMetadata> metadata) {
        ClusterHistory history = clusterHistoryRepository.findById(clusterId)
                .orElse(new ClusterHistory(clusterId, new HashSet<>(), new HashSet<>()));

        metadata.forEach(m -> history.applications().add(new ApplicationMetadata(m.name(), m.kind(), false)));

        clusterHistoryRepository.save(history);
    }

}
