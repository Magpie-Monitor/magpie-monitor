package pl.pwr.zpi.notifications.discord;

public interface DiscordNotificationService {
    void sendTestMessage(String webhookUrl);
    void sendMessageAboutNewReport(String webhookUrl, String message);
}
