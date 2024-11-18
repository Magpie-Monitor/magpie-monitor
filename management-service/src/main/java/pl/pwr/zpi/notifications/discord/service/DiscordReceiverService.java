package pl.pwr.zpi.notifications.discord.service;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.notifications.common.ConfidentialTextEncoder;
import pl.pwr.zpi.notifications.discord.dto.DiscordReceiverDTO;
import pl.pwr.zpi.notifications.discord.entity.DiscordReceiver;
import pl.pwr.zpi.notifications.discord.repository.DiscordRepository;

import java.util.List;
import java.util.stream.Collectors;
import java.util.stream.Stream;

@Service
@RequiredArgsConstructor
@Slf4j
public class DiscordReceiverService {

    private final DiscordRepository discordRepository;
    private final ConfidentialTextEncoder confidentialTextEncoder;

    public List<DiscordReceiver> getAllDiscordIntegrations() {
        return discordRepository.findAll();
    }

    public void deleteDiscordReceiver(Long receiverId) {
        discordRepository.deleteById(receiverId);
    }

    public void createDiscordReceiver(DiscordReceiverDTO discordIntegration) throws Exception {
        long now = System.currentTimeMillis();

        String encryptedWebhookUrl = confidentialTextEncoder.encrypt(discordIntegration.getWebhookUrl());
        checkIfWebhookExists(encryptedWebhookUrl);
        DiscordReceiver receiver = DiscordReceiver.builder()
                .receiverName(discordIntegration.getName())
                .webhookUrl(encryptedWebhookUrl)
                .createdAt(now)
                .updatedAt(now)
                .build();
        discordRepository.save(receiver);
    }

    public DiscordReceiver updateDiscordIntegration(Long id, DiscordReceiverDTO discordReceiver) throws Exception {
        var receiver = getDiscordReceiver(id);
        String encryptedWebhookUrl = confidentialTextEncoder.encrypt(discordReceiver.getWebhookUrl());
        checkIfUserCanUpdateWebhookUrl(encryptedWebhookUrl, id);

        receiver.setReceiverName(discordReceiver.getName());
        receiver.setWebhookUrl(encryptedWebhookUrl);
        receiver.setUpdatedAt(System.currentTimeMillis());
        return discordRepository.save(receiver);
    }

    public DiscordReceiver getDiscordReceiver(Long receiverWebhookId) {
        DiscordReceiver receiver = discordRepository.findById(receiverWebhookId)
                .orElseThrow(() -> new IllegalArgumentException("Webhook with given Id not found"));

        receiver.setWebhookUrl(getAnonymizedWebhookUrl(receiver.getWebhookUrl()));
        return receiver;
    }

    private String getAnonymizedWebhookUrl(String webhookUrl) {
        String[] webhookParts = webhookUrl.split("/");
        String authToken = webhookParts[webhookParts.length - 1];
        return joinWebhookWithoutAuthToken(webhookParts) + authToken.replace(".", "*");
    }

    private String joinWebhookWithoutAuthToken(String[] webhookParts) {
        return Stream.of(webhookParts)
                .limit(webhookParts.length - 2)
                .collect(Collectors.joining("/"));
    }

    private void checkIfWebhookExists(String webhookUrl) {
        if (discordRepository.existsByWebhookUrl(webhookUrl)) {
            throw new IllegalArgumentException("Webhook already exists");
        }
    }

    private void checkIfUserCanUpdateWebhookUrl(String webhookUrl, Long id) {
        if (discordRepository.existsByWebhookUrl(webhookUrl) && !discordRepository.findById(id).get().getWebhookUrl().equals(webhookUrl)) {
            throw new IllegalArgumentException("Webhook is already assigned to other entry");
        }
    }
}
