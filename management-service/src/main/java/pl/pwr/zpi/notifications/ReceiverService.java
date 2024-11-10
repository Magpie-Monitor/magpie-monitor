package pl.pwr.zpi.notifications;

import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.notifications.slack.entity.SlackReceiver;
import pl.pwr.zpi.notifications.slack.service.SlackReceiverService;

@Service
@RequiredArgsConstructor
public class ReceiverService {

    private final SlackReceiverService slackReceiverService;

    public boolean slackReceiverExists(Long receiverId) {
        return slackReceiverService.existsById(receiverId);
    }

    public SlackReceiver getReceiverById(Long receiverId) {
        return slackReceiverService.getById(receiverId);
    }
}
