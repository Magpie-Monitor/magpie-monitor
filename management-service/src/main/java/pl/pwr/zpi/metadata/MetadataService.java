package pl.pwr.zpi.metadata;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.metadata.dto.application.AggregatedApplicationMetadata;
import pl.pwr.zpi.metadata.dto.application.ApplicationMetadata;
import pl.pwr.zpi.metadata.dto.cluster.AggregatedClusterMetadata;
import pl.pwr.zpi.metadata.dto.cluster.ClusterMetadata;
import pl.pwr.zpi.metadata.dto.node.AggregatedNodeMetadata;
import pl.pwr.zpi.metadata.dto.node.Node;
import pl.pwr.zpi.metadata.repository.*;

import java.util.*;
import java.util.stream.Collectors;

@Slf4j
@Service
@RequiredArgsConstructor
public class MetadataService {

    private final AggregatedApplicationMetadataRepository applicationMetadataRepository;
    private final AggregatedNodeMetadataRepository nodeMetadataRepository;
    private final AggregatedClusterMetadataRepository clusterMetadataRepository;

    public List<ClusterMetadata> getAllClusters() {
        List<ClusterMetadataProjection> metadata = clusterMetadataRepository.findAllByOrderByCollectedAtMsDesc();

        if (metadata.isEmpty()) {
            return List.of();
        }

        Set<String> activeClusterIds = metadata.getFirst().getMetadata().stream()
                .map(ClusterMetadata::name)
                .collect(Collectors.toSet());

        Set<String> clusterIds = metadata.stream()
                .map(ClusterMetadataProjection::getMetadata)
                .flatMap(Collection::stream)
                .map(ClusterMetadata::name).
                collect(Collectors.toSet());

        return clusterIds.stream()
                .map(id -> new ClusterMetadata(id, activeClusterIds.contains(id)))
                .toList();
    }

    public List<ClusterMetadata> getActiveClusters() {
        return clusterMetadataRepository.findFirstByOrderByCollectedAtMsDesc()
                .map(AggregatedClusterMetadata::metadata)
                .orElse(Collections.emptyList());
    }

    public Optional<ClusterMetadata> getClusterById(String clusterId) {
        Optional<AggregatedClusterMetadata> metadata = clusterMetadataRepository.findFirstByMetadataName(clusterId);
        if (metadata.isEmpty()) {
            return Optional.empty();
        }

        Optional<AggregatedClusterMetadata> activeClustersMetadata = clusterMetadataRepository.findFirstByOrderByCollectedAtMsDesc();
        if (activeClustersMetadata.isEmpty()) {
            return Optional.empty();
        }

        boolean running = activeClustersMetadata.get().metadata().stream()
                .map(ClusterMetadata::name)
                .collect(Collectors.toSet())
                .contains(clusterId);

        return Optional.of(new ClusterMetadata(clusterId, running));
    }

    public List<Node> getClusterActiveNodes(String clusterId) {
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
