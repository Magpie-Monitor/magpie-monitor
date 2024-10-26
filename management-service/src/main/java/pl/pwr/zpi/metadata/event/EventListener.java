package pl.pwr.zpi.metadata.event;

import com.fasterxml.jackson.core.JsonParser;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.ObjectMapper;
import lombok.extern.slf4j.Slf4j;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.metadata.MetadataService;
import pl.pwr.zpi.metadata.event.dto.ApplicationMetadataUpdated;
import pl.pwr.zpi.metadata.event.dto.ClusterMetadataUpdated;
import pl.pwr.zpi.metadata.event.dto.NodeMetadataUpdated;

@Slf4j
@Service
public class EventListener {

    private final MetadataService metadataService;
    private final ObjectMapper mapper;

    public EventListener(MetadataService metadataService) {
        this.metadataService = metadataService;
        this.mapper = new ObjectMapper();
        this.mapper.configure(JsonParser.Feature.INCLUDE_SOURCE_IN_LOCATION, true);
        this.mapper.configure(DeserializationFeature.ACCEPT_SINGLE_VALUE_AS_ARRAY, true);
    }

    // TODO - add wrapper around obj mapper
    @KafkaListener(topics = "ApplicationMetadataUpdated")
    public void listenForApplicationMetadataStateUpdate(String message) {
        try {
            ApplicationMetadataUpdated metadata = mapper.readValue(message, ApplicationMetadataUpdated.class);
            log.info("Application metadata updated: {}", metadata);
            metadataService.saveApplicationMetadata(metadata.metadata());
        } catch (JsonProcessingException e) {
            log.error("Error parsing ApplicationMetadataUpdated event {}", e.getMessage());
            throw new RuntimeException(e);
        }
    }

    @KafkaListener(topics = "NodeMetadataUpdated")
    public void listenForNodeMetadataStateUpdate(String message) {
        try {
            NodeMetadataUpdated metadata = mapper.readValue(message, NodeMetadataUpdated.class);
            log.info("Node metadata updated {}", metadata);
            metadataService.saveNodeMetadata(metadata.metadata());
        } catch (JsonProcessingException e) {
            log.error("Error parsing NodeMetadataUpdated event {}", e.getMessage());
            throw new RuntimeException(e);
        }
    }

    @KafkaListener(topics = "ClusterMetadataUpdated")
    public void listenForClusterMetadataStateUpdate(String message) {
        try {
            ClusterMetadataUpdated metadata = mapper.readValue(message, ClusterMetadataUpdated.class);
            log.info("Cluster metadata updated: {}", metadata);
            metadataService.saveClusterMetadata(metadata.metadata());
        } catch (JsonProcessingException e) {
            log.error("Error parsing NodeMetadataUpdated event {}", e.getMessage());
            throw new RuntimeException(e);
        }
    }
}
