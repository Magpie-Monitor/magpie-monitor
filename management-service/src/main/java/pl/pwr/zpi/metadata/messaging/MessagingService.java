package pl.pwr.zpi.metadata.messaging;

import com.fasterxml.jackson.core.JsonParser;
import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.ObjectMapper;
import lombok.RequiredArgsConstructor;
import lombok.SneakyThrows;
import lombok.extern.slf4j.Slf4j;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.metadata.messaging.event.ApplicationMetadataUpdated;
import pl.pwr.zpi.metadata.messaging.event.NodeMetadataUpdated;

@Slf4j
@Service
public class MessagingService {

    private final KafkaTemplate<String, String> kafka;
    private final ObjectMapper mapper;

    public MessagingService(KafkaTemplate<String, String> kafka) {
        mapper = new ObjectMapper();
        this.kafka = kafka;
        this.mapper.configure(JsonParser.Feature.INCLUDE_SOURCE_IN_LOCATION, true);
        this.mapper.configure(DeserializationFeature.ACCEPT_SINGLE_VALUE_AS_ARRAY, true);
    }

    @SneakyThrows
    @KafkaListener(topics = "ApplicationMetadataUpdated")
    public void processApplicationMetadataStateChange(String message) {
        ApplicationMetadataUpdated metadata = mapper.readValue(message, ApplicationMetadataUpdated.class);
        log.info("Application metadata updated: {}", metadata);
    }

    @SneakyThrows
    @KafkaListener(topics = "NodeMetadataUpdated")
    public void processNodeMetadataStateChange(String message) {
        NodeMetadataUpdated metadata = mapper.readValue(message, NodeMetadataUpdated.class);
        log.info("Node metadata updated: {}", metadata);
    }
}
