package pl.pwr.zpi.notifications.slack.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;
import pl.pwr.zpi.notifications.slack.entity.SlackReceiver;

@Repository
public interface SlackRepository extends JpaRepository<SlackReceiver, Long> {
    boolean existsByWebhookUrl(String webhookUrl);
}
