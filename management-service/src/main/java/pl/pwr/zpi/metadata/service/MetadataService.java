package pl.pwr.zpi.metadata.service;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.metadata.dto.application.AggregatedApplicationMetadata;
import pl.pwr.zpi.metadata.dto.application.ApplicationMetadata;
import pl.pwr.zpi.metadata.dto.cluster.AggregatedClusterMetadata;
import pl.pwr.zpi.metadata.dto.cluster.ClusterMetadata;
import pl.pwr.zpi.metadata.dto.node.AggregatedNodeMetadata;
import pl.pwr.zpi.metadata.dto.node.Node;
import pl.pwr.zpi.metadata.repository.AggregatedApplicationMetadataRepository;
import pl.pwr.zpi.metadata.repository.AggregatedClusterMetadataRepository;
import pl.pwr.zpi.metadata.repository.AggregatedNodeMetadataRepository;

import java.util.*;
import java.util.stream.Collectors;

@Slf4j
@Service
@RequiredArgsConstructor
public class MetadataService {

    private final AggregatedApplicationMetadataRepository applicationMetadataRepository;
    private final AggregatedNodeMetadataRepository nodeMetadataRepository;
    private final AggregatedClusterMetadataRepository clusterMetadataRepository;
    private final MetadataHistoryService metadataHistoryService;

    public List<ClusterMetadata> getAllClusters() {
        List<ClusterMetadata> activeClusters = getActiveClusters();

        Set<String> activeClustersIds = activeClusters.stream().map(ClusterMetadata::clusterId).collect(Collectors.toSet());

        Set<ClusterMetadata> inactiveClusters = metadataHistoryService.getClustersHistory().stream()
                .map(clusterHistory -> new ClusterMetadata(clusterHistory.id(), false))
                .filter(clusterMetadata -> !activeClustersIds.contains(clusterMetadata.clusterId()))
                .collect(Collectors.toSet());


        activeClusters.addAll(inactiveClusters);
        return activeClusters;
    }

    public List<ClusterMetadata> getActiveClusters() {
        return clusterMetadataRepository.findFirstByOrderByCollectedAtMsDesc()
                .map(AggregatedClusterMetadata::metadata)
                .orElse(Collections.emptyList());
    }

    public Optional<ClusterMetadata> getClusterById(String clusterId) {
        Optional<AggregatedClusterMetadata> metadata = clusterMetadataRepository.findFirstByMetadataClusterId(clusterId);
        if (metadata.isEmpty()) {
            return Optional.empty();
        }

        boolean running = !getActiveClusters().stream()
                .filter(clusterMetadata -> clusterMetadata.clusterId().equals(clusterId))
                .toList()
                .isEmpty();

        return Optional.of(new ClusterMetadata(clusterId, running));
    }

    public List<Node> getClusterNodes(String clusterId) {
        List<Node> activeNodes = nodeMetadataRepository.findFirstByClusterIdOrderByCollectedAtMs(clusterId)
                .map(aggregatedNodeMetadata -> aggregatedNodeMetadata.metadata().stream()
                        .map(nodeMetadata -> new Node(nodeMetadata.name(), true))
                        .collect(Collectors.toCollection(ArrayList::new))
                ).orElse(new ArrayList<>());

        Set<String> activeNodeIds = activeNodes.stream().map(Node::name).collect(Collectors.toSet());

        Set<Node> inactiveNodes = metadataHistoryService.getNodeHistory(clusterId).stream()
                .filter(node -> !activeNodeIds.contains(node.name()))
                .collect(Collectors.toSet());

        activeNodes.addAll(inactiveNodes);

        return activeNodes;
    }

    public List<ApplicationMetadata> getClusterApplications(String clusterId) {
        List<ApplicationMetadata> activeApplications = applicationMetadataRepository.findFirstByClusterIdOrderByCollectedAtMs(clusterId)
                .map(AggregatedApplicationMetadata::metadata)
                .orElse(Collections.emptyList());

        Set<String> activeApplicationIds = activeApplications.stream().map(ApplicationMetadata::name).collect(Collectors.toSet());

        Set<ApplicationMetadata> inactiveApplications = metadataHistoryService.getApplicationHistory(clusterId).stream()
                .filter(applicationMetadata -> !activeApplicationIds.contains(applicationMetadata.name()))
                .collect(Collectors.toSet());

        activeApplications.addAll(inactiveApplications);

        return activeApplications;
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
