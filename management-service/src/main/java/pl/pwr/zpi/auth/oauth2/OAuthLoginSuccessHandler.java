package pl.pwr.zpi.auth.oauth2;

import com.google.api.client.googleapis.auth.oauth2.GoogleIdToken.Payload;
import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.ResponseCookie;
import org.springframework.security.core.Authentication;
import org.springframework.security.oauth2.client.OAuth2AuthorizedClient;
import org.springframework.security.oauth2.client.OAuth2AuthorizedClientService;
import org.springframework.security.oauth2.client.authentication.OAuth2AuthenticationToken;
import org.springframework.security.oauth2.core.OAuth2RefreshToken;
import org.springframework.security.oauth2.core.oidc.user.DefaultOidcUser;
import org.springframework.security.web.authentication.SimpleUrlAuthenticationSuccessHandler;
import org.springframework.stereotype.Component;
import pl.pwr.zpi.security.cookie.CookieService;
import pl.pwr.zpi.user.data.User;
import pl.pwr.zpi.user.dto.Provider;
import pl.pwr.zpi.user.service.UserService;
import pl.pwr.zpi.utils.jwt.JWTUtils;

import java.io.IOException;
import java.time.Instant;
import java.util.Optional;

@Component
@RequiredArgsConstructor
@Slf4j
public class OAuthLoginSuccessHandler extends SimpleUrlAuthenticationSuccessHandler {

    @Value("${oauth2.google.redirect-uri}")
    private String redirectUri;

    private final UserService userService;
    private final OAuth2AuthorizedClientService authorizedClientService;
    private final CookieService cookieService;
    private final JWTUtils jwtUtils;

    @Override
    public void onAuthenticationSuccess(HttpServletRequest request, HttpServletResponse response, Authentication authentication) throws IOException, ServletException {
        if (!(authentication instanceof OAuth2AuthenticationToken)) {
            throw new ServletException("Unsupported authentication type.");
        }

        OAuth2AuthenticationToken oauthToken = (OAuth2AuthenticationToken) authentication;
        String registrationId = oauthToken.getAuthorizedClientRegistrationId();
        if (registrationId == null) {
            throw new ServletException("Missing registration ID.");
        }

        OAuth2AuthorizedClient authorizedClient = authorizedClientService.loadAuthorizedClient(registrationId, authentication.getName());
        DefaultOidcUser oidcUser = (DefaultOidcUser) authentication.getPrincipal();
        OAuth2RefreshToken refreshToken = authorizedClient.getRefreshToken();

        if (refreshToken == null) {
            throw new RuntimeException("Refresh token is null. Please contact the administrator.");
        }

        addCookiesToResponse(response, oidcUser.getIdToken().getTokenValue(), oidcUser.getIdToken().getExpiresAt(), refreshToken.getTokenValue());

        createOrUpdateUser(oidcUser);

        getRedirectStrategy().sendRedirect(request, response, redirectUri);
    }

    private void addCookiesToResponse(HttpServletResponse response, String authToken, Instant expiresAt, String refreshToken) {
        ResponseCookie authCookie = cookieService.createAuthCookie(authToken, expiresAt);
        ResponseCookie refreshCookie = cookieService.createRefreshCookie(refreshToken);

        response.addHeader("Set-Cookie", authCookie.toString());
        response.addHeader("Set-Cookie", refreshCookie.toString());
    }

    private void createOrUpdateUser(DefaultOidcUser oidcUser) {
        Payload payload = jwtUtils.decodeToken(oidcUser.getIdToken().getTokenValue());

        String email = payload.getEmail();
        String nickname = (String) payload.get("name");
        String fallbackNickname = nickname != null ? nickname : email.split("@")[0];

        userService.findByEmail(email)
                .ifPresentOrElse(
                        user -> log.info("User exists: {}", email),
                        () -> userService.saveUser(User.builder()
                                .email(email)
                                .nickname(fallbackNickname)
                                .provider(Provider.GOOGLE)
                                .build())
                );
    }
}
