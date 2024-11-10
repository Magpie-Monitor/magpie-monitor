package pl.pwr.zpi.notifications.email.service;

public interface EmailMessagingService {
    void sendMessage(String receiver, String subject, String body, boolean isHtml);
}
