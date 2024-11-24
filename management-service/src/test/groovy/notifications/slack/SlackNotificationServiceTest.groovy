package notifications.slack

import org.springframework.beans.factory.annotation.Value
import pl.pwr.zpi.notifications.slack.SlackNotificationService
import pl.pwr.zpi.notifications.slack.entity.SlackReceiver
import pl.pwr.zpi.notifications.slack.service.SlackMessagingService
import pl.pwr.zpi.notifications.slack.service.SlackReceiverService
import pl.pwr.zpi.notifications.common.ConfidentialTextEncoder
import spock.lang.Ignore
import spock.lang.Specification
import spock.lang.Subject

class SlackNotificationServiceTest extends Specification {

    def slackService = Mock(SlackMessagingService)
    def receiverService = Mock(SlackReceiverService)
    def confidentialTextEncoder = Mock(ConfidentialTextEncoder)

    @Value("\${magpie.monitor.client.base.url}")
    def MAGPIE_MONITOR_CLIENT_BASE_URL = "http://localhost"

    @Subject
    def slackNotificationService = new SlackNotificationService(slackService, receiverService, confidentialTextEncoder)

    @Ignore
    def "should send test message successfully"() {
        given:
        def receiverSlackId = 1L
        def receiver = buildSlackReceiver(receiverSlackId, "https://webhook.url")
        receiverService.getById(receiverSlackId) >> receiver
        confidentialTextEncoder.decrypt(receiver.getWebhookUrl()) >> receiver.getWebhookUrl()

        when:
        slackNotificationService.sendTestMessage(receiverSlackId)

        then:
        1 * receiverService.getById(receiverSlackId)
        1 * confidentialTextEncoder.decrypt(receiver.getWebhookUrl())
        1 * slackService.sendMessage(_, receiver.getWebhookUrl())
    }

    def "should send test message by webhook URL"() {
        given:
        def webhookUrl = "https://webhook.url"

        when:
        slackNotificationService.sendTestMessage(webhookUrl)

        then:
        1 * slackService.sendMessage(_, webhookUrl)
    }

    @Ignore
    def "should notify on report generated successfully"() {
        given:
        def receiverId = 1L
        def reportId = "report123"
        def receiver = buildSlackReceiver(receiverId, "https://webhook.url")
        receiverService.getEncodedWebhookUrl(receiverId) >> receiver
        slackService.sendMessage(_, _) >> {}

        when:
        slackNotificationService.notifyOnReportGenerated(receiverId, reportId)

        then:
        1 * receiverService.getEncodedWebhookUrl(receiverId)
        1 * slackService.sendMessage(_, "http://localhost/reports/${reportId}")
    }

    def "should throw exception when notify on report generated fails"() {
        given:
        def receiverId = 1L
        def reportId = "report123"
        receiverService.getEncodedWebhookUrl(receiverId) >> { throw new Exception("Error retrieving webhook") }

        when:
        slackNotificationService.notifyOnReportGenerated(receiverId, reportId)

        then:
        thrown(RuntimeException)
    }

    private SlackReceiver buildSlackReceiver(Long id, String webhookUrl) {
        return SlackReceiver.builder()
                .id(id)
                .receiverName("Receiver Name")
                .webhookUrl(webhookUrl)
                .updatedAt(System.currentTimeMillis())
                .createdAt(System.currentTimeMillis())
                .build()
    }
}