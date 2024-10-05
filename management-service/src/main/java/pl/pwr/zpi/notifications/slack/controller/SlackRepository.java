package pl.pwr.zpi.notifications.slack.controller;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface SlackRepository extends JpaRepository<SlackReceiver, Long> {
    boolean existsByWebhookUrl(String webhookUrl);
}
