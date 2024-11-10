package pl.pwr.zpi.notifications.slack.api;

import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import pl.pwr.zpi.notifications.slack.dto.SlackReceiverDTO;
import pl.pwr.zpi.notifications.slack.entity.SlackReceiver;
import pl.pwr.zpi.notifications.slack.service.SlackReceiverService;

import java.util.List;

@RequiredArgsConstructor
@RestController
@RequestMapping("/api/v1/notification-channels/slack")
public class SlackController {
    private final SlackReceiverService slackReceiverService;

    @GetMapping
    public ResponseEntity<List<SlackReceiver>> getAllSlackIntegrations() {
        return ResponseEntity.ok().body(slackReceiverService.getAllSlackIntegrations());
    }

    @PostMapping
    public ResponseEntity<SlackReceiver> addSlackIntegration(@Valid @RequestBody SlackReceiverDTO slackReceiver) throws Exception {
        slackReceiverService.addNewSlackIntegration(slackReceiver);
        return ResponseEntity.ok().build();
    }

    @PatchMapping("/{id}")
    public ResponseEntity<SlackReceiver> updateSlackIntegration(@PathVariable Long id, @Valid @RequestBody SlackReceiverDTO slackReceiver) throws Exception {
        return ResponseEntity.ok().body(slackReceiverService.updateSlackIntegration(id, slackReceiver));
    }

    @PostMapping("/{id}/test-notification")
    public ResponseEntity<SlackReceiver> sendTestMessage(@PathVariable Long id) throws Exception {
//        slackReceiverService.sendTestMessage(id);
        return ResponseEntity.ok().build();
    }

    @GetMapping("/{id}/webhook-url")
    public ResponseEntity<SlackReceiver> getWebhookUrl(@PathVariable Long id) throws Exception {
        return ResponseEntity.ok().body(slackReceiverService.getEncodedWebhookUrl(id));
    }
}
