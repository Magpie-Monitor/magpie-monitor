package pl.pwr.zpi.notifications.discord.dto;

public record UpdateDiscordReceiverRequest(
        String name,
        String webhookUrl
) {
}
