package pl.pwr.zpi.security.cookie;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.ResponseCookie;
import org.springframework.stereotype.Service;

@Service
public class CookieService {
    @Value("${server.domainname}")
    private String PAGE_DOMAIN;
    @Value("${response-cookie.secure}")
    private boolean RESPONSE_COOKIE_SECURE;

    public ResponseCookie createAuthCookie(String token) {
        return ResponseCookie.from("authToken", token)
                .httpOnly(true)
                .secure(RESPONSE_COOKIE_SECURE)
                .domain(PAGE_DOMAIN)
                .path("/")
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
