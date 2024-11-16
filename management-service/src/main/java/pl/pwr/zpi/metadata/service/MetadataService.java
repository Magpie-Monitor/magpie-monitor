package pl.pwr.zpi.metadata.service;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.metadata.broker.dto.application.AggregatedApplicationMetadata;
import pl.pwr.zpi.metadata.broker.dto.cluster.AggregatedClusterMetadata;
import pl.pwr.zpi.metadata.broker.dto.node.AggregatedNodeMetadata;
import pl.pwr.zpi.metadata.dto.application.ApplicationMetadataDTO;
import pl.pwr.zpi.metadata.dto.cluster.ClusterMetadataDTO;
import pl.pwr.zpi.metadata.dto.node.NodeMetadataDTO;
import pl.pwr.zpi.metadata.repository.AggregatedApplicationMetadataRepository;
import pl.pwr.zpi.metadata.repository.AggregatedClusterMetadataRepository;
import pl.pwr.zpi.metadata.repository.AggregatedNodeMetadataRepository;

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

    public List<ClusterMetadataDTO> getAllClusters() {
        List<ClusterMetadataDTO> activeClusterMetadataDTOS = getActiveClusters();

        Set<ClusterMetadataDTO> inactiveClusterMetadataDTOS = filterInactiveClusters(activeClusterMetadataDTOS);
        activeClusterMetadataDTOS.addAll(inactiveClusterMetadataDTOS);

        return activeClusterMetadataDTOS;
    }

    public List<ClusterMetadataDTO> getActiveClusters() {
        return clusterMetadataRepository.findFirstByOrderByCollectedAtMsDesc()
                .map(metadata -> metadata.metadata()
                        .stream()
                        .map(cluster -> new ClusterMetadataDTO(cluster.clusterId(), true))
                        .collect(Collectors.toCollection(ArrayList::new))
                )
                .orElse(new ArrayList<>());
    }

    private Set<ClusterMetadataDTO> filterInactiveClusters(List<ClusterMetadataDTO> activeClusterMetadataDTOS) {
        Set<String> activeClusterIds = activeClusterMetadataDTOS.stream().map(ClusterMetadataDTO::clusterId).collect(Collectors.toSet());

        return metadataHistoryService.getClustersHistory().stream()
                .map(clusterHistory -> new ClusterMetadataDTO(clusterHistory.id(), false))
                .filter(clusterMetadataDTO -> !activeClusterIds.contains(clusterMetadataDTO.clusterId()))
                .collect(Collectors.toSet());
    }

    public Optional<ClusterMetadataDTO> getClusterById(String clusterId) {
        if (!clusterMetadataRepository.existsByMetadataClusterId(clusterId)) {
            return Optional.empty();
        }
        return Optional.of(new ClusterMetadataDTO(clusterId, getActiveClusters().contains(new ClusterMetadataDTO(clusterId, false))));
    }

    public List<NodeMetadataDTO> getClusterNodes(String clusterId) {
        List<NodeMetadataDTO> activeNodeMetadataDTOS = getActiveNodesForClusterId(clusterId);

        Set<NodeMetadataDTO> inactiveNodeMetadataDTOS = filterInactiveNodesForClusterId(clusterId, activeNodeMetadataDTOS);
        activeNodeMetadataDTOS.addAll(inactiveNodeMetadataDTOS);

        return activeNodeMetadataDTOS;
    }

    private List<NodeMetadataDTO> getActiveNodesForClusterId(String clusterId) {
        return nodeMetadataRepository.findFirstByClusterIdOrderByCollectedAtMsDesc(clusterId)
                .map(aggregatedNodeMetadata -> aggregatedNodeMetadata.metadata().stream()
                        .map(nodeMetadata -> new NodeMetadataDTO(nodeMetadata.name(), true))
                        .collect(Collectors.toCollection(ArrayList::new))
                ).orElse(new ArrayList<>());
    }

    private Set<NodeMetadataDTO> filterInactiveNodesForClusterId(String clusterId, List<NodeMetadataDTO> activeNodeMetadataDTOS) {
        Set<String> activeNodeIs = activeNodeMetadataDTOS.stream().map(NodeMetadataDTO::name).collect(Collectors.toSet());

        return metadataHistoryService.getNodeHistory(clusterId)
                .stream()
                .filter(nodeMetadataDTO -> !activeNodeIs.contains(nodeMetadataDTO.name()))
                .collect(Collectors.toSet());
    }

    public List<ApplicationMetadataDTO> getClusterApplications(String clusterId) {
        List<ApplicationMetadataDTO> activeApplicationMetadataDTOS = getActiveApplicationsForClusterId(clusterId);

        Set<ApplicationMetadataDTO> inactiveApplicationMetadataDTOS = filterInactiveApplicationsForClusterId(clusterId, activeApplicationMetadataDTOS);
        activeApplicationMetadataDTOS.addAll(inactiveApplicationMetadataDTOS);

        return activeApplicationMetadataDTOS;
    }

    private List<ApplicationMetadataDTO> getActiveApplicationsForClusterId(String clusterId) {
        return applicationMetadataRepository.findFirstByClusterIdOrderByCollectedAtMsDesc(clusterId)
                .map(metadata -> metadata.metadata().stream()
                        .map(application -> new ApplicationMetadataDTO(application.name(), application.kind(), true))
                        .collect(Collectors.toCollection(ArrayList::new))
                )
                .orElse(new ArrayList<>());
    }

    private Set<ApplicationMetadataDTO> filterInactiveApplicationsForClusterId(String clusterId, List<ApplicationMetadataDTO> activeApplicationMetadataDTOS) {
        Set<String> activeApplicationIds = activeApplicationMetadataDTOS.stream()
                .map(ApplicationMetadataDTO::name)
                .collect(Collectors.toSet());

        return metadataHistoryService.getApplicationHistory(clusterId).stream()
                .filter(applicationMetadataDTOMetadata -> !activeApplicationIds.contains(applicationMetadataDTOMetadata.name()))
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
