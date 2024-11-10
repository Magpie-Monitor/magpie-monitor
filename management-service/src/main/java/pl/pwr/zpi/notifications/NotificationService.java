package pl.pwr.zpi.notifications;

import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.notifications.slack.service.SlackNotificationService;

import java.util.List;

@Service
@RequiredArgsConstructor
public class NotificationService {

    private final SlackNotificationService slackNotificationService;

    public void notifySlack(List<Long> receiverIds, String reportId) {
        receiverIds.forEach(id -> slackNotificationService.notifyOnReportCreated(id, reportId));
    }

    public void notifyDiscord(List<Long> receiverIds, String reportId) {

    }

    public void notifyEmail(List<Long> receiverIds, String reportId) {

    }
}
