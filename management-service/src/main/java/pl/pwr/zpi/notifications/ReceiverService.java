package pl.pwr.zpi.notifications;

import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.notifications.discord.entity.DiscordReceiver;
import pl.pwr.zpi.notifications.discord.service.DiscordReceiverService;
import pl.pwr.zpi.notifications.email.entity.EmailReceiver;
import pl.pwr.zpi.notifications.email.service.EmailReceiverService;
import pl.pwr.zpi.notifications.slack.entity.SlackReceiver;
import pl.pwr.zpi.notifications.slack.service.SlackReceiverService;

@Service
@RequiredArgsConstructor
public class ReceiverService {

    private final SlackReceiverService slackReceiverService;
    private final DiscordReceiverService discordReceiverService;
    private final EmailReceiverService emailReceiverService;

    public SlackReceiver getSlackReceiverById(Long receiverId) {
        return slackReceiverService.getById(receiverId);
    }

    public DiscordReceiver getDiscordReceiverById(Long receiverId) {
        return discordReceiverService.getDiscordReceiver(receiverId);
    }

    public EmailReceiver getEmailReceiverById(Long receiverId) {
        return emailReceiverService.getEmailReceiver(receiverId);
    }
}
