package pl.pwr.zpi.auth.service;

import com.google.api.client.googleapis.auth.oauth2.GoogleIdToken.Payload;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.auth.dto.TokenExpTime;
import pl.pwr.zpi.user.dto.UserDTO;
import pl.pwr.zpi.utils.jwt.JWTUtils;

@RequiredArgsConstructor
@Service
public class AuthenticationService {

    private final JWTUtils jwtUtils;

    public UserDTO getUserDetails(String authToken) {
        Payload payload = jwtUtils.decodeToken(authToken);
        return UserDTO.builder()
                .nickname((String) payload.get("name"))
                .email(payload.getEmail())
                .build();
    }

    public TokenExpTime getTokenValidationTime(String authToken) {
        long expMillis = jwtUtils.decodeToken(authToken).getExpirationTimeSeconds() * 1000;
        return new TokenExpTime(expMillis - System.currentTimeMillis());
    }
}
