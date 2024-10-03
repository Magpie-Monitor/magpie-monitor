package pl.pwr.zpi.email.controller;

import jakarta.transaction.Transactional;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.email.NotificationEmailService;

import java.time.LocalDateTime;
import java.util.List;

@Service
@RequiredArgsConstructor
@Slf4j
public class EmailService {

    private final NotificationEmailService notificationEmailService;
    private final EmailRepository emailRepository;

    public void sendTestEmail(Long receiverEmailId) {
        var receiver = getEmailReceiver(receiverEmailId);
        log.info("Sending test email to: {}", receiver.getReceiverEmail());
        notificationEmailService.sendTestEmail(receiver.getReceiverEmail());
    }

    public List<EmailReceiver> getAllEmails() {
        return emailRepository.findAll();
    }

    public void addNewEmail(EmailReceiverDTO emailReceiver) {
        checkIfEmailExists(emailReceiver);
        EmailReceiver receiver = EmailReceiver.builder()
                .receiverName(emailReceiver.getName())
                .receiverEmail(emailReceiver.getEmail())
                .createdAt(LocalDateTime.now())
                .build();
        emailRepository.save(receiver);
    }

    public EmailReceiver updateEmail(Long id, EmailReceiverDTO emailReceiver) {
        var receiver = getEmailReceiver(id);

        checkIfUserCanUpdateEmail(emailReceiver.getEmail(), id);

        receiver.setReceiverName(emailReceiver.getName());
        receiver.setReceiverEmail(emailReceiver.getEmail());
        receiver.setUpdatedAt(LocalDateTime.now());
        return emailRepository.save(receiver);
    }

    private EmailReceiver getEmailReceiver(Long receiverEmailId) {
        return emailRepository.findById(receiverEmailId)
                .orElseThrow(() -> new IllegalArgumentException("Email with given id not found"));
    }

    private void checkIfEmailExists(EmailReceiverDTO emailReceiver) {
        if (emailRepository.existsByReceiverEmail(emailReceiver.getEmail())) {
            throw new IllegalArgumentException("Email already exists");
        }
    }

    private void checkIfUserCanUpdateEmail(String email, Long id) {
        if (emailRepository.existsByReceiverEmail(email) && !emailRepository.findById(id).get().getReceiverEmail().equals(email)) {
            throw new IllegalArgumentException("Email is already assigned to other user");
        }
    }
}
