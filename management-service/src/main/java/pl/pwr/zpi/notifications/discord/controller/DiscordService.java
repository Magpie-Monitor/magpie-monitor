package pl.pwr.zpi.notifications.discord.controller;

import lombok.RequiredArgsConstructor;
import lombok.SneakyThrows;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.notifications.common.ConfidentialTextEncoder;
import pl.pwr.zpi.notifications.discord.DiscordNotificationService;

import java.time.LocalDateTime;
import java.util.List;

@Service
@RequiredArgsConstructor
@Slf4j
public class DiscordService {

    private final DiscordNotificationService discordNotificationService;
    private final DiscordRepository discordRepository;
    private final ConfidentialTextEncoder confidentialTextEncoder;

    public void sendTestMessage(Long receiverDiscordId) throws Exception {
        var receiver = getDiscordReceiver(receiverDiscordId);
        String decodedWebhookUrl = confidentialTextEncoder.decrypt(receiver.getWebhookUrl());
        discordNotificationService.sendTestMessage(decodedWebhookUrl);
    }

    public List<DiscordReceiver> getAllDiscordIntegrations() {
        return discordRepository.findAll();
    }

    public void addNewDiscordIntegration(DiscordReceiverDTO discordIntegration) throws Exception {
        String encryptedWebhookUrl = confidentialTextEncoder.encrypt(discordIntegration.getWebhookUrl());
        checkIfWebhookExists(encryptedWebhookUrl);
        DiscordReceiver receiver = DiscordReceiver.builder()
                .receiverName(discordIntegration.getName())
                .webhookUrl(encryptedWebhookUrl)
                .createdAt(LocalDateTime.now())
                .build();
        discordRepository.save(receiver);
    }

    public DiscordReceiver updateDiscordIntegration(Long id, DiscordReceiverDTO discordReceiver) throws Exception {
        var receiver = getDiscordReceiver(id);
        String encryptedWebhookUrl = confidentialTextEncoder.encrypt(discordReceiver.getWebhookUrl());
        checkIfUserCanUpdateWebhookUrl(encryptedWebhookUrl, id);

        receiver.setReceiverName(discordReceiver.getName());
        receiver.setWebhookUrl(encryptedWebhookUrl);
        receiver.setUpdatedAt(LocalDateTime.now());
        return discordRepository.save(receiver);
    }

    private DiscordReceiver getDiscordReceiver(Long receiverWebhookId) {
        return discordRepository.findById(receiverWebhookId)
                .orElseThrow(() -> new IllegalArgumentException("Webhook with given id not found"));
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

    public DiscordReceiver getEncodedWebhookUrl(Long id) throws Exception {
        var receiver = getDiscordReceiver(id);
        receiver.setWebhookUrl(confidentialTextEncoder.decrypt(receiver.getWebhookUrl()));
        return receiver;
    }
}
