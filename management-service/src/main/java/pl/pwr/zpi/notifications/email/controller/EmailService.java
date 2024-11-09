package pl.pwr.zpi.notifications.email.controller;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.notifications.email.EmailNotificationService;

import java.util.List;

@Service
@RequiredArgsConstructor
@Slf4j
public class EmailService {

    private final EmailNotificationService emailNotificationService;
    private final EmailRepository emailRepository;

    public void sendTestEmail(Long receiverEmailId) {
        var receiver = getEmailReceiver(receiverEmailId);
        log.info("Sending test email to: {}", receiver.getReceiverEmail());
        emailNotificationService.sendTestEmail(receiver.getReceiverEmail());
    }

    public List<EmailReceiver> getAllEmails() {
        return emailRepository.findAll();
    }

    public void addNewEmail(EmailReceiverDTO emailReceiver) {
        checkIfEmailExists(emailReceiver);
        EmailReceiver receiver = EmailReceiver.builder()
                .receiverName(emailReceiver.getName())
                .receiverEmail(emailReceiver.getEmail())
                .createdAt(System.currentTimeMillis())
                .build();
        emailRepository.save(receiver);
    }

    public EmailReceiver updateEmail(Long id, EmailReceiverDTO emailReceiver) {
        var receiver = getEmailReceiver(id);

        checkIfUserCanUpdateEmail(emailReceiver.getEmail(), id);

        receiver.setReceiverName(emailReceiver.getName());
        receiver.setReceiverEmail(emailReceiver.getEmail());
        receiver.setUpdatedAt(System.currentTimeMillis());
        return emailRepository.save(receiver);
    }

    private EmailReceiver getEmailReceiver(Long receiverEmailId) {
        return emailRepository.findById(receiverEmailId)
                .orElseThrow(() -> new IllegalArgumentException("Email with given clusterId not found"));
    }

    private void checkIfEmailExists(EmailReceiverDTO emailReceiver) {
        if (emailRepository.existsByReceiverEmail(emailReceiver.getEmail())) {
            throw new IllegalArgumentException("Email already exists");
        }
    }

    private void checkIfUserCanUpdateEmail(String email, Long id) {
        if (emailRepository.existsByReceiverEmail(email) && !emailRepository.findById(id).get().getReceiverEmail().equals(email)) {
            throw new IllegalArgumentException("Email is already assigned to other entry");
        }
    }
}
