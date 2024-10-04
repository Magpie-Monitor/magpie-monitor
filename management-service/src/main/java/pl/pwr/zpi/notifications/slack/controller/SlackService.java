package pl.pwr.zpi.notifications.slack.controller;

import lombok.RequiredArgsConstructor;
import lombok.SneakyThrows;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.notifications.slack.SlackNotificationService;
import pl.pwr.zpi.notifications.common.ConfidentialTextEncoder;

import java.time.LocalDateTime;
import java.util.List;

@Service
@RequiredArgsConstructor
@Slf4j
public class SlackService {

    private final SlackNotificationService slackNotificationService;
    private final SlackRepository slackRepository;
    private final ConfidentialTextEncoder confidentialTextEncoder;

    @SneakyThrows
    public void sendTestMessage(Long receiverSlackId) {
        var receiver = getSlackReceiver(receiverSlackId);
        String decodedWebhookUrl = confidentialTextEncoder.decrypt(receiver.getWebhookUrl());
        slackNotificationService.sendTestMessage(decodedWebhookUrl);
    }

    public List<SlackReceiver> getAllSlackIntegrations() {
        return slackRepository.findAll();
    }

    @SneakyThrows
    public void addNewSlackIntegration(SlackReceiverDTO slackIntegration) {
        String decodedWebhookUrl = confidentialTextEncoder.encrypt(slackIntegration.getWebhookUrl());
        checkIfWebhookExists(decodedWebhookUrl);
        SlackReceiver receiver = SlackReceiver.builder()
                .receiverName(slackIntegration.getName())
                .webhookUrl(decodedWebhookUrl)
                .createdAt(LocalDateTime.now())
                .build();
        slackRepository.save(receiver);
    }

    public SlackReceiver updateSlackIntegration(Long id, SlackReceiverDTO slackReceiver) {
        var receiver = getSlackReceiver(id);

        checkIfUserCanUpdateWebhookUrl(slackReceiver.getWebhookUrl(), id);

        receiver.setReceiverName(slackReceiver.getName());
        receiver.setWebhookUrl(slackReceiver.getWebhookUrl());
        receiver.setUpdatedAt(LocalDateTime.now());
        return slackRepository.save(receiver);
    }

    private SlackReceiver getSlackReceiver(Long receiverWebhookId) {
        return slackRepository.findById(receiverWebhookId)
                .orElseThrow(() -> new IllegalArgumentException("Webhook with given id not found"));
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
}
