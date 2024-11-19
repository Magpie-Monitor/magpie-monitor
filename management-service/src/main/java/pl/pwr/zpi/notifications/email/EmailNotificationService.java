package pl.pwr.zpi.notifications.email;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.notifications.ReportNotifier;
import pl.pwr.zpi.notifications.email.entity.EmailReceiver;
import pl.pwr.zpi.notifications.email.internalization.service.LocalizedMessageService;
import pl.pwr.zpi.notifications.email.service.EmailMessagingService;
import pl.pwr.zpi.notifications.email.service.EmailReceiverService;
import pl.pwr.zpi.notifications.email.utils.EmailUtils;

@Slf4j
@Service("emailNotificationService")
@RequiredArgsConstructor
public class EmailNotificationService implements ReportNotifier {

    private final LocalizedMessageService localizedTestMailServiceImpl;
    private final LocalizedMessageService localizedNewReportMailServiceImpl;
    private final EmailMessagingService emailService;
    private final EmailUtils emailUtils;
    private final EmailReceiverService emailReceiverService;

    @Value("${magpie.monitor.client.base.url}")
    private String MAGPIE_MONITOR_CLIENT_BASE_URL;

    public void sendTestEmail(Long receiverEmailId) {
        var receiver = emailReceiverService.getEmailReceiver(receiverEmailId);
        log.info("Sending test email to: {}", receiver.getReceiverEmail());
        sendTestEmail(receiver.getReceiverEmail());
    }

    private void sendTestEmail(String receiverEmail) {
        String title = localizedTestMailServiceImpl.getMessage("test.title",
                localizedTestMailServiceImpl.getLanguageFromContextOrDefault());

        emailService.sendMessage(receiverEmail,
                title,
                emailUtils.createTestEmailTemplate(),
                true);
    }

    @Override
    public void notifyOnReportGenerated(Long receiverId, String reportId) {
        EmailReceiver emailReceiver = emailReceiverService.getEmailReceiver(receiverId);

        String title = localizedNewReportMailServiceImpl.getMessage("new-report.title",
                localizedNewReportMailServiceImpl.getLanguageFromContextOrDefault());

        emailService.sendMessage(emailReceiver.getReceiverEmail(),
                title,
                emailUtils.createNewReportEmailTemplate(String.format("%s/reports/%s", MAGPIE_MONITOR_CLIENT_BASE_URL, reportId)),
                true);
    }

    @Override
    public void notifyOnReportGenerationFailed(Long receiverId, String clusterId) {

    }
}
