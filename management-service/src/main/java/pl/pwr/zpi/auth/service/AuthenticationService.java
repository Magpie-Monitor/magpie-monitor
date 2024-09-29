package pl.pwr.zpi.auth.service;

import lombok.RequiredArgsConstructor;
import org.springframework.security.core.Authentication;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.auth.dto.TokenExpTime;
import pl.pwr.zpi.user.data.User;
import pl.pwr.zpi.user.dto.UserDTO;
import pl.pwr.zpi.utils.AuthenticationUtils;

import java.time.Instant;

@RequiredArgsConstructor
@Service
public class AuthenticationService {

    private final AuthenticationUtils authenticationUtils;

    public UserDTO getUserDetails(Authentication authentication) {
        User oAuthUser = authenticationUtils.getOAuthUserFromAuthentication(authentication);
        return UserDTO.toUserDTO(oAuthUser);
    }

    public TokenExpTime getTokenValidationTime(Authentication authentication) {
        User oAuthUser = authenticationUtils.getOAuthUserFromAuthentication(authentication);
        long tokenExpiration = oAuthUser.getAuthTokenExpDate().toEpochMilli();
        long currentTimeMillis = Instant.now().toEpochMilli();
        long timeDifferenceMillis = tokenExpiration - currentTimeMillis;
        return new TokenExpTime(timeDifferenceMillis);
    }
}
