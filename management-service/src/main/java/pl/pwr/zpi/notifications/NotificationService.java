package pl.pwr.zpi.notifications;

import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.notifications.slack.controller.SlackService;

import java.util.List;

@Service
@RequiredArgsConstructor
public class NotificationService {

    private final SlackService slackService;

    public List<SlackChannel> getSlackReceivers() {
        return slackService.getAllSlackIntegrations().stream()
                .map(SlackChannel::of)
                .toList();
    }

    // TODO - implement notifications
    public void notifySlack(List<Long> receiverIds) {

    }

    public void notifyDiscord(List<Long> receiverIds) {

    }

    public void notifyEmail(List<Long> receiverIds) {

    }
}
