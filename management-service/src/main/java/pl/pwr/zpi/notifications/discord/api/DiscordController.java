package pl.pwr.zpi.notifications.discord.api;

import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import pl.pwr.zpi.notifications.discord.DiscordNotificationService;
import pl.pwr.zpi.notifications.discord.entity.DiscordReceiver;
import pl.pwr.zpi.notifications.discord.dto.DiscordReceiverDTO;
import pl.pwr.zpi.notifications.discord.service.DiscordReceiverService;

import java.util.List;

@RequiredArgsConstructor
@RestController
@RequestMapping("/api/v1/notification-channels/discord")
public class DiscordController {
    private final DiscordReceiverService discordReceiverService;
    private final DiscordNotificationService discordNotificationService;

    @GetMapping
    public ResponseEntity<List<DiscordReceiver>> getAllDiscordIntegrations() {
        return ResponseEntity.ok().body(discordReceiverService.getAllDiscordIntegrations());
    }

    @PostMapping
    public ResponseEntity<DiscordReceiver> addDiscordReceiver(@Valid @RequestBody DiscordReceiverDTO discordReceiver) throws Exception {
        discordReceiverService.createDiscordReceiver(discordReceiver);
        return ResponseEntity.ok().build();
    }

    @PatchMapping("/{id}")
    public ResponseEntity<DiscordReceiver> updateDiscordIntegration(@PathVariable Long id, @Valid @RequestBody DiscordReceiverDTO DiscordReceiver) throws Exception {
        return ResponseEntity.ok().body(discordReceiverService.updateDiscordIntegration(id, DiscordReceiver));
    }

    @PostMapping("/{id}/test-notification")
    public ResponseEntity<DiscordReceiver> sendTestMessage(@PathVariable Long id) throws Exception {
        discordNotificationService.sendTestMessage(id);
        return ResponseEntity.ok().build();
    }

    @DeleteMapping("/{id}")
    public ResponseEntity<Void> deleteDiscordReceiver(@PathVariable Long id) {
        discordReceiverService.deleteDiscordReceiver(id);
        return ResponseEntity.ok().build();
    }
}
