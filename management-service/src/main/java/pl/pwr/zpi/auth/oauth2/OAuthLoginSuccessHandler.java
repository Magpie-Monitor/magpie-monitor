package pl.pwr.zpi.auth.oauth2;

import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.ResponseCookie;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.authority.SimpleGrantedAuthority;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.security.oauth2.client.OAuth2AuthorizedClient;
import org.springframework.security.oauth2.client.OAuth2AuthorizedClientService;
import org.springframework.security.oauth2.client.authentication.OAuth2AuthenticationToken;
import org.springframework.security.oauth2.core.OAuth2AccessToken;
import org.springframework.security.oauth2.core.OAuth2RefreshToken;
import org.springframework.security.oauth2.core.oidc.user.DefaultOidcUser;
import org.springframework.security.oauth2.core.user.OAuth2User;
import org.springframework.security.web.authentication.SimpleUrlAuthenticationSuccessHandler;
import org.springframework.stereotype.Component;
import pl.pwr.zpi.security.cookie.CookieService;
import pl.pwr.zpi.user.data.User;
import pl.pwr.zpi.user.dto.Provider;
import pl.pwr.zpi.user.service.UserService;
import pl.pwr.zpi.utils.AuthenticationUtils;

import java.io.IOException;
import java.util.ArrayList;
import java.util.List;
import java.util.stream.Collectors;


@Component
@RequiredArgsConstructor
@Slf4j
public class OAuthLoginSuccessHandler extends SimpleUrlAuthenticationSuccessHandler {

    @Value("${oauth2.google.redirect-uri}")
    private String REDIRECT_URI;

    private final UserService userService;
    private final OAuth2AuthorizedClientService authorizedClientService;
    private final AuthenticationUtils authenticationUtils;
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

        createOrUpdateUser(authentication);

        ResponseCookie authCookie = cookieService.createAuthCookie(oAuth2AccessToken.getTokenValue(), oAuth2AccessToken.getExpiresAt());
        response.addHeader("Set-Cookie", authCookie.toString());
        ResponseCookie refreshCookie = cookieService.createRefreshCookie(oAuth2RefreshToken.getTokenValue());
        response.addHeader("Set-Cookie", refreshCookie.toString());

        OAuth2User oAuth2User = (OAuth2User) authentication.getPrincipal();

        List<GrantedAuthority> updatedAuthorities = new ArrayList<>(authentication.getAuthorities());
        Authentication updatedAuthentication = new OAuth2AuthenticationToken(
                oAuth2User,
                updatedAuthorities,
                registrationId
        );

        SecurityContextHolder.getContext().setAuthentication(updatedAuthentication);

        String targetUrl = REDIRECT_URI;
        getRedirectStrategy().sendRedirect(request, response, targetUrl);
    }

    private void createOrUpdateUser(Authentication authentication) {
        authenticationUtils.getOAuthUserFromAuthentication(authentication)
                .ifPresentOrElse(
                        (storedUser) -> {},
                        () -> userService.saveUser(User.builder()
                                .email(((DefaultOidcUser) authentication.getPrincipal()).getEmail())
                                .nickname(((DefaultOidcUser) authentication.getPrincipal()).getEmail().split("@")[0])
                                .provider(Provider.GOOGLE)
                                .build())
                );
    }

}