package pl.pwr.zpi.metadata.messaging;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.stereotype.Service;

@Slf4j
@Service
@RequiredArgsConstructor
public class MessagingService {

    private final KafkaTemplate<String, String> kafka;

    @KafkaListener(topics = "ApplicationMetadataUpdated")
    public void processApplicationMetadataStateChange(String message) {
        log.info("Received updated application metadata state: {}", message);
    }

    @KafkaListener(topics = "NodeMetadataUpdated")
    public void processNodeMetadataStateChange(String message) {
        log.info("Received updated node metadata state: {}", message);
    }
}
