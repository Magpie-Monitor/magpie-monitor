package pl.pwr.zpi.auth.oauth2;

import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.core.env.Environment;
import org.springframework.http.ResponseCookie;
import org.springframework.security.core.Authentication;
import org.springframework.security.oauth2.client.OAuth2AuthorizedClient;
import org.springframework.security.oauth2.client.OAuth2AuthorizedClientService;
import org.springframework.security.oauth2.client.authentication.OAuth2AuthenticationToken;
import org.springframework.security.oauth2.core.OAuth2AccessToken;
import org.springframework.security.oauth2.core.OAuth2RefreshToken;
import org.springframework.security.oauth2.core.oidc.user.DefaultOidcUser;
import org.springframework.security.web.authentication.SimpleUrlAuthenticationSuccessHandler;
import org.springframework.stereotype.Component;
import pl.pwr.zpi.security.cookie.CookieService;
import pl.pwr.zpi.user.data.User;
import pl.pwr.zpi.user.dto.Provider;
import pl.pwr.zpi.user.service.OAuthUserService;
import pl.pwr.zpi.user.service.UserService;
import pl.pwr.zpi.utils.AuthenticationUtils;

import java.io.IOException;


@Component
@RequiredArgsConstructor
@Slf4j
public class OAuthLoginSuccessHandler extends SimpleUrlAuthenticationSuccessHandler {

    @Value("${oauth2.google.redirect-uri}")
    private String REDIRECT_URI;

    public final OAuthUserService oAuthUserService;
    public final UserService userService;
    private final OAuth2AuthorizedClientService authorizedClientService;
    public final AuthenticationUtils authenticationUtils;
    private final Environment environment;
    private final CookieService cookieService;


    @Override
    public void onAuthenticationSuccess(HttpServletRequest request, HttpServletResponse response, Authentication authentication) throws IOException, ServletException {
        OAuth2AuthenticationToken oauthToken = (OAuth2AuthenticationToken) authentication;
        String registrationId = oauthToken.getAuthorizedClientRegistrationId();

        if (registrationId == null) {
            throw new ServletException("Failed to find the registrationId from the authorities");
        }
        OAuth2AuthorizedClient authorizedClient = authorizedClientService.loadAuthorizedClient(registrationId, authentication.getName());

        OAuth2AccessToken oAuth2AccessToken = authorizedClient.getAccessToken();
        OAuth2RefreshToken oAuth2RefreshToken = authorizedClient.getRefreshToken();

        User storedUser = authenticationUtils.getOAuthUserFromAuthentication(authentication);
        if (storedUser == null) {
            User newUser = User.builder()
                    .email(((DefaultOidcUser) authentication.getPrincipal()).getEmail())
                    .nickname(((DefaultOidcUser) authentication.getPrincipal()).getEmail().split("@")[0])
                    .provider(Provider.GOOGLE)
                    .authTokenExpDate(oAuth2AccessToken.getExpiresAt())
                    .build();

            userService.saveUser(newUser);
        } else {
            storedUser.setAuthTokenExpDate(oAuth2AccessToken.getExpiresAt());
            userService.saveUser(storedUser);
        }

        ResponseCookie authCookie = cookieService.createAuthCookie(oAuth2AccessToken.getTokenValue());
        response.addHeader("Set-Cookie", authCookie.toString());
        ResponseCookie refreshCookie = cookieService.createRefreshCookie(oAuth2RefreshToken.getTokenValue());
        response.addHeader("Set-Cookie", refreshCookie.toString());

        String targetUrl = REDIRECT_URI;
        getRedirectStrategy().sendRedirect(request, response, targetUrl);
    }
}