package pl.pwr.zpi.auth.service;

import jakarta.servlet.http.Cookie;
import jakarta.servlet.http.HttpServletRequest;
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
        User oAuthUser = authenticationUtils.getOAuthUserFromAuthentication(authentication).orElseThrow(() ->
                new RuntimeException("User not found"));
        return UserDTO.toUserDTO(oAuthUser);
    }

    public TokenExpTime getTokenValidationTime(HttpServletRequest request) {
        Cookie[] cookies = request.getCookies();
        if (cookies == null) {
            throw new RuntimeException("No cookies found");
        }

        for (Cookie cookie : cookies) {
            if (cookie.getName().equals("authToken")) {
                String[] cookieParts = cookie.getValue().split("\\|");
                if (cookieParts.length == 2) {
                    long expiresAtMillis = Long.parseLong(cookieParts[1]);
                    long currentTimeMillis = Instant.now().toEpochMilli();
                    long timeDifferenceMillis = expiresAtMillis - currentTimeMillis;
                    return new TokenExpTime(timeDifferenceMillis);
                } else {
                    throw new RuntimeException("Invalid authToken format");
                }
            }
        }
        return new TokenExpTime(0L);
    }
}
