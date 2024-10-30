package pl.pwr.zpi.metadata.event;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.stereotype.Component;
import pl.pwr.zpi.metadata.service.MetadataHistoryService;
import pl.pwr.zpi.metadata.MetadataService;
import pl.pwr.zpi.metadata.event.dto.application.ApplicationMetadataUpdated;
import pl.pwr.zpi.metadata.event.dto.cluster.ClusterMetadataUpdated;
import pl.pwr.zpi.metadata.event.dto.node.NodeMetadataUpdated;
import pl.pwr.zpi.utils.mapper.JsonMapper;

@Slf4j
@RequiredArgsConstructor
@Component
public class MetadataEventListener {

    private final MetadataService metadataService;
    private final MetadataHistoryService metadataHistoryService;
    private final JsonMapper mapper;

    @KafkaListener(topics = "ApplicationMetadataUpdated")
    public void listenForApplicationMetadataStateUpdate(String message) {
        ApplicationMetadataUpdated metadata = mapper.fromJson(message, ApplicationMetadataUpdated.class);
        log.info("Application metadata updated: {}", metadata);
        metadataService.saveApplicationMetadata(metadata.metadata());
        metadataHistoryService.updateApplicationHistory(metadata.clusterId(), metadata.applicationMetadata());
    }

    @KafkaListener(topics = "NodeMetadataUpdated")
    public void listenForNodeMetadataStateUpdate(String message) {
        NodeMetadataUpdated metadata = mapper.fromJson(message, NodeMetadataUpdated.class);
        log.info("Node metadata updated {}", metadata);
        metadataService.saveNodeMetadata(metadata.metadata());
        metadataHistoryService.updateNodeHistory(metadata.clusterId(), metadata.nodeMetadata());
    }

    @KafkaListener(topics = "ClusterMetadataUpdated")
    public void listenForClusterMetadataStateUpdate(String message) {
        ClusterMetadataUpdated metadata = mapper.fromJson(message, ClusterMetadataUpdated.class);
        log.info("Cluster metadata updated: {}", metadata);
        metadataService.saveClusterMetadata(metadata.metadata());
        metadataHistoryService.updateClustersHistory(metadata.clusterMetadata());
    }
}
