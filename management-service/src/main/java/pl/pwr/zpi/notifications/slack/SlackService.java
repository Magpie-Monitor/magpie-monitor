package pl.pwr.zpi.notifications.slack;

import lombok.SneakyThrows;

public interface SlackService {
    @SneakyThrows
    void sendMessage(String message, String webhookUrl);
}
