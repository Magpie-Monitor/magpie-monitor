package pl.pwr.zpi.notifications.slack.controller;

import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RequiredArgsConstructor
@RestController
@RequestMapping("/api/v1/notification-channels/slack")
public class SlackController {
    private final SlackService slackService;

    @GetMapping
    public ResponseEntity<List<SlackReceiver>> getAllSlackIntegrations() {
        return ResponseEntity.ok().body(slackService.getAllSlackIntegrations());
    }

    @PostMapping
    public ResponseEntity<SlackReceiver> addSlackIntegration(@Valid @RequestBody SlackReceiverDTO slackReceiver) throws Exception {
        slackService.addNewSlackIntegration(slackReceiver);
        return ResponseEntity.ok().build();
    }

    @PatchMapping("/{id}")
    public ResponseEntity<SlackReceiver> updateSlackIntegration(@PathVariable Long id, @Valid @RequestBody SlackReceiverDTO slackReceiver) throws Exception {
        return ResponseEntity.ok().body(slackService.updateSlackIntegration(id, slackReceiver));
    }

    @GetMapping("/{id}/test-notification")
    public ResponseEntity<SlackReceiver> sendTestMessage(@PathVariable Long id) throws Exception {
        slackService.sendTestMessage(id);
        return ResponseEntity.ok().build();
    }
}
