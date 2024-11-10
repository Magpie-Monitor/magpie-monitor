package pl.pwr.zpi.notifications.slack.service;

import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.notifications.common.ResourceLoaderUtils;

@Service
@RequiredArgsConstructor
public class SlackNotificationServiceImpl implements SlackNotificationService {

    private final String TEST_MESSAGE_PATH = "slack/test-message.txt";
    private final String NEW_REPORT_MESSAGE_PATH = "slack/new-report-message.txt";
    private final String MAGPIE_MONITOR_CLIENT_BASE_URL = "https://magpie-monitor.rolo-labs.xyz/reports";

    private final SlackMessagingService slackService;
    private final SlackReceiverService receiverService;

//    public void sendTestMessage(Long receiverSlackId) {
//        var receiver = getById(receiverSlackId);
//        String decodedWebhookUrl = confidentialTextEncoder.decrypt(receiver.getWebhookUrl());
//        slackNotificationService.sendTestMessage(decodedWebhookUrl);
//    }

    @Override
    public void sendTestMessage(String webhookUrl) {
        slackService.sendMessage(
                loadResource(TEST_MESSAGE_PATH),
                webhookUrl
        );
    }

    @Override
    public void notifyOnReportCreated(Long receiverId, String reportId) {
        try {
            String webhookUrl = receiverService.getEncodedWebhookUrl(receiverId).getWebhookUrl();
            slackService.sendMessage(
                    String.format("%s: %s", loadResource(NEW_REPORT_MESSAGE_PATH), getReportUrl(reportId)),
                    webhookUrl
            );
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
    }

    private String loadResource(String resourcePath) {
        return ResourceLoaderUtils.loadResourceToString(resourcePath);
    }

    private String getReportUrl(String reportId) {
        return String.format("%s/%s", MAGPIE_MONITOR_CLIENT_BASE_URL, reportId);
    }
}
