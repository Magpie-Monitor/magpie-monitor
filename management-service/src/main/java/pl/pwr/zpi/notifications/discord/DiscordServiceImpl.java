package pl.pwr.zpi.notifications.discord;

import kong.unirest.HttpResponse;
import kong.unirest.Unirest;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;

@Service
@RequiredArgsConstructor
@Slf4j
public class DiscordServiceImpl implements DiscordService {

    @Override
    public void sendMessage(String message, String webhookUrl) {
        try {
            HttpResponse<String> response = Unirest.post(webhookUrl)
                    .header("Content-Type", "application/json")
                    .body("{\"content\":\"" + message + "\"}")
                    .asString();

            if (!response.isSuccess()) {
                log.error("Failed to send message to Discord. Status: {}", response.getStatus());
            }
        } catch (Exception e) {
            log.error("Error sending message to Discord: {}", e.getMessage(), e);
        }
    }
}
