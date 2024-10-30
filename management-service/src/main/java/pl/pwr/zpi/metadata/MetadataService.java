package pl.pwr.zpi.metadata;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.metadata.dto.application.Application;
import pl.pwr.zpi.metadata.dto.cluster.Cluster;
import pl.pwr.zpi.metadata.dto.node.Node;
import pl.pwr.zpi.metadata.event.dto.application.AggregatedApplicationMetadata;
import pl.pwr.zpi.metadata.event.dto.cluster.AggregatedClusterMetadata;
import pl.pwr.zpi.metadata.event.dto.node.AggregatedNodeMetadata;
import pl.pwr.zpi.metadata.repository.AggregatedApplicationMetadataRepository;
import pl.pwr.zpi.metadata.repository.AggregatedClusterMetadataRepository;
import pl.pwr.zpi.metadata.repository.AggregatedNodeMetadataRepository;
import pl.pwr.zpi.metadata.service.MetadataHistoryService;

import java.util.ArrayList;
import java.util.List;
import java.util.Optional;
import java.util.Set;
import java.util.stream.Collectors;

@Slf4j
@Service
@RequiredArgsConstructor
public class MetadataService {

    private final AggregatedApplicationMetadataRepository applicationMetadataRepository;
    private final AggregatedNodeMetadataRepository nodeMetadataRepository;
    private final AggregatedClusterMetadataRepository clusterMetadataRepository;
    private final MetadataHistoryService metadataHistoryService;

    public List<Cluster> getAllClusters() {
        List<Cluster> activeClusters = getActiveClusters();

        Set<Cluster> inactiveClusters = filterInactiveClusters(activeClusters);
        activeClusters.addAll(inactiveClusters);

        return activeClusters;
    }

    public List<Cluster> getActiveClusters() {
        return clusterMetadataRepository.findFirstByOrderByCollectedAtMsDesc()
                .map(metadata -> metadata.metadata()
                        .stream()
                        .map(cluster -> new Cluster(cluster.clusterId(), true))
                        .collect(Collectors.toCollection(ArrayList::new))
                )
                .orElse(new ArrayList<>());
    }

    private Set<Cluster> filterInactiveClusters(List<Cluster> activeClusters) {
        Set<String> activeClusterIds = activeClusters.stream().map(Cluster::clusterId).collect(Collectors.toSet());

        return metadataHistoryService.getClustersHistory().stream()
                .map(clusterHistory -> new Cluster(clusterHistory.id(), false))
                .filter(cluster -> !activeClusterIds.contains(cluster.clusterId()))
                .collect(Collectors.toSet());
    }

    public Optional<Cluster> getClusterById(String clusterId) {
        if (!clusterMetadataRepository.existsByMetadataClusterId(clusterId)) {
            return Optional.empty();
        }
        return Optional.of(new Cluster(clusterId, getActiveClusters().contains(new Cluster(clusterId, false))));
    }

    public List<Node> getClusterNodes(String clusterId) {
        List<Node> activeNodes = getActiveNodesForClusterId(clusterId);

        Set<Node> inactiveNodes = filterInactiveNodesForClusterId(clusterId, activeNodes);
        activeNodes.addAll(inactiveNodes);

        return activeNodes;
    }

    private List<Node> getActiveNodesForClusterId(String clusterId) {
        return nodeMetadataRepository.findFirstByClusterIdOrderByCollectedAtMsDesc(clusterId)
                .map(aggregatedNodeMetadata -> aggregatedNodeMetadata.metadata().stream()
                        .map(nodeMetadata -> new Node(nodeMetadata.name(), true))
                        .collect(Collectors.toCollection(ArrayList::new))
                ).orElse(new ArrayList<>());
    }

    private Set<Node> filterInactiveNodesForClusterId(String clusterId, List<Node> activeNodes) {
        Set<String> activeNodeIs = activeNodes.stream().map(Node::name).collect(Collectors.toSet());

        return metadataHistoryService.getNodeHistory(clusterId)
                .stream()
                .filter(node -> !activeNodeIs.contains(node.name()))
                .collect(Collectors.toSet());
    }

    public List<Application> getClusterApplications(String clusterId) {
        List<Application> activeApplications = getActiveApplicationsForClusterId(clusterId);

        Set<Application> inactiveApplications = filterInactiveApplicationsForClusterId(clusterId, activeApplications);
        activeApplications.addAll(inactiveApplications);

        return activeApplications;
    }

    private List<Application> getActiveApplicationsForClusterId(String clusterId) {
        return applicationMetadataRepository.findFirstByClusterIdOrderByCollectedAtMsDesc(clusterId)
                .map(metadata -> metadata.metadata().stream()
                        .map(application -> new Application(application.name(), application.kind(), true))
                        .collect(Collectors.toCollection(ArrayList::new))
                )
                .orElse(new ArrayList<>());
    }

    private Set<Application> filterInactiveApplicationsForClusterId(String clusterId, List<Application> activeApplications) {
        Set<String> activeApplicationIds = activeApplications.stream()
                .map(Application::name)
                .collect(Collectors.toSet());

        return metadataHistoryService.getApplicationHistory(clusterId).stream()
                .filter(applicationMetadata -> !activeApplicationIds.contains(applicationMetadata.name()))
                .collect(Collectors.toSet());
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
