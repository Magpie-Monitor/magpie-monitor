package pl.pwr.zpi.notifications.discord.service;

import kong.unirest.HttpResponse;
import kong.unirest.Unirest;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;

@Service
@RequiredArgsConstructor
@Slf4j
public class DiscordMessagingService {

    public void sendMessage(String message, String webhookUrl) {
        try {
            HttpResponse<String> response = Unirest.post(webhookUrl)
                    .header("Content-Type", "application/json")
                    .body("{\"content\":\"" + message + "\"}")
                    .asString();

            if (!response.isSuccess()) {
                log.error("Failed to send message to Discord. Status: {}", response.getStatus());
                throw new Exception("Failed to send message to Discord. Status: " + response.getStatus());
            }
        } catch (Exception e) {
            log.error("Error sending message to Discord: {}", e.getMessage(), e);
            throw new RuntimeException("Error sending message to Discord: " + e.getMessage());
        }
    }
}
