package pl.pwr.zpi.notifications.email.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;
import pl.pwr.zpi.notifications.email.entity.EmailReceiver;

@Repository
public interface EmailRepository extends JpaRepository<EmailReceiver, Long> {
    boolean existsByReceiverEmail(String email);
}
