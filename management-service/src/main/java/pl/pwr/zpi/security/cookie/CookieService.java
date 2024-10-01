package pl.pwr.zpi.security.cookie;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.ResponseCookie;
import org.springframework.stereotype.Service;

import java.time.Duration;
import java.time.Instant;

@Service
public class CookieService {
    @Value("${server.domainname}")
    private String PAGE_DOMAIN;
    @Value("${response-cookie.secure}")
    private boolean RESPONSE_COOKIE_SECURE;

    public ResponseCookie createAuthCookie(String token, Instant expiresAt) {
        long maxAgeInSeconds = Duration.between(Instant.now(), expiresAt).getSeconds();

        String cookieValue = token + "|" + expiresAt.toEpochMilli();

        return ResponseCookie.from("authToken", cookieValue)
                .httpOnly(true)
                .secure(RESPONSE_COOKIE_SECURE)
                .domain(PAGE_DOMAIN)
                .path("/")
                .maxAge(maxAgeInSeconds)
                .build();
    }


    public ResponseCookie createRefreshCookie(String token) {
        return ResponseCookie.from("refreshToken", token)
                .httpOnly(true)
                .secure(RESPONSE_COOKIE_SECURE)
                .domain(PAGE_DOMAIN)
                .path("/") //api/v1/auth/refresh
                .build();
    }

}
