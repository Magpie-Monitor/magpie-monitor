package pl.pwr.zpi.email.controller;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.email.NotificationEmailService;

@Service
@RequiredArgsConstructor
@Slf4j
public class EmailService {

    private final NotificationEmailService notificationEmailService;

    public void sendTestEmail(String receiverEmail) {
        log.info("Sending test email to: {}", receiverEmail);
        notificationEmailService.sendTestEmail(receiverEmail);
    }
}
