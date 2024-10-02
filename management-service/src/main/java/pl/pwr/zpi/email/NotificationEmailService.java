package pl.pwr.zpi.email;


public interface NotificationEmailService {

    void sendTestEmail(String receiverEmail);

    void sendNewReportNotification(String receiverEmail, String reportUrl);
}
