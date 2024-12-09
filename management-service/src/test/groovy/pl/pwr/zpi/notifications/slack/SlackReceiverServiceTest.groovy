package pl.pwr.zpi.notifications.slack

import pl.pwr.zpi.notifications.common.ConfidentialTextEncoder
import pl.pwr.zpi.notifications.slack.dto.SlackReceiverDTO
import pl.pwr.zpi.notifications.slack.dto.UpdateSlackReceiverRequest
import pl.pwr.zpi.notifications.slack.entity.SlackReceiver
import pl.pwr.zpi.notifications.slack.repository.SlackRepository
import pl.pwr.zpi.notifications.slack.service.SlackReceiverService
import spock.lang.Specification
import spock.lang.Subject

class SlackReceiverServiceTest extends Specification {

    def slackRepository
    def confidentialTextEncoder

    @Subject
    def slackReceiverService

    def setup() {
        slackRepository = Mock(SlackRepository)
        confidentialTextEncoder = Mock(ConfidentialTextEncoder)
        slackReceiverService = new SlackReceiverService(slackRepository, confidentialTextEncoder)
        slackReceiverService.WEBHOOK_URL_REGEX = "https://hooks.slack.com/services/[A-Z0-9]+/[A-Z0-9]+/[a-zA-Z0-9]+"
    }

    def "getAllSlackIntegrations should anonymize webhook URL for each receiver"() {
        given:
        def receiver1 = new SlackReceiver(1L, "Receiver1", "encryptedUrl1", System.currentTimeMillis(), System.currentTimeMillis())
        def receiver2 = new SlackReceiver(2L, "Receiver2", "encryptedUrl2", System.currentTimeMillis(), System.currentTimeMillis())
        slackRepository.findAll() >> [receiver1, receiver2]
        confidentialTextEncoder.decrypt(_) >> { args -> args[0] }

        when:
        def result = slackReceiverService.getAllSlackIntegrations()

        then:
        1 * confidentialTextEncoder.decrypt("encryptedUrl1") >> "https://slack.com/receiver1/token1"
        1 * confidentialTextEncoder.decrypt("encryptedUrl2") >> "https://slack.com/receiver2/token2"
        result.size() == 2
        result[0].webhookUrl == "https://slack.com/receiver1/******"
        result[1].webhookUrl == "https://slack.com/receiver2/******"
    }

    def "getEncodedWebhookUrl should return receiver with decoded webhook URL"() {
        given:
        def receiver = new SlackReceiver(id: 1L, receiverName: "Receiver1", webhookUrl: "encryptedUrl")
        slackRepository.findById(1L) >> Optional.of(receiver)
        confidentialTextEncoder.decrypt(_) >> "https://hooks.slack.com/services/T04PB0Y4K8Q/B07QG098S7M/Xk3uMvmSOCsFhhTWPSGA"

        when:
        def result = slackReceiverService.getEncodedWebhookUrl(1L)

        then:
        result.webhookUrl == "https://hooks.slack.com/services/T04PB0Y4K8Q/B07QG098S7M/Xk3uMvmSOCsFhhTWPSGA"
    }

    def "validateReceiverName should throw exception for short name"() {
        given:
        def shortName = "A"

        when:
        slackReceiverService.validateReceiverName(shortName)

        then:
        thrown(RuntimeException)
    }

    def "validateWebhookUrl should throw exception for invalid URL"() {
        given:
        def invalidUrl = "invalidUrl"

        when:
        slackReceiverService.validateWebhookUrl(invalidUrl)

        then:
        thrown(RuntimeException)
    }

    def "checkIfUserCanUpdateWebhookUrl should throw exception if URL is already used"() {
        given:
        slackRepository.findById(1L) >> Optional.of(new SlackReceiver(id: 1L, receiverName: "Receiver1", webhookUrl: "encryptedUrl"))
        slackRepository.existsByWebhookUrl("newWebhookUrl") >> true

        when:
        slackReceiverService.checkIfUserCanUpdateWebhookUrl("newWebhookUrl", 1L)

        then:
        thrown(IllegalArgumentException)
    }

    def "patchReceiver should update receiver name and webhook URL"() {
        given:
        def receiver = new SlackReceiver(id: 1L, receiverName: "OldReceiver", webhookUrl: "newencryptedUrl")
        def updateRequest = new UpdateSlackReceiverRequest("UpdatedReceiver", "https://hooks.slack.com/services/T04PB0Y4K8Q/B07QG098S7M/Xk3uMvmSOCsFhhTWPSGA")
        slackRepository.findById(1L) >> Optional.of(receiver)
        slackRepository.existsByWebhookUrl("EncryptedUrl") >> true
        confidentialTextEncoder.encrypt(_) >> "newEncryptedUrl"
        slackRepository.save(receiver) >> receiver
        confidentialTextEncoder.decrypt(_) >> "https://hooks.slack.com/services/T04PB0Y4K8Q/B07QG098S7M/Xk3uMvmSOCsFhhTWPSGA"

        when:
        slackReceiverService.updateSlackIntegration(1L, updateRequest)

        then:
        receiver.receiverName == "UpdatedReceiver"
        receiver.webhookUrl == "https://hooks.slack.com/services/T04PB0Y4K8Q/B07QG098S7M/********************"
    }

    def "addNewSlackIntegration should save new slack integration"() {
        given:
        def slackReceiverDTO = buildSlackReceiverDTO("NewReceiver", "https://slack.com/webhook/abc123")
        def encryptedUrl = "encryptedUrl"
        confidentialTextEncoder.encrypt(_) >> encryptedUrl
        slackRepository.existsByWebhookUrl(_) >> false

        when:
        slackReceiverService.addNewSlackIntegration(slackReceiverDTO)

        then:
        1 * slackRepository.save(_ as SlackReceiver)
    }

    def "addNewSlackIntegration should throw exception if webhook already exists"() {
        given:
        def slackReceiverDTO = buildSlackReceiverDTO("NewReceiver", "https://slack.com/webhook/abc123")
        def encryptedUrl = "encryptedUrl"
        confidentialTextEncoder.encrypt(_) >> encryptedUrl
        slackRepository.existsByWebhookUrl(_) >> true

        when:
        slackReceiverService.addNewSlackIntegration(slackReceiverDTO)

        then:
        thrown(IllegalArgumentException)
    }

    def "getById should return the SlackReceiver"() {
        given:
        def receiver = new SlackReceiver(id: 1L, receiverName: "Receiver1", webhookUrl: "encryptedUrl")
        slackRepository.findById(1L) >> Optional.of(receiver)

        when:
        def result = slackReceiverService.getById(1L)

        then:
        result.id == 1L
        result.receiverName == "Receiver1"
    }

    def "getById should throw exception if receiver is not found"() {
        given:
        slackRepository.findById(1L) >> Optional.empty()

        when:
        slackReceiverService.getById(1L)

        then:
        thrown(IllegalArgumentException)
    }

    def "deleteSlackReceiver should delete the receiver"() {
        given:
        slackRepository.existsById(1L) >> true

        when:
        slackReceiverService.deleteSlackReceiver(1L)

        then:
        1 * slackRepository.deleteById(1L)
    }

    def "deleteSlackReceiver should throw exception if receiver is not found"() {
        given:
        slackRepository.existsById(1L) >> false

        when:
        slackReceiverService.deleteSlackReceiver(1L)

        then:
        thrown(IllegalArgumentException)
    }

    private buildSlackReceiverDTO(String name, String webhookUrl) {
        return SlackReceiverDTO.builder()
                .name(name)
                .webhookUrl(webhookUrl)
                .build()
    }
}
