package pl.pwr.zpi.notifications.slack;

import com.slack.api.Slack;
import com.slack.api.model.Attachment;
import com.slack.api.webhook.Payload;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;

import java.io.IOException;

@Service
@RequiredArgsConstructor
@Slf4j
public class SlackServiceImpl implements SlackService {
    @Override
    public void sendMessage(String message, String webhookUrl) {
        Payload payload = Payload.builder()
                .text(message)
                .build();
        try {
            Slack.getInstance().send(webhookUrl, payload);
        } catch (IOException e) {
            log.error("Error while sending!");

        }
    }
}
