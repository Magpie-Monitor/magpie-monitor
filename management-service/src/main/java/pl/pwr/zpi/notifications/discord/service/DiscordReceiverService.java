package pl.pwr.zpi.notifications.discord.service;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.notifications.common.ConfidentialTextEncoder;
import pl.pwr.zpi.notifications.discord.dto.DiscordReceiverDTO;
import pl.pwr.zpi.notifications.discord.dto.UpdateDiscordReceiverRequest;
import pl.pwr.zpi.notifications.discord.entity.DiscordReceiver;
import pl.pwr.zpi.notifications.discord.repository.DiscordRepository;

import java.util.List;
import java.util.regex.Pattern;
import java.util.stream.Collectors;
import java.util.stream.Stream;

@Service
@RequiredArgsConstructor
@Slf4j
public class DiscordReceiverService {

    private final DiscordRepository discordRepository;
    private final ConfidentialTextEncoder confidentialTextEncoder;

    private final String WEBHOOK_URL_REGEX = "https://discord.com/api/webhooks/[0-9]+/[a-zA-Z0-9\\-]+";

    public List<DiscordReceiver> getAllDiscordIntegrations() {
        List<DiscordReceiver> receivers = discordRepository.findAll();
        receivers.forEach(receiver -> receiver.setWebhookUrl(
                getAnonymizedWebhookUrl(receiver.getWebhookUrl())
        ));
        return receivers;
    }

    public void deleteDiscordReceiver(Long receiverId) {
        discordRepository.deleteById(receiverId);
    }

    public void createDiscordReceiver(DiscordReceiverDTO discordIntegration) {
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

    public DiscordReceiver updateDiscordIntegration(Long id, UpdateDiscordReceiverRequest updateDiscordReceiverRequest) {
        var receiver = getDiscordReceiver(id);

        String encryptedWebhookUrl = confidentialTextEncoder.encrypt(updateDiscordReceiverRequest.webhookUrl());
        checkIfUserCanUpdateWebhookUrl(encryptedWebhookUrl, id);

        patchReceiver(receiver, updateDiscordReceiverRequest);

        receiver.setUpdatedAt(System.currentTimeMillis());
        return discordRepository.save(receiver);
    }

    private void patchReceiver(DiscordReceiver discordReceiver, UpdateDiscordReceiverRequest updateRequest) {
        if (updateRequest.name() != null) {
            validateReceiverName(updateRequest.name());
            discordReceiver.setReceiverName(updateRequest.name());
        }

        if (updateRequest.webhookUrl() != null) {
            validateWebhookUrl(updateRequest.webhookUrl());
            discordReceiver.setWebhookUrl(updateRequest.webhookUrl());
        }
    }

    private void validateReceiverName(String name) {
        if (name.length() < 2 || name.length() > 100) {
            throw new RuntimeException("Receiver name length has to be >= 2 and <= 100");
        }
    }

    private void validateWebhookUrl(String webhookUrl) {
        if (!Pattern.matches(WEBHOOK_URL_REGEX, webhookUrl)) {
            throw new RuntimeException(String.format("webhookUrl has to satisfy the following regex - %s", WEBHOOK_URL_REGEX));
        }
    }

    public DiscordReceiver getDiscordReceiver(Long receiverWebhookId) {
        DiscordReceiver receiver = discordRepository.findById(receiverWebhookId)
                .orElseThrow(() -> new IllegalArgumentException("Webhook with given Id not found"));

        receiver.setWebhookUrl(getAnonymizedWebhookUrl(receiver.getWebhookUrl()));
        return receiver;
    }

    private String getAnonymizedWebhookUrl(String webhookUrl) {
        webhookUrl = confidentialTextEncoder.decrypt(webhookUrl);

        String[] webhookParts = webhookUrl.split("/");
        String authToken = webhookParts[webhookParts.length - 1];
        return joinWebhookWithoutAuthToken(webhookParts) + "/" + "*".repeat(authToken.length());
    }

    private String joinWebhookWithoutAuthToken(String[] webhookParts) {
        return Stream.of(webhookParts)
                .limit(webhookParts.length - 1)
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
