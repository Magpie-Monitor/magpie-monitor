package pl.pwr.zpi.notifications.slack.service;

public interface SlackNotificationService {
    void sendTestMessage(String webhookUrl);
    void notifyOnReportCreated(Long receiverId, String message);
}
