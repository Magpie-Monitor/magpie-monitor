package pl.pwr.zpi.notifications.slack;

public interface SlackService {
    void sendMessage(String message, String webhookUrl);
}
