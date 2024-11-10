package pl.pwr.zpi.notifications.discord.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;
import pl.pwr.zpi.notifications.discord.entity.DiscordReceiver;

@Repository
public interface DiscordRepository extends JpaRepository<DiscordReceiver, Long> {
    boolean existsByWebhookUrl(String webhookUrl);
}
