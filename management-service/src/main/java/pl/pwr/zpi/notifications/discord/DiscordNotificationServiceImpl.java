package pl.pwr.zpi.notifications.discord;

import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.notifications.common.ResourceLoaderUtils;

@RequiredArgsConstructor
@Service
public class DiscordNotificationServiceImpl implements DiscordNotificationService {
    private final String TEST_MESSAGE_PATH = "discord/test-message.txt";
    private final String NEW_REPORT_MESSAGE_PATH = "discord/new-report-message.txt";

    private final DiscordService discordService;

    @Override
    public void sendTestMessage(String webhookUrl) {
        discordService.sendMessage(
                ResourceLoaderUtils.loadResourceToString(TEST_MESSAGE_PATH),
                webhookUrl);
    }

    @Override
    public void sendMessageAboutNewReport(String webhookUrl, String reportUrl) {
        discordService.sendMessage(
                ResourceLoaderUtils.loadResourceToString(NEW_REPORT_MESSAGE_PATH) + reportUrl,
                webhookUrl);
    }
}
