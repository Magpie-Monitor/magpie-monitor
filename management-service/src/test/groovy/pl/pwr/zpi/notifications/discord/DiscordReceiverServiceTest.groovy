package pl.pwr.zpi.notifications.discord

import pl.pwr.zpi.notifications.common.ConfidentialTextEncoder
import pl.pwr.zpi.notifications.discord.dto.DiscordReceiverDTO
import pl.pwr.zpi.notifications.discord.dto.UpdateDiscordReceiverRequest
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
        discordReceiverService.WEBHOOK_URL_REGEX = "https://discord.com/api/webhooks/[0-9]+/[a-zA-Z0-9_-]+"
    }

    def "should get all discord integrations"() {
        given:
        def encryptedWebhookUrl = "encryptedWebhook1"
        def decryptedWebhookUrl = "https://discord.com/api/webhooks/1234554321/xKh5vF0Som55bSex4q9slwOApmB0VXjcUoVS5Z9v9vu89snl-XeedfHj"
        def discordReceiverList = List.of(createDiscordReceiver(1L, "Receiver 1", encryptedWebhookUrl))

        when:
        def result = discordReceiverService.getAllDiscordIntegrations()

        then:
        result == discordReceiverList
        1 * discordRepository.findAll() >> discordReceiverList
        1 * confidentialTextEncoder.decrypt(encryptedWebhookUrl) >> decryptedWebhookUrl
    }

    def "should add new discord integration successfully"() {
        given:
        def encryptedWebhookUrl = "encryptedWebhook1"
        def decryptedWebhookUrl = "https://discord.com/api/webhooks/1234554321/xKh5vF0Som55bSex4q9slwOApmB0VXjcUoVS5Z9v9vu89snl-XeedfHj"
        def discordReceiverDTO = createDiscordReceiverDTO("Receiver 1", decryptedWebhookUrl)

        confidentialTextEncoder.encrypt(_) >> encryptedWebhookUrl
        discordRepository.existsByWebhookUrl(_) >> false

        when:
        discordReceiverService.createDiscordReceiver(discordReceiverDTO)

        then:
        1 * confidentialTextEncoder.encrypt(discordReceiverDTO.getWebhookUrl()) >> encryptedWebhookUrl
        1 * discordRepository.existsByWebhookUrl(encryptedWebhookUrl) >> false
        1 * discordRepository.save(_ as DiscordReceiver)
    }

    def "should throw exception if webhook already exists when adding new discord integration"() {
        given:
        def encryptedWebhookUrl = "encryptedWebhook1"
        def decryptedWebhookUrl = "https://discord.com/api/webhooks/1234554321/xKh5vF0Som55bSex4q9slwOApmB0VXjcUoVS5Z9v9vu89snl-XeedfHj"
        def discordReceiverDTO = createDiscordReceiverDTO("Receiver 1", decryptedWebhookUrl)

        confidentialTextEncoder.encrypt(_) >> encryptedWebhookUrl
        discordRepository.existsByWebhookUrl(_) >> true

        when:
        discordReceiverService.createDiscordReceiver(discordReceiverDTO)

        then:
        1 * confidentialTextEncoder.encrypt(discordReceiverDTO.getWebhookUrl()) >> encryptedWebhookUrl
        1 * discordRepository.existsByWebhookUrl(encryptedWebhookUrl) >> true
        0 * discordRepository.save(_ as DiscordReceiver)
        thrown(IllegalArgumentException)
    }

    def "should update discord integration successfully"() {
        given:
        def id = 1L
        def encryptedWebhookUrl = "https://discord.com/api/webhooks/1234554321/****"
        def decryptedWebhookUrl = "https://discord.com/api/webhooks/1234554321/xKh5vF0Som55bSex4q9slwOApmB0VXjcUoVS5Z9v9vu89snl-XeedfHj"
        def discordReceiverUpdateRequest = new UpdateDiscordReceiverRequest("Updated Receiver", decryptedWebhookUrl)
        def existingReceiver = createDiscordReceiver(id, "Receiver 1", encryptedWebhookUrl)

        discordRepository.findById(id) >> Optional.of(existingReceiver)
        confidentialTextEncoder.decrypt(encryptedWebhookUrl) >> decryptedWebhookUrl
        confidentialTextEncoder.encrypt(decryptedWebhookUrl) >> encryptedWebhookUrl
        discordRepository.existsByWebhookUrl(_) >> false
        discordRepository.save(_) >> existingReceiver

        when:
        def updatedReceiver = discordReceiverService.updateDiscordIntegration(id, discordReceiverUpdateRequest)

        then:
        updatedReceiver != null
        updatedReceiver.receiverName == "Updated Receiver"
        updatedReceiver.webhookUrl == encryptedWebhookUrl
        1 * discordRepository.findById(id) >> Optional.of(existingReceiver)
        2 * confidentialTextEncoder.encrypt(decryptedWebhookUrl) >> encryptedWebhookUrl
        1 * discordRepository.existsByWebhookUrl(encryptedWebhookUrl) >> false
        1 * discordRepository.save(_ as DiscordReceiver) >> existingReceiver
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
