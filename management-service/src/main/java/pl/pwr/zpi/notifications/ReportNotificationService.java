package pl.pwr.zpi.notifications;

import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
public class ReportNotificationService {

    private final ReportNotifier slackNotificationService;
    private final ReportNotifier emailNotificationService;
    private final ReportNotifier discordNotificationService;

    public ReportNotificationService(
            @Qualifier("slackNotificationService")
            ReportNotifier slackNotificationService,
            @Qualifier("emailNotificationService")
            ReportNotifier emailNotificationService,
            @Qualifier("discordNotificationService")
            ReportNotifier discordNotificationService
    ) {
        this.slackNotificationService = slackNotificationService;
        this.emailNotificationService = emailNotificationService;
        this.discordNotificationService = discordNotificationService;
    }

    public void notifySlackOnReportCreated(List<Long> receiverIds, String reportId) {
        receiverIds.forEach(id -> slackNotificationService.notifyOnReportGenerated(id, reportId));
    }

    public void notifyDiscordOnReportCreated(List<Long> receiverIds, String reportId) {
        receiverIds.forEach(id -> discordNotificationService.notifyOnReportGenerated(id, reportId));
    }

    public void notifyEmailOnReportCreated(List<Long> receiverIds, String reportId) {
        receiverIds.forEach(id -> emailNotificationService.notifyOnReportGenerated(id, reportId));
    }
}
