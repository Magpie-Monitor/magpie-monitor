package pl.pwr.zpi.notifications.slack.service;

import com.slack.api.Slack;
import com.slack.api.webhook.Payload;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;

import java.io.IOException;

@Service
@RequiredArgsConstructor
@Slf4j
public class SlackMessagingServiceImpl implements SlackMessagingService {

    @Override
    public void sendMessage(String message, String webhookUrl) {
        Payload payload = Payload.builder()
                .text(message)
                .build();
        try {
            Slack.getInstance().send(webhookUrl, payload);
        } catch (IOException e) {
            log.error("Error sending message to Slack: {}", e.getMessage(), e);
            throw new RuntimeException("Error sending message to Slack: " + e.getMessage());
        }
    }
}
