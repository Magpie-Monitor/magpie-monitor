package pl.pwr.zpi.notifications.slack;

import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.notifications.ReportNotifier;
import pl.pwr.zpi.notifications.common.ConfidentialTextEncoder;
import pl.pwr.zpi.notifications.common.ResourceLoaderUtils;
import pl.pwr.zpi.notifications.slack.entity.SlackReceiver;
import pl.pwr.zpi.notifications.slack.service.SlackMessagingService;
import pl.pwr.zpi.notifications.slack.service.SlackReceiverService;

@Service("slackNotificationService")
@RequiredArgsConstructor
public class SlackNotificationService implements ReportNotifier {

    private final String TEST_MESSAGE_PATH = "slack/test-message.txt";
    private final String NEW_REPORT_MESSAGE_PATH = "slack/new-report-message.txt";
    private final String MAGPIE_MONITOR_CLIENT_BASE_URL = "https://magpie-monitor.rolo-labs.xyz/reports";

    private final SlackMessagingService slackService;
    private final SlackReceiverService receiverService;
    private final ConfidentialTextEncoder confidentialTextEncoder;

    public void sendTestMessage(Long receiverSlackId) throws Exception {
        SlackReceiver receiver = receiverService.getById(receiverSlackId);
        String decodedWebhookUrl = confidentialTextEncoder.decrypt(receiver.getWebhookUrl());
        sendTestMessage(decodedWebhookUrl);
    }

    public void sendTestMessage(String webhookUrl) {
        slackService.sendMessage(
                loadResource(TEST_MESSAGE_PATH),
                webhookUrl
        );
    }


    @Override
    public void notifyOnReportGenerated(Long receiverId, String reportId) {
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

    // TODO - implement
    @Override
    public void notifyOnReportGenerationFailed(Long receiverId, String clusterId) {
    }

    private String loadResource(String resourcePath) {
        return ResourceLoaderUtils.loadResourceToString(resourcePath);
    }

    private String getReportUrl(String reportId) {
        return String.format("%s/%s", MAGPIE_MONITOR_CLIENT_BASE_URL, reportId);
    }

}
