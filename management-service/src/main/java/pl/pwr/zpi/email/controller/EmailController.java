package pl.pwr.zpi.email.controller;

import lombok.RequiredArgsConstructor;
import org.springframework.web.bind.annotation.*;

@RequiredArgsConstructor
@RestController
@RequestMapping("/api/v1/notification-channels/mails")
public class EmailController {
    private final EmailService emailService;

    @GetMapping("/test-notification")
    public void sendTestEmail(@RequestParam String receiverEmail) {
        emailService.sendTestEmail(receiverEmail);
    }
}

