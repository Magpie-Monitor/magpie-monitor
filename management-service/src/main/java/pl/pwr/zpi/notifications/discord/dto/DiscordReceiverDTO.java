package pl.pwr.zpi.notifications.discord.dto;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Pattern;
import jakarta.validation.constraints.Size;
import lombok.Builder;
import lombok.Data;

@Data
@Builder
public class DiscordReceiverDTO {
    @NotBlank(message = "Discord integration name cannot be empty")
    @Size(min = 2, max = 100, message = "The discord integration name must be from 2 to 100 characters.")
    private String name;
    @NotBlank(message = "Webhook url cannot be empty")
    @Pattern(regexp = "https://discord.com/api/webhooks/[0-9]+/[a-zA-Z0-9\\-_]+", message = "Provided webhook url is invalid")
    private String webhookUrl;
}