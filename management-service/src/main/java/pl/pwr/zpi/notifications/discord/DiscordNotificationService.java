package pl.pwr.zpi.notifications.discord;

import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.notifications.ReportNotifier;
import pl.pwr.zpi.notifications.common.ResourceLoaderUtils;
import pl.pwr.zpi.notifications.discord.entity.DiscordReceiver;
import pl.pwr.zpi.notifications.discord.service.DiscordMessagingService;
import pl.pwr.zpi.notifications.discord.service.DiscordReceiverService;

@Service("discordNotificationService")
@RequiredArgsConstructor
public class DiscordNotificationService implements ReportNotifier {

    private final String TEST_MESSAGE_PATH = "discord/test-message.txt";
    private final String NEW_REPORT_MESSAGE_PATH = "discord/new-report-message.txt";
    private final String MAGPIE_MONITOR_CLIENT_BASE_URL = "https://magpie-monitor.rolo-labs.xyz/reports";

    private final DiscordMessagingService discordMessagingService;
    private final DiscordReceiverService discordReceiverService;

    public void sendTestMessage(String webhookUrl) {
        discordMessagingService.sendMessage(
                ResourceLoaderUtils.loadResourceToString(TEST_MESSAGE_PATH),
                webhookUrl);
    }

    @Override
    public void notifyOnReportGenerated(Long receiverId, String reportId) {
        DiscordReceiver discordReceiver = discordReceiverService.getDiscordReceiver(receiverId);

        discordMessagingService.sendMessage(
                ResourceLoaderUtils.loadResourceToString(NEW_REPORT_MESSAGE_PATH) +
                        String.format("%s/%s", MAGPIE_MONITOR_CLIENT_BASE_URL, reportId),
                discordReceiver.getWebhookUrl());
    }

    @Override
    public void notifyOnReportGenerationFailed(Long receiverId, String reportId) {

    }
}
