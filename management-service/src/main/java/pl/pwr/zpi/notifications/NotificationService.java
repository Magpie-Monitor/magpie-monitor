package pl.pwr.zpi.notifications;

import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.notifications.slack.controller.SlackReceiver;
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
}
