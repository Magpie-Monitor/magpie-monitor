package pl.pwr.zpi.notifications.slack.dto;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Pattern;
import jakarta.validation.constraints.Size;
import lombok.Builder;
import lombok.Data;

@Data
@Builder
public class SlackReceiverDTO {
    @NotBlank(message = "Slack integration name cannot be empty")
    @Size(min = 2, max = 100, message = "The slack integration name must be from 2 to 100 characters.")
    private String name;
    @NotBlank(message = "Webhook url cannot be empty")
    @Pattern(regexp = "https://hooks.slack.com/services/[A-Z0-9]+/[A-Z0-9]+/[a-zA-Z0-9]+", message = "Provided webhook url is invalid")
    private String webhookUrl;
}
