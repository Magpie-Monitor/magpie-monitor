package pl.pwr.zpi.google.config;

import com.google.api.client.googleapis.auth.oauth2.GoogleAuthorizationCodeFlow;
import lombok.Getter;
import lombok.Setter;
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.context.annotation.Configuration;

import java.util.List;

@Configuration
@Getter
@Setter
@ConfigurationProperties(prefix = "spring.security.oauth2.client.registration.google")
public class GoogleApiProperties {

    private String clientId;
    private String clientSecret;
    private String scope;

    public String buildAuthorizationUri(String redirectUri,
                                        String state,
                                        String accessType,
                                        String email,
                                        List<String> requiredScopes,
                                        GoogleAuthorizationCodeFlow googleAuthorizationCodeFlow) {

        return googleAuthorizationCodeFlow
                .newAuthorizationUrl()
                .setClientId(clientId)
                .setRedirectUri(redirectUri)
                .setScopes(requiredScopes)
                .setState(state)
                .setAccessType(accessType)
                .set("login_hint", email)
                .set("prompt", "consent")
                .build();
    }

}