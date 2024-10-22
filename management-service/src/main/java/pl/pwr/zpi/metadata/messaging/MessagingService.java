package pl.pwr.zpi.metadata.messaging;

import com.fasterxml.jackson.core.JsonParser;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.ObjectMapper;
import lombok.extern.slf4j.Slf4j;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.metadata.MetadataService;
import pl.pwr.zpi.metadata.messaging.event.ApplicationMetadataUpdated;
import pl.pwr.zpi.metadata.messaging.event.NodeMetadataUpdated;

@Slf4j
@Service
public class MessagingService {

    private final MetadataService metadataService;
    private final ObjectMapper mapper;

    public MessagingService(MetadataService metadataService) {
        this.metadataService = metadataService;
        this.mapper = new ObjectMapper();
        this.mapper.configure(JsonParser.Feature.INCLUDE_SOURCE_IN_LOCATION, true);
        this.mapper.configure(DeserializationFeature.ACCEPT_SINGLE_VALUE_AS_ARRAY, true);
    }

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
            log.info("Node metadata updated: {}", metadata);
            metadataService.saveNodeMetadata(metadata.metadata());
        } catch (JsonProcessingException e) {
            log.error("Error parsing NodeMetadataUpdated event {}", e.getMessage());
            throw new RuntimeException(e);
        }
    }
}
