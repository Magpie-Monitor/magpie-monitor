package pl.pwr.zpi.notifications;

import pl.pwr.zpi.notifications.slack.controller.SlackReceiver;

import java.time.LocalDateTime;

public record SlackChannel(
        String name,
        String server,
        LocalDateTime updatedAt,
        LocalDateTime createdAt) {

    public static SlackChannel of(SlackReceiver receiver) {
        return new SlackChannel(
                receiver.getReceiverName(),
                receiver.getWebhookUrl(),
                receiver.getUpdatedAt(),
                receiver.getCreatedAt()
        );
    }
}
