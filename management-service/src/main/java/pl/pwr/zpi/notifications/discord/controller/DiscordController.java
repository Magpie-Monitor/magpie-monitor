package pl.pwr.zpi.notifications.discord.controller;

import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RequiredArgsConstructor
@RestController
@RequestMapping("/api/v1/notification-channels/discord")
public class DiscordController {
    private final DiscordService discordService;

    @GetMapping
    public ResponseEntity<List<DiscordReceiver>> getAllDiscordIntegrations() {
        return ResponseEntity.ok().body(discordService.getAllDiscordIntegrations());
    }

    @PostMapping
    public ResponseEntity<DiscordReceiver> addDiscordIntegration(@Valid @RequestBody DiscordReceiverDTO discordReceiver) throws Exception {
        discordService.addNewDiscordIntegration(discordReceiver);
        return ResponseEntity.ok().build();
    }

    @PatchMapping("/{id}")
    public ResponseEntity<DiscordReceiver> updateDiscordIntegration(@PathVariable Long id, @Valid @RequestBody DiscordReceiverDTO DiscordReceiver) throws Exception {
        return ResponseEntity.ok().body(discordService.updateDiscordIntegration(id, DiscordReceiver));
    }

    @PostMapping("/{id}/test-notification")
    public ResponseEntity<DiscordReceiver> sendTestMessage(@PathVariable Long id) throws Exception {
        discordService.sendTestMessage(id);
        return ResponseEntity.ok().build();
    }

    @GetMapping("/{id}/webhook-url")
    public ResponseEntity<DiscordReceiver> getWebhookUrl(@PathVariable Long id) throws Exception {
        return ResponseEntity.ok().body(discordService.getEncodedWebhookUrl(id));
    }
}
