package pl.pwr.zpi.notifications.email.service;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.notifications.email.dto.EmailReceiverDTO;
import pl.pwr.zpi.notifications.email.entity.EmailReceiver;
import pl.pwr.zpi.notifications.email.repository.EmailRepository;

import java.util.List;

@Slf4j
@Service
@RequiredArgsConstructor
public class EmailReceiverService {

    private final EmailRepository emailRepository;

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

    public EmailReceiver getEmailReceiver(Long receiverEmailId) {
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
