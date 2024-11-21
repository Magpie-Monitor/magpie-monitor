package pl.pwr.zpi.notifications.email.api;

import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import pl.pwr.zpi.notifications.email.EmailNotificationService;
import pl.pwr.zpi.notifications.email.dto.EmailReceiverDTO;
import pl.pwr.zpi.notifications.email.dto.EmailReceiverUpdateRequest;
import pl.pwr.zpi.notifications.email.entity.EmailReceiver;
import pl.pwr.zpi.notifications.email.service.EmailReceiverService;
import pl.pwr.zpi.notifications.slack.entity.SlackReceiver;

import java.util.List;

@RequiredArgsConstructor
@RestController
@RequestMapping("/api/v1/notification-channels/mails")
public class EmailController {

    private final EmailReceiverService emailReceiverService;
    private final EmailNotificationService emailNotificationService;

    @GetMapping
    public ResponseEntity<List<EmailReceiver>> getAllEmails() {
        return ResponseEntity.ok().body(emailReceiverService.getAllEmails());
    }

    @PostMapping
    public ResponseEntity<EmailReceiver> addEmail(@Valid @RequestBody EmailReceiverDTO emailReceiver) {
        emailReceiverService.addNewEmail(emailReceiver);
        return ResponseEntity.ok().build();
    }

    @PatchMapping("/{id}")
    public ResponseEntity<EmailReceiver> updateEmail(
            @PathVariable Long id, @RequestBody EmailReceiverUpdateRequest updateRequest) {
        return ResponseEntity.ok().body(emailReceiverService.updateEmail(id, updateRequest));
    }

    @PostMapping("/{id}/test-notification")
    public ResponseEntity<SlackReceiver> sendTestEmail(@PathVariable Long id) {
        emailNotificationService.sendTestEmail(id);
        return ResponseEntity.ok().build();
    }

    @DeleteMapping("/{id}")
    public ResponseEntity<Void> deleteEmailReceiver(@PathVariable Long id) {
        emailReceiverService.deleteEmailReceiver(id);
        return ResponseEntity.ok().build();
    }
}
