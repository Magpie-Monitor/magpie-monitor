package pl.pwr.zpi.metadata.messaging;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Service;

@Slf4j
@Service
@RequiredArgsConstructor
public class MessagingService {

    private final KafkaTemplate<String, String> kafka;

    @KafkaListener(topics = "applications")
    public void listen(String message) {
        log.info("received message: {}", message);
    }

    @Scheduled(fixedDelay = 1000)
    public void publish() {
        log.info("published message");
        kafka.send("applications", "hello");
    }
}
