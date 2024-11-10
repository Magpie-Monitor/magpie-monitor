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

    public void notifySlackOnReportGenerationFailed(List<Long> receiverIds, String clusterId) {
        receiverIds.forEach(id -> slackNotificationService.notifyOnReportGenerationFailed(id, clusterId));
    }

    public void notifyDiscordOnReportCreated(List<Long> receiverIds, String reportId) {
        receiverIds.forEach(id -> discordNotificationService.notifyOnReportGenerated(id, reportId));
    }

    public void notifyDiscordOnReportGenerationFailed(List<Long> receiverIds, String clusterId) {
        receiverIds.forEach(id -> discordNotificationService.notifyOnReportGenerationFailed(id, clusterId));
    }

    public void notifyEmailOnReportCreated(List<Long> receiverIds, String reportId) {
        receiverIds.forEach(id -> emailNotificationService.notifyOnReportGenerated(id, reportId));
    }

    public void notifyEmailOnReportGenerationFailed(List<Long> receiverIds, String clusterId) {
        receiverIds.forEach(id -> emailNotificationService.notifyOnReportGenerationFailed(id, clusterId));
    }
}
