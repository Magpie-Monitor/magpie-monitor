package pl.pwr.zpi.auth.oauth2;

import com.google.api.client.googleapis.auth.oauth2.GoogleRefreshTokenRequest;
import com.google.api.client.googleapis.auth.oauth2.GoogleTokenResponse;
import com.google.api.client.http.javanet.NetHttpTransport;
import com.google.api.client.json.gson.GsonFactory;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.ResponseCookie;
import org.springframework.security.core.Authentication;
import org.springframework.security.oauth2.client.OAuth2AuthorizedClient;
import org.springframework.security.oauth2.client.OAuth2AuthorizedClientService;
import org.springframework.security.oauth2.client.authentication.OAuth2AuthenticationToken;
import org.springframework.security.oauth2.core.OAuth2AccessToken;
import org.springframework.security.oauth2.core.OAuth2RefreshToken;
import org.springframework.security.oauth2.core.oidc.user.DefaultOidcUser;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.security.cookie.CookieService;
import pl.pwr.zpi.user.service.UserService;

import java.io.IOException;
import java.time.Instant;
import java.time.temporal.TemporalAmount;

@Service
@RequiredArgsConstructor
@Slf4j
public class OauthRefreshTokenService {

    @Value("${spring.security.oauth2.client.registration.google.client-id}")
    private String clientId;
    @Value("${spring.security.oauth2.client.registration.google.client-secret}")
    private String clientSecret;
    @Value("${google.oauth.cookie.exp-time}")
    private Long cookieExpTime;

    private final OAuth2AuthorizedClientService authorizedClientService;
    private final CookieService cookieService;


    public ResponseCookie updateAuthToken(Authentication authentication) {
        OAuth2AuthenticationToken oauthToken = (OAuth2AuthenticationToken) authentication;
        String registrationId = oauthToken.getAuthorizedClientRegistrationId();

        OAuth2AuthorizedClient authorizedClient = authorizedClientService.loadAuthorizedClient(registrationId, authentication.getName());

        OAuth2RefreshToken oAuth2RefreshToken = authorizedClient.getRefreshToken();

        if (oAuth2RefreshToken == null) {
            throw new RuntimeException("Refresh token is null");
        }

        return cookieService.createAuthCookie(refreshAccessToken(oAuth2RefreshToken.getTokenValue()), Instant.now().plusSeconds(cookieExpTime));
    }


    private String refreshAccessToken(String refreshToken) {
        GsonFactory jsonFactory = GsonFactory.getDefaultInstance();
        GoogleTokenResponse tokenResponse;
        try {
            tokenResponse = new GoogleRefreshTokenRequest(
                    new NetHttpTransport(),
                    jsonFactory,
                    refreshToken,
                    clientId,
                    clientSecret)
                    .execute();
            return tokenResponse.getAccessToken();

        } catch (IOException e) {
            throw new RuntimeException(e);
        }

    }
}
