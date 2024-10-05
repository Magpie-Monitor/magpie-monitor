package pl.pwr.zpi.notifications.email;


public interface EmailNotificationService {

    void sendTestEmail(String receiverEmail);

    void sendNewReportNotification(String receiverEmail, String reportUrl);
}
