package pl.pwr.zpi.notifications;

import pl.pwr.zpi.notifications.slack.controller.SlackReceiver;

public record SlackChannel(
        String name,
        String server,
        Long updatedAt,
        Long createdAt) {

    public static SlackChannel of(SlackReceiver receiver) {
        return new SlackChannel(
                receiver.getReceiverName(),
                receiver.getWebhookUrl(),
                receiver.getUpdatedAt(),
                receiver.getCreatedAt()
        );
    }
}
