package pl.pwr.zpi.notifications.email.dto;

public record EmailReceiverUpdateRequest(
        String name,
        String email
) {
}
