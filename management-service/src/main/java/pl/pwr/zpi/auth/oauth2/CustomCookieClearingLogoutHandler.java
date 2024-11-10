package pl.pwr.zpi.auth.oauth2;

import jakarta.servlet.http.Cookie;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.security.core.Authentication;
import org.springframework.security.web.authentication.logout.LogoutHandler;
import org.springframework.stereotype.Component;

@Component
public class CustomCookieClearingLogoutHandler implements LogoutHandler {
    @Value("${server.domainname}")
    private String PAGE_DOMAIN;
    @Value("${response-cookie.secure}")
    private boolean RESPONSE_COOKIE_SECURE;

    @Override
    public void logout(HttpServletRequest request, HttpServletResponse response, Authentication authentication) {
        clearCookie(response, "authToken", "/");
        clearCookie(response, "refreshToken", "/api/v1/auth/refresh-token");
        clearCookie(response, "JSESSIONID", "/");
    }

    private void clearCookie(HttpServletResponse response, String name, String path) {
        Cookie cookie = new Cookie(name, null);
        cookie.setPath(path);
        cookie.setDomain(PAGE_DOMAIN);
        cookie.setMaxAge(0);
        cookie.setSecure(RESPONSE_COOKIE_SECURE);
        cookie.setHttpOnly(true);
        response.addCookie(cookie);
    }
}
