package pl.pwr.zpi.auth.oauth2;

import jakarta.servlet.*;
import jakarta.servlet.http.HttpServletRequest;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Component;
import pl.pwr.zpi.security.cookie.CookieService;
import pl.pwr.zpi.utils.exception.AuthenticationException;
import pl.pwr.zpi.utils.jwt.JWTUtils;


@Component
@Slf4j
@RequiredArgsConstructor
public class OauthAuthenticator implements Filter {

    private final JWTUtils jwtUtils;
    private final CookieService cookieService;

    @Override
    public void doFilter(ServletRequest servletRequest, ServletResponse servletResponse, FilterChain filterChain) {
        HttpServletRequest httpRequest = (HttpServletRequest) servletRequest;
        String idToken = cookieService.getCookieValue(httpRequest, "authToken");

        if (idToken == null) {
            throw new AuthenticationException("Missing authentication token");
        }

        try {
            jwtUtils.decodeToken(idToken);
            filterChain.doFilter(servletRequest, servletResponse);

        } catch (Exception e) {
            throw new AuthenticationException("Token validation failed");
        }
    }
}
