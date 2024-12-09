package pl.pwr.zpi.metadata.service;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.cluster.service.ClusterService;
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
import java.util.stream.Stream;

@Slf4j
@Service
@RequiredArgsConstructor
public class MetadataService {

    private final AggregatedApplicationMetadataRepository applicationMetadataRepository;
    private final AggregatedNodeMetadataRepository nodeMetadataRepository;
    private final AggregatedClusterMetadataRepository clusterMetadataRepository;
    private final MetadataHistoryService metadataHistoryService;

    public List<ClusterMetadataDTO> getAllClusters() {
        return Stream.concat(
                        getActiveClusters().stream(),
                        filterInactiveClusters(getActiveClusters()).stream())
                .collect(Collectors.toList());
    }

    public List<ClusterMetadataDTO> getActiveClusters() {
        return clusterMetadataRepository.findFirstByOrderByCollectedAtMsDesc()
                .map(metadata -> metadata.metadata()
                        .stream()
                        .map(cluster -> ClusterMetadataDTO.of(cluster.clusterId(), true))
                        .collect(Collectors.toCollection(ArrayList::new))
                )
                .orElse(new ArrayList<>());
    }

    private Set<ClusterMetadataDTO> filterInactiveClusters(List<ClusterMetadataDTO> activeClusterMetadataDTOS) {
        Set<String> activeClusterIds = activeClusterMetadataDTOS.stream().map(ClusterMetadataDTO::getClusterId).collect(Collectors.toSet());

        return metadataHistoryService.getClustersHistory().stream()
                .map(clusterHistory -> ClusterMetadataDTO.of(clusterHistory.id(), false))
                .filter(clusterMetadataDTO -> !activeClusterIds.contains(clusterMetadataDTO.getClusterId()))
                .collect(Collectors.toSet());
    }

    public Optional<ClusterMetadataDTO> getClusterById(String clusterId) {
        if (!clusterMetadataRepository.existsByMetadataClusterId(clusterId)) {
            return Optional.empty();
        }
        return Optional.of(ClusterMetadataDTO.of(clusterId, getActiveClusters().contains(ClusterMetadataDTO.of(clusterId, false))));
    }

    public List<NodeMetadataDTO> getClusterNodes(String clusterId) {
        return Stream.concat(
                        getActiveNodesForClusterId(clusterId).stream(),
                        filterInactiveNodesForClusterId(clusterId, getActiveNodesForClusterId(clusterId)).stream()
                )
                .collect(Collectors.toList());
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
