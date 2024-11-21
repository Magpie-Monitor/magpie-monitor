package pl.pwr.zpi.notifications.email.service;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.apache.commons.validator.routines.EmailValidator;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.notifications.email.dto.EmailReceiverDTO;
import pl.pwr.zpi.notifications.email.dto.EmailReceiverUpdateRequest;
import pl.pwr.zpi.notifications.email.entity.EmailReceiver;
import pl.pwr.zpi.notifications.email.repository.EmailRepository;

import java.util.List;

@Slf4j
@Service
@RequiredArgsConstructor
public class EmailReceiverService {

    private final EmailRepository emailRepository;
    private final EmailValidator emailValidator = EmailValidator.getInstance();

    public List<EmailReceiver> getAllEmails() {
        return emailRepository.findAll();
    }

    public void addNewEmail(EmailReceiverDTO emailReceiver) {
        long now = System.currentTimeMillis();

        checkIfEmailExists(emailReceiver);
        EmailReceiver receiver = EmailReceiver.builder()
                .receiverName(emailReceiver.getName())
                .receiverEmail(emailReceiver.getEmail())
                .createdAt(now)
                .updatedAt(now)
                .build();
        emailRepository.save(receiver);
    }

    public EmailReceiver updateEmail(Long id, EmailReceiverUpdateRequest updateRequest) {
        var receiver = getEmailReceiver(id);

        checkIfUserCanUpdateEmail(updateRequest.email(), id);
        patchEmail(receiver, updateRequest);

        return emailRepository.save(receiver);
    }

    private void patchEmail(EmailReceiver emailReceiver, EmailReceiverUpdateRequest updateRequest) {
        if(updateRequest.email() != null) {
            validateEmail(updateRequest.email());
            emailReceiver.setReceiverEmail(updateRequest.email());
        }

        if(updateRequest.name() != null) {
            validateReceiverName(updateRequest.name());
            emailReceiver.setReceiverName(updateRequest.name());
        }

        emailReceiver.setUpdatedAt(System.currentTimeMillis());
    }

    private void validateEmail(String email) {
        if(!emailValidator.isValid(email)) {
            throw new RuntimeException("Invalid email");
        }
    }

    private void validateReceiverName(String name) {
        if(name.isEmpty()) {
           throw new RuntimeException("Receiver name is empty");
        }
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

    public void deleteEmailReceiver(Long receiverId) {
        emailRepository.deleteById(receiverId);
    }
}
