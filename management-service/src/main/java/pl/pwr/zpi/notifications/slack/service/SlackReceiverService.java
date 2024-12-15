package pl.pwr.zpi.notifications.slack.service;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.notifications.common.ConfidentialTextEncoder;
import pl.pwr.zpi.notifications.slack.dto.SlackReceiverDTO;
import pl.pwr.zpi.notifications.slack.dto.UpdateSlackReceiverRequest;
import pl.pwr.zpi.notifications.slack.entity.SlackReceiver;
import pl.pwr.zpi.notifications.slack.repository.SlackRepository;

import java.util.List;
import java.util.regex.Pattern;

@Service
@RequiredArgsConstructor
@Slf4j
public class SlackReceiverService {

    private final SlackRepository slackRepository;
    private final ConfidentialTextEncoder confidentialTextEncoder;

    @Value("${slack.webhook.url.regex}")
    private String WEBHOOK_URL_REGEX;

    public List<SlackReceiver> getAllSlackIntegrations() {
        List<SlackReceiver> receivers = slackRepository.findAll();
        receivers.forEach(receiver -> receiver.setWebhookUrl(
                getAnonymizedWebhookUrl(receiver.getWebhookUrl()))
        );
        return receivers;
    }

    public void addNewSlackIntegration(SlackReceiverDTO slackIntegration) {
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

    public SlackReceiver updateSlackIntegration(Long id, UpdateSlackReceiverRequest updateRequest) {
        var receiver = getById(id);

        checkIfUserCanUpdateWebhookUrl(updateRequest.webhookUrl(), id);

        patchReceiver(receiver, updateRequest);
        receiver.setUpdatedAt(System.currentTimeMillis());

        return getAnonymizedSlackReceiver(receiver);
    }

    private SlackReceiver getAnonymizedSlackReceiver(SlackReceiver receiver) {
        receiver = slackRepository.save(receiver);
        receiver.setWebhookUrl(getAnonymizedWebhookUrl(receiver.getWebhookUrl()));
        return receiver;
    }

    private void patchReceiver(SlackReceiver slackReceiver, UpdateSlackReceiverRequest updateRequest) {
        if (updateRequest.name() != null) {
            validateReceiverName(updateRequest.name());
            slackReceiver.setReceiverName(updateRequest.name());
        }

        if (updateRequest.webhookUrl() != null) {
            validateWebhookUrl(updateRequest.webhookUrl());
            slackReceiver.setWebhookUrl(confidentialTextEncoder.encrypt(updateRequest.webhookUrl()));
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

    private String getAnonymizedWebhookUrl(String webhookUrl) {
        String decryptedUrl = confidentialTextEncoder.decrypt(webhookUrl);
        int lastSlashIndex = decryptedUrl.lastIndexOf('/');

        return decryptedUrl.substring(0, lastSlashIndex + 1) + "****";
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

    public SlackReceiver getEncodedWebhookUrl(Long id) {
        var receiver = getById(id);
        receiver.setWebhookUrl(confidentialTextEncoder.decrypt(receiver.getWebhookUrl()));
        return receiver;
    }

    public void deleteSlackReceiver(Long receiverId) {
        checkIfReceiverExist(receiverId);
        slackRepository.deleteById(receiverId);
    }

    private void checkIfReceiverExist(Long receiverId) {
        if (!slackRepository.existsById(receiverId)) {
            throw new IllegalArgumentException("Webhook with given Id not found");
        }
    }
}
