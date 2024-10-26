package pl.pwr.zpi.metadata;

import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.metadata.dto.Cluster;
import pl.pwr.zpi.metadata.dto.application.AggregatedApplicationMetadata;
import pl.pwr.zpi.metadata.dto.application.ApplicationMetadata;
import pl.pwr.zpi.metadata.dto.cluster.AggregatedClusterMetadata;
import pl.pwr.zpi.metadata.dto.cluster.ClusterMetadata;
import pl.pwr.zpi.metadata.dto.node.AggregatedNodeMetadata;
import pl.pwr.zpi.metadata.dto.node.Node;
import pl.pwr.zpi.metadata.repository.AggregatedApplicationMetadataRepository;
import pl.pwr.zpi.metadata.repository.AggregatedClusterMetadataRepository;
import pl.pwr.zpi.metadata.repository.AggregatedNodeMetadataRepository;

import java.util.Collections;
import java.util.List;
import java.util.Optional;

@Service
@RequiredArgsConstructor
public class MetadataService {

    private final AggregatedApplicationMetadataRepository applicationMetadataRepository;
    private final AggregatedNodeMetadataRepository nodeMetadataRepository;
    private final AggregatedClusterMetadataRepository clusterMetadataRepository;

    public List<ClusterMetadata> getClusters() {
        return clusterMetadataRepository.findFirstByOrderByCollectedAtMsDesc()
                .map(AggregatedClusterMetadata::metadata)
                .orElse(Collections.emptyList());
    }

    public Optional<ClusterMetadata> getClusterById(String clusterId) {
        return getClusters().stream()
                .filter(clusterMetadata -> clusterMetadata.name().equals(clusterId))
                .findFirst();
    }

    public List<Node> getClusterNodes(String clusterId) {
        return nodeMetadataRepository.findFirstByClusterIdOrderByCollectedAtMs(clusterId)
                .map(aggregatedNodeMetadata -> aggregatedNodeMetadata.metadata().stream()
                        .map(Node::of)
                        .toList()
                ).orElse(Collections.emptyList());
    }

    public List<ApplicationMetadata> getClusterApplications(String clusterId) {
        return applicationMetadataRepository.findFirstByClusterIdOrderByCollectedAtMs(clusterId)
                .map(AggregatedApplicationMetadata::metadata)
                .orElse(Collections.emptyList());
    }

    public void saveClusterMetadata(AggregatedClusterMetadata clusterMetadata) {
        clusterMetadataRepository.save(clusterMetadata);
    }

    public void saveApplicationMetadata(AggregatedApplicationMetadata applicationMetadata) {
        applicationMetadataRepository.save(applicationMetadata);
    }

    public void saveNodeMetadata(AggregatedNodeMetadata nodeMetadata) {
        nodeMetadataRepository.save(nodeMetadata);
    }
}
