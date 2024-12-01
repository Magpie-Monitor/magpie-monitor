package pl.pwr.zpi.notifications.slack.dto;

public record UpdateSlackReceiverRequest(
        String name,
        String webhookUrl
) {
}
