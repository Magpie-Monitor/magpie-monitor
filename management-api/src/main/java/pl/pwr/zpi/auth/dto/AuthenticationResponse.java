package pl.pwr.zpi.auth.dto;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;
import pl.pwr.zpi.security.jwt.JwtToken;

@Data
@Builder
@AllArgsConstructor
@NoArgsConstructor
public class  AuthenticationResponse {

    private JwtToken token;

}
