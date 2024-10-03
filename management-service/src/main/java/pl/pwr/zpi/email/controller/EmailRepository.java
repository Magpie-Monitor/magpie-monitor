package pl.pwr.zpi.email.controller;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface EmailRepository extends JpaRepository<EmailReceiver, Long> {
    boolean existsByReceiverEmail(String email);
}
