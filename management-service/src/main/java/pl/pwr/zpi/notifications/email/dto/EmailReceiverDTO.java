package pl.pwr.zpi.notifications.email.dto;

import jakarta.validation.constraints.Email;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Size;
import lombok.Builder;
import lombok.Data;

import java.util.Optional;

@Data
@Builder
public class EmailReceiverDTO {
//    @NotBlank(message = "Name cannot be empty")
//    @Size(min = 2, max = 100, message = "The email name must be from 2 to 100 characters.")
    private String name;
////    @NotBlank(message = "Email cannot be empty")
//    @Email(message = "Provided email is invalid")
    private String email;
}
