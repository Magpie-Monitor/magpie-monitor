package pl.pwr.zpi.notifications.discord;

public interface DiscordService {
    void sendMessage(String message, String webhookUrl);
}
