package pl.pwr.zpi.notifications.slack.service;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.notifications.common.ConfidentialTextEncoder;
import pl.pwr.zpi.notifications.slack.dto.SlackReceiverDTO;
import pl.pwr.zpi.notifications.slack.entity.SlackReceiver;
import pl.pwr.zpi.notifications.slack.repository.SlackRepository;

import java.util.List;
import java.util.stream.Collectors;
import java.util.stream.Stream;

@Service
@RequiredArgsConstructor
@Slf4j
public class SlackReceiverService {

    private final SlackRepository slackRepository;
    private final ConfidentialTextEncoder confidentialTextEncoder;

    public List<SlackReceiver> getAllSlackIntegrations() {
        return slackRepository.findAll();
    }

    public void addNewSlackIntegration(SlackReceiverDTO slackIntegration) throws Exception {
        long now = System.currentTimeMillis();

        String encryptedWebhookUrl = confidentialTextEncoder.encrypt(slackIntegration.getWebhookUrl());
        checkIfWebhookExists(encryptedWebhookUrl);
        SlackReceiver receiver = SlackReceiver.builder()
                .receiverName(slackIntegration.getName())
                .webhookUrl(encryptedWebhookUrl)
                .createdAt(now)
                .updatedAt(now)
                .build();
        slackRepository.save(receiver);
    }

    public SlackReceiver updateSlackIntegration(Long id, SlackReceiverDTO slackReceiver) {
        var receiver = getById(id);

        receiver.setReceiverName(slackReceiver.getName());
        receiver.setWebhookUrl(getAnonymizedWebhookUrl(slackReceiver.getWebhookUrl()));
        receiver.setUpdatedAt(System.currentTimeMillis());

        return slackRepository.save(receiver);
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

    public SlackReceiver getById(Long receiverId) {
        return slackRepository.findById(receiverId)
                .orElseThrow(() -> new IllegalArgumentException("Webhook with given clusterId not found"));
    }

    private void checkIfWebhookExists(String webhookUrl) {
        if (slackRepository.existsByWebhookUrl(webhookUrl)) {
            throw new IllegalArgumentException("Webhook already exists");
        }
    }

    private void checkIfUserCanUpdateWebhookUrl(String webhookUrl, Long id) {
        if (slackRepository.existsByWebhookUrl(webhookUrl) && !slackRepository.findById(id).get().getWebhookUrl().equals(webhookUrl)) {
            throw new IllegalArgumentException("Webhook is already assigned to other entry");
        }
    }

    public SlackReceiver getEncodedWebhookUrl(Long id) throws Exception {
        var receiver = getById(id);
        receiver.setWebhookUrl(confidentialTextEncoder.decrypt(receiver.getWebhookUrl()));
        return receiver;
    }

    public void deleteSlackReceiver(Long receiverId) {
        slackRepository.deleteById(receiverId);
    }
}
