package pl.pwr.zpi.email.controller;

import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RequiredArgsConstructor
@RestController
@RequestMapping("/api/v1/notification-channels/mails")
public class EmailController {
    private final EmailService emailService;

    @GetMapping
    public ResponseEntity<List<EmailReceiver>> getAllEmails() {
        return ResponseEntity.ok().body(emailService.getAllEmails());
    }

    @PostMapping
    public ResponseEntity<EmailReceiver> addEmail(@Valid @RequestBody EmailReceiverDTO emailReceiver) {
        emailService.addNewEmail(emailReceiver);
        return ResponseEntity.ok().build();
    }

    @PatchMapping("/{id}")
    public ResponseEntity<EmailReceiver> updateEmail(@PathVariable Long id, @Valid @RequestBody EmailReceiverDTO emailReceiver) {
        emailService.updateEmail(id, emailReceiver);
        return ResponseEntity.ok().build();
    }

    @GetMapping("/{id}/test-notification")
    public void sendTestEmail(@PathVariable Long id) {
        emailService.sendTestEmail(id);
    }
}
