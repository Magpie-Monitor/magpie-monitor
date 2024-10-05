package pl.pwr.zpi.notifications.discord;

import lombok.SneakyThrows;

public interface DiscordService {
    @SneakyThrows
    void sendMessage(String message, String webhookUrl);
}
