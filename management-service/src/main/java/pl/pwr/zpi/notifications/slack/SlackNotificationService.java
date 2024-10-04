package pl.pwr.zpi.notifications.slack;

public interface SlackNotificationService {
    void sendTestMessage(String webhookUrl);
    void sendMessageAboutNewReport(String webhookUrl, String message);
}
