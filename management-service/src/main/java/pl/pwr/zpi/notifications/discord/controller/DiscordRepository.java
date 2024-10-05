package pl.pwr.zpi.notifications.discord.controller;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface DiscordRepository extends JpaRepository<DiscordReceiver, Long> {
    boolean existsByWebhookUrl(String webhookUrl);
}
