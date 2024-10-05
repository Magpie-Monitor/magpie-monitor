package pl.pwr.zpi.notifications.email;

public interface EmailService {
    void sendMessage(String receiver, String subject, String body, boolean isHtml);
}
