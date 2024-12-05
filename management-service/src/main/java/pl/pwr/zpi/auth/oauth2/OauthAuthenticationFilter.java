package pl.pwr.zpi.auth.oauth2;

import jakarta.servlet.FilterChain;
import jakarta.servlet.ServletException;
import jakarta.servlet.http.Cookie;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import org.springframework.http.HttpStatus;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.web.filter.OncePerRequestFilter;
import pl.pwr.zpi.utils.exception.AuthenticationException;

import java.io.IOException;
import java.security.GeneralSecurityException;

public class OauthAuthenticationFilter extends OncePerRequestFilter {

    private final GoogleOauthTokenService googleOauthTokenService;

    public OauthAuthenticationFilter(GoogleOauthTokenService googleOauthTokenService) {
        this.googleOauthTokenService = googleOauthTokenService;
    }

    @Override
    protected void doFilterInternal(HttpServletRequest request, HttpServletResponse response, FilterChain filterChain)
            throws ServletException, IOException {
        try {
            String authToken = getAuthTokenFromCookie(request);
            if (authToken != null) {
                googleOauthTokenService.validateToken(authToken);

                Authentication authentication = new UsernamePasswordAuthenticationToken(authToken, null, null);
                SecurityContextHolder.getContext().setAuthentication(authentication);
            }
        } catch (AuthenticationException e) {
            response.setStatus(HttpStatus.UNAUTHORIZED.value());
            response.getWriter().write("Authentication failed: " + e.getMessage());
            return;
        } catch (GeneralSecurityException e) {
            throw new AuthenticationException("Error while validating token");
        }

        filterChain.doFilter(request, response);
    }

    private String getAuthTokenFromCookie(HttpServletRequest request) {
        if (request.getCookies() != null) {
            for (Cookie cookie : request.getCookies()) {
                if ("authToken".equals(cookie.getName())) {
                    return cookie.getValue();
                }
            }
        }
        return null;
    }
}
