package pl.pwr.zpi.notifications.discord

import pl.pwr.zpi.notifications.common.ConfidentialTextEncoder
import pl.pwr.zpi.notifications.discord.dto.DiscordReceiverDTO
import pl.pwr.zpi.notifications.discord.entity.DiscordReceiver
import pl.pwr.zpi.notifications.discord.repository.DiscordRepository
import pl.pwr.zpi.notifications.discord.service.DiscordReceiverService
import spock.lang.Specification

class DiscordReceiverServiceTest extends Specification {

    DiscordRepository discordRepository
    ConfidentialTextEncoder confidentialTextEncoder
    DiscordReceiverService discordReceiverService

    def setup() {
        confidentialTextEncoder = Mock()
        discordRepository = Mock()
        discordReceiverService = new DiscordReceiverService(discordRepository, confidentialTextEncoder)
    }

    def "should get all discord integrations"() {
        given:
        def discordReceiverList = List.of(createDiscordReceiver(1L, "Receiver 1", "encryptedWebhook1"))

        when:
        def result = discordReceiverService.getAllDiscordIntegrations()

        then:
        result == discordReceiverList
        1 * discordRepository.findAll() >> discordReceiverList
    }

    def "should add new discord integration successfully"() {
        given:
        def discordReceiverDTO =createDiscordReceiverDTO("Receiver 1", "http://webhook1")
        def encryptedWebhookUrl = "encryptedWebhook1"

        confidentialTextEncoder.encrypt(_) >> encryptedWebhookUrl
        discordRepository.existsByWebhookUrl(_) >> false

        when:
        discordReceiverService.addNewDiscordIntegration(discordReceiverDTO)

        then:
        1 * confidentialTextEncoder.encrypt(discordReceiverDTO.getWebhookUrl()) >> encryptedWebhookUrl
        1 * discordRepository.existsByWebhookUrl(encryptedWebhookUrl) >> false
        1 * discordRepository.save(_ as DiscordReceiver)
    }

    def "should throw exception if webhook already exists when adding new discord integration"() {
        given:
        def discordReceiverDTO = createDiscordReceiverDTO("Receiver 1", "http://webhook1")
        def encryptedWebhookUrl = "encryptedWebhook1"

        confidentialTextEncoder.encrypt(_) >> encryptedWebhookUrl
        discordRepository.existsByWebhookUrl(_) >> true

        when:
        discordReceiverService.addNewDiscordIntegration(discordReceiverDTO)

        then:
        1 * confidentialTextEncoder.encrypt(discordReceiverDTO.getWebhookUrl()) >> encryptedWebhookUrl
        1 * discordRepository.existsByWebhookUrl(encryptedWebhookUrl) >> true
        0 * discordRepository.save(_ as DiscordReceiver)
        thrown(IllegalArgumentException)
    }

    def "should update discord integration successfully"() {
        given:
        def id = 1L
        def discordReceiverDTO = createDiscordReceiverDTO("Updated Receiver", "http://webhook1")
        def encryptedWebhookUrl = "encryptedUpdatedWebhook"
        def existingReceiver = createDiscordReceiver(id, "Receiver 1", "oldEncryptedWebhook")

        discordRepository.findById(id) >> Optional.of(existingReceiver)
        confidentialTextEncoder.encrypt(discordReceiverDTO.getWebhookUrl()) >> encryptedWebhookUrl
        discordRepository.existsByWebhookUrl(_) >> false
        discordRepository.save(_) >> existingReceiver

        when:
        def updatedReceiver = discordReceiverService.updateDiscordIntegration(id, discordReceiverDTO)

        then:
        updatedReceiver != null
        updatedReceiver.receiverName == "Updated Receiver"
        updatedReceiver.webhookUrl == encryptedWebhookUrl
        1 * discordRepository.findById(id) >> Optional.of(existingReceiver)
        1 * confidentialTextEncoder.encrypt(discordReceiverDTO.getWebhookUrl()) >> encryptedWebhookUrl
        1 * discordRepository.existsByWebhookUrl(encryptedWebhookUrl) >> false
        1 * discordRepository.save(_ as DiscordReceiver) >> existingReceiver
    }

    def "should throw exception if webhook URL is assigned to another entry when updating discord integration"() {
        given:
        def id = 1L
        def discordReceiverDTO = createDiscordReceiverDTO("Receiver 1", "https://discord.com/api/webhooks/1234554321/xKh5vF0Som55bSex4q9slwOApmB0VXjcUoVS5Z9v9vu89snl-XeedfHj")
        def encryptedWebhookUrl = "encryptedUpdatedWebhook"
        def existingReceiver = DiscordReceiver.builder()
                .id(id)
                .receiverName("Receiver 1")
                .webhookUrl("encryptedWebhook")
                .createdAt(System.currentTimeMillis())
                .updatedAt(System.currentTimeMillis())
                .build()

        discordRepository.findById(id) >> Optional.of(existingReceiver)
        confidentialTextEncoder.encrypt(discordReceiverDTO.getWebhookUrl()) >> encryptedWebhookUrl
        discordRepository.existsByWebhookUrl(encryptedWebhookUrl) >> true

        when:
        discordReceiverService.updateDiscordIntegration(id, discordReceiverDTO)

        then:
        thrown(IllegalArgumentException)
        2 * discordRepository.findById(id) >> Optional.of(existingReceiver)
        1 * confidentialTextEncoder.encrypt(_ as String) >> encryptedWebhookUrl
        1 * discordRepository.existsByWebhookUrl(encryptedWebhookUrl) >> true
        0 * discordRepository.save(_ as DiscordReceiver)
    }

    def "should throw exception if discord receiver not found"() {
        given:
        def id = 1L
        discordRepository.findById(id) >> Optional.empty()

        when:
        discordReceiverService.getDiscordReceiver(id)

        then:
        thrown(IllegalArgumentException)
        1 * discordRepository.findById(id) >> Optional.empty()
    }

    def "should get encoded webhook URL successfully"() {
        given:
        def id = 1L
        def discordReceiver = createDiscordReceiver(id, "Receiver 1", "encryptedWebhook1")

        def decryptedWebhookUrl = "https://discord.com/api/webhooks/1234554321/xKh5vF0Som55bSex4q9slwOApmB0VXjcUoVS5Z9v9vu89snl-XeedfHj"

        discordRepository.findById(id) >> Optional.of(discordReceiver)
        confidentialTextEncoder.decrypt(_) >> decryptedWebhookUrl

        when:
        def receiver = discordReceiverService.getEncodedWebhookUrl(id)

        then:
        receiver.webhookUrl == decryptedWebhookUrl
        1 * discordRepository.findById(id) >> Optional.of(discordReceiver)
        1 * confidentialTextEncoder.decrypt(discordReceiver.webhookUrl) >> decryptedWebhookUrl
    }

    private DiscordReceiverDTO createDiscordReceiverDTO(String name, String webhookUrl) {
        return DiscordReceiverDTO.builder()
                .name(name)
                .webhookUrl(webhookUrl)
                .build()
    }

    private DiscordReceiver createDiscordReceiver(Long id, String name, String webhookUrl) {
        return DiscordReceiver.builder()
                .id(id)
                .receiverName(name)
                .webhookUrl(webhookUrl)
                .createdAt(System.currentTimeMillis())
                .updatedAt(System.currentTimeMillis())
                .build()
    }
}
