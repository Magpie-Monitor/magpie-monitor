package pl.pwr.zpi.notifications.slack.api;

import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import pl.pwr.zpi.notifications.slack.SlackNotificationService;
import pl.pwr.zpi.notifications.slack.dto.SlackReceiverDTO;
import pl.pwr.zpi.notifications.slack.dto.UpdateSlackReceiverRequest;
import pl.pwr.zpi.notifications.slack.entity.SlackReceiver;
import pl.pwr.zpi.notifications.slack.service.SlackReceiverService;

import java.util.List;

@RequiredArgsConstructor
@RestController
@RequestMapping("/api/v1/notification-channels/slack")
public class SlackController {

    private final SlackReceiverService slackReceiverService;
    private final SlackNotificationService slackNotificationService;

    @GetMapping
    public ResponseEntity<List<SlackReceiver>> getAllSlackIntegrations() {
        return ResponseEntity.ok().body(slackReceiverService.getAllSlackIntegrations());
    }

    @PostMapping
    public ResponseEntity<SlackReceiver> addSlackIntegration(@Valid @RequestBody SlackReceiverDTO slackReceiver) {
        slackReceiverService.addNewSlackIntegration(slackReceiver);
        return ResponseEntity.ok().build();
    }

    @PatchMapping("/{id}")
    public ResponseEntity<SlackReceiver> updateSlackIntegration(
            @PathVariable Long id, @Valid @RequestBody UpdateSlackReceiverRequest updateRequest) {
        return ResponseEntity.ok().body(slackReceiverService.updateSlackIntegration(id, updateRequest));
    }

    @PostMapping("/{id}/test-notification")
    public ResponseEntity<SlackReceiver> sendTestMessage(@PathVariable Long id) {
        slackNotificationService.sendTestMessage(id);
        return ResponseEntity.ok().build();
    }

    @DeleteMapping("/{id}")
    public ResponseEntity<Void> deleteSlackReceiver(@PathVariable Long id) {
        slackReceiverService.deleteSlackReceiver(id);
        return ResponseEntity.ok().build();
    }
}
