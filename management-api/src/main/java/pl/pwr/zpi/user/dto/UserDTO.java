package pl.pwr.zpi.user.dto;

import lombok.Builder;
import lombok.Data;
import pl.pwr.zpi.user.data.User;

@Data
@Builder
public class UserDTO {
    private Long id;
    private String nickname;
    private String email;

    public static UserDTO toUserDTO(User user) {
        return UserDTO.builder()
                .id(user.getId())
                .nickname(user.getNickname())
                .email(user.getEmail())
                .build();
    }
}
