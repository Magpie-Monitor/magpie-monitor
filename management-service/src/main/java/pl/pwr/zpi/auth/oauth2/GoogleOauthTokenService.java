package pl.pwr.zpi.auth.oauth2;

import com.google.api.client.googleapis.auth.oauth2.GoogleIdToken;
import com.google.api.client.googleapis.auth.oauth2.GoogleIdTokenVerifier;
import com.google.api.client.googleapis.auth.oauth2.GoogleRefreshTokenRequest;
import com.google.api.client.googleapis.auth.oauth2.GoogleTokenResponse;
import com.google.api.client.http.javanet.NetHttpTransport;
import com.google.api.client.json.gson.GsonFactory;
import lombok.RequiredArgsConstructor;
import lombok.SneakyThrows;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.ResponseCookie;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.security.cookie.CookieService;
import pl.pwr.zpi.utils.exception.AuthenticationException;

import java.io.IOException;
import java.security.GeneralSecurityException;
import java.time.Instant;
import java.util.Collections;

@Service
@Slf4j
public class GoogleOauthTokenService {

    private final String googleClientId;
    private final String googleClientSecret;
    private final CookieService cookieService;
    private final GsonFactory jsonFactory;
    private final NetHttpTransport transport;
    private final GoogleIdTokenVerifier verifier;

    public GoogleOauthTokenService(@Value("${spring.security.oauth2.client.registration.google.client-id}") String googleClientId,
                                   @Value("${spring.security.oauth2.client.registration.google.client-secret}") String googleClientSecret,
                                   CookieService cookieService) {
        this.cookieService = cookieService;
        this.googleClientId = googleClientId;
        this.googleClientSecret = googleClientSecret;
        this.jsonFactory = GsonFactory.getDefaultInstance();
        this.transport = new NetHttpTransport();
        this.verifier = new GoogleIdTokenVerifier.Builder(transport, jsonFactory)
                .setAudience(Collections.singletonList(googleClientId))
                .build();
    }

    public ResponseCookie updateAuthToken(String refreshToken) {
        var refreshedIdToken = refreshIdToken(refreshToken);
        return cookieService.createAuthCookie(refreshedIdToken.getIdToken(), Instant.ofEpochSecond(refreshedIdToken.getExpiresInSeconds()));
    }

    @SneakyThrows
    public GoogleIdToken.Payload decodeToken(String token) {
        GoogleIdToken idToken = verifier.verify(token);
        if (idToken != null) {
            return idToken.getPayload();
        }
        throw new AuthenticationException("Token validation failed");
    }

    public void validateToken(String token) throws GeneralSecurityException, IOException {
        verifier.verify(token);
    }

    private GoogleTokenResponse refreshIdToken(String refreshToken) {
        try {
            return new GoogleRefreshTokenRequest(
                    transport,
                    jsonFactory,
                    refreshToken,
                    googleClientId,
                    googleClientSecret)
                    .execute();

        } catch (IOException e) {
            throw new RuntimeException(e);
        }

    }
}
