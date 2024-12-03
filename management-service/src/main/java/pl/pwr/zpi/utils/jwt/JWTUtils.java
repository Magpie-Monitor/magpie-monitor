package pl.pwr.zpi.utils.jwt;

import com.google.api.client.googleapis.auth.oauth2.GoogleIdToken;
import com.google.api.client.googleapis.auth.oauth2.GoogleIdToken.Payload;
import com.google.api.client.googleapis.auth.oauth2.GoogleIdTokenVerifier;
import com.google.api.client.http.javanet.NetHttpTransport;
import com.google.api.client.json.gson.GsonFactory;
import lombok.RequiredArgsConstructor;
import lombok.SneakyThrows;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;
import pl.pwr.zpi.utils.exception.AuthenticationException;

import java.util.Collections;

@Component
@RequiredArgsConstructor
public class JWTUtils {

    @Value("${spring.security.oauth2.client.registration.google.client-id}")
    private String googleClientId;
    private static final NetHttpTransport transport = new NetHttpTransport();
    private static final GsonFactory jsonFactory = GsonFactory.getDefaultInstance();

    @SneakyThrows
    public Payload decodeToken(String token) {
        GoogleIdTokenVerifier verifier = new GoogleIdTokenVerifier.Builder(transport, jsonFactory)
                .setAudience(Collections.singletonList(googleClientId))
                .build();


        GoogleIdToken idToken = verifier.verify(token);
        if (idToken != null) {
            return idToken.getPayload();
        }
        throw new AuthenticationException("Token validation failed");
    }

}