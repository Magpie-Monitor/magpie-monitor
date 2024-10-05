package pl.pwr.zpi.notifications.email;

import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.notifications.email.iternalization.service.LocalizedMessageService;

@RequiredArgsConstructor
@Service
public class EmailNotificationServiceImpl implements EmailNotificationService {

    private final LocalizedMessageService localizedTestMailServiceImpl;
    private final LocalizedMessageService localizedNewReportMailServiceImpl;
    private final EmailService emailService;
    private final EmailUtils emailUtils;

    @Override
    public void sendTestEmail(String receiverEmail) {
        String title = localizedTestMailServiceImpl.getMessage("test.title",
                localizedTestMailServiceImpl.getLanguageFromContextOrDefault());

        emailService.sendMessage(receiverEmail,
                title,
                emailUtils.createTestEmailTemplate(),
                true);
    }

    @Override
    public void sendNewReportNotification(String receiverEmail, String reportUrl) {
        String title = localizedNewReportMailServiceImpl.getMessage("new-report.title",
                localizedNewReportMailServiceImpl.getLanguageFromContextOrDefault());

        emailService.sendMessage(receiverEmail,
                title,
                emailUtils.createNewReportEmailTemplate(reportUrl),
                true);
    }
}
