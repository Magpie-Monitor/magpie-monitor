package pl.pwr.zpi.notifications.slack.service;

public interface SlackMessagingService {
    void sendMessage(String message, String webhookUrl);
}
