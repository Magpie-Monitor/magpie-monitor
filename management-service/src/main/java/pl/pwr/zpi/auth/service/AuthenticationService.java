package pl.pwr.zpi.auth.service;

import com.google.api.client.googleapis.auth.oauth2.GoogleIdToken.Payload;
import lombok.RequiredArgsConstructor;
import lombok.extern.log4j.Log4j2;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.auth.dto.TokenExpTime;
import pl.pwr.zpi.auth.oauth2.GoogleOauthTokenService;
import pl.pwr.zpi.user.dto.UserDTO;

@RequiredArgsConstructor
@Service
@Log4j2
public class AuthenticationService {

    private final GoogleOauthTokenService googleOauthTokenService;

    public UserDTO getUserDetails(String authToken) {
        Payload payload = googleOauthTokenService.decodeToken(authToken);
        return UserDTO.builder()
                .nickname((String) payload.get("name"))
                .email(payload.getEmail())
                .build();
    }

    public TokenExpTime getTokenValidationTime(String authToken) {
        long expMillis = googleOauthTokenService.decodeToken(authToken).getExpirationTimeSeconds() * 1000;
        return new TokenExpTime(expMillis - System.currentTimeMillis());
    }
}
