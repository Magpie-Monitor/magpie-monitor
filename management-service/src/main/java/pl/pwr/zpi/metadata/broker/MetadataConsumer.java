package pl.pwr.zpi.metadata.broker;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.stereotype.Component;
import pl.pwr.zpi.metadata.service.MetadataHistoryService;
import pl.pwr.zpi.metadata.service.MetadataService;
import pl.pwr.zpi.metadata.broker.dto.application.ApplicationMetadataUpdated;
import pl.pwr.zpi.metadata.broker.dto.cluster.ClusterMetadataUpdated;
import pl.pwr.zpi.metadata.broker.dto.node.NodeMetadataUpdated;
import pl.pwr.zpi.utils.mapper.JsonMapper;

@Slf4j
@RequiredArgsConstructor
@Component
public class MetadataConsumer {

    private final MetadataService metadataService;
    private final MetadataHistoryService metadataHistoryService;
    private final JsonMapper mapper;

    @KafkaListener(topics = "${kafka.cluster.metadata.application.topic}")
    public void listenForApplicationMetadataStateUpdate(String message) {
        log.info("Application metadata updated: {}", message);
        ApplicationMetadataUpdated metadata = mapper.fromJson(message, ApplicationMetadataUpdated.class);
        metadataService.saveApplicationMetadata(metadata.metadata());
        metadataHistoryService.updateApplicationHistory(metadata.clusterId(), metadata.applicationMetadata());
    }

    @KafkaListener(topics = "${kafka.cluster.metadata.node.topic}")
    public void listenForNodeMetadataStateUpdate(String message) {
        log.info("Node metadata updated {}", message);
        NodeMetadataUpdated metadata = mapper.fromJson(message, NodeMetadataUpdated.class);
        metadataService.saveNodeMetadata(metadata.metadata());
        metadataHistoryService.updateNodeHistory(metadata.clusterId(), metadata.nodeMetadata());
    }

    @KafkaListener(topics = "${kafka.cluster.metadata.cluster.topic}")
    public void listenForClusterMetadataStateUpdate(String message) {
        log.info("Cluster metadata updated: {}", message);
        ClusterMetadataUpdated metadata = mapper.fromJson(message, ClusterMetadataUpdated.class);
        metadataService.saveClusterMetadata(metadata.metadata());
        metadataHistoryService.updateClustersHistory(metadata.clusterMetadata());
    }
}
