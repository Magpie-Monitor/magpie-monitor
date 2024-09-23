package pl.pwr.zpi.user.service;

import com.google.api.client.googleapis.auth.oauth2.GoogleRefreshTokenRequest;
import com.google.api.client.googleapis.auth.oauth2.GoogleTokenResponse;
import com.google.api.client.googleapis.javanet.GoogleNetHttpTransport;
import com.google.api.client.http.GenericUrl;
import com.google.api.client.http.HttpRequest;
import com.google.api.client.http.HttpRequestFactory;
import com.google.api.client.http.HttpResponseException;
import com.google.api.client.http.javanet.NetHttpTransport;
import com.google.api.client.json.gson.GsonFactory;
import lombok.AllArgsConstructor;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.autoconfigure.condition.ConditionalOnExpression;
import org.springframework.security.oauth2.core.OAuth2AccessToken;
import org.springframework.security.oauth2.core.OAuth2RefreshToken;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.user.data.User;
import pl.pwr.zpi.user.repository.UserRepository;

import java.io.IOException;
import java.security.GeneralSecurityException;
import java.time.Instant;

@Service
@AllArgsConstructor
public class OAuthUserServiceImpl implements OAuthUserService {

    private final UserRepository userRepository;

//    @Value("${spring.security.oauth2.client.registration.google.client-id}")
//    private String clientId;
//
//    @Value("${spring.security.oauth2.client.registration.google.client-secret}")
//    private String clientSecret;

    private static final String APPLICATION_NAME = "Your-Application-Name";
    private static final GsonFactory JSON_FACTORY = GsonFactory.getDefaultInstance();
    private static NetHttpTransport HTTP_TRANSPORT;
    static {
        try {
            HTTP_TRANSPORT = GoogleNetHttpTransport.newTrustedTransport();
        } catch (GeneralSecurityException | IOException e) {
            e.printStackTrace();
            HTTP_TRANSPORT = null;
        }
    }


    @Override
    public User findById(Long id) {
        return userRepository.findById(id).orElseThrow();
    }

    @Override
    public User findByEmail(String email) {
        return userRepository.findByEmail(email);
    }

    @Override
    @ConditionalOnExpression("!T(site.easy.to.build.crm.util.StringUtils).isEmpty('${spring.security.oauth2.client.registration.google.client-id:}')")
    public String refreshAccessTokenIfNeeded(User oauthUser) {
//        Instant now = Instant.now();
//        if (now.isBefore(oauthUser.getAccessTokenExpiration())) {
//            return oauthUser.getAccessToken();
//        }
//
//        GsonFactory jsonFactory = GsonFactory.getDefaultInstance();
//
//        // Create a new GoogleTokenResponse
//        GoogleTokenResponse tokenResponse;
//        try {
//            tokenResponse = new GoogleRefreshTokenRequest(
//                    new NetHttpTransport(),
//                    jsonFactory,
//                    oauthUser.getRefreshToken(),
//                    clientId,
//                    clientSecret)
//                    .execute();
//            String newAccessToken = tokenResponse.getAccessToken();
//            long expiresIn = tokenResponse.getExpiresInSeconds();
//            Instant expiresAt = Instant.now().plusSeconds(expiresIn);
//
//            oauthUser.setAccessToken(newAccessToken);
//            oauthUser.setAccessTokenExpiration(expiresAt);
//
//            oAuthUserRepository.save(oauthUser);
//        } catch (IOException e) {
//            throw new RuntimeException(e);
//        }
//
//        return oauthUser.getAccessToken();
        return null;
    }

    @Override
    @ConditionalOnExpression("!T(site.easy.to.build.crm.util.StringUtils).isEmpty('${spring.security.oauth2.client.registration.google.client-id:}')")
    public void revokeAccess(User oAuthUser) {
//        try {
//            final NetHttpTransport httpTransport = GoogleNetHttpTransport.newTrustedTransport();
//            HttpRequestFactory requestFactory = httpTransport.createRequestFactory();
//
//            GenericUrl url = new GenericUrl("https://accounts.google.com/o/oauth2/revoke");
//            url.set("token", oAuthUser.getAccessToken());
//
//            HttpRequest request = requestFactory.buildGetRequest(url);
//            request.execute();
//        } catch (HttpResponseException e) {
//            // Handle the error response if needed
//        } catch (GeneralSecurityException | IOException e) {
//            throw new RuntimeException(e);
//        }
    }

    @Override
    public void save(User oAuthUser) {
        userRepository.save(oAuthUser);
    }

    @Override
    public void deleteById(int id) {

    }

    public void updateOAuthUserTokens(User oAuthUser, OAuth2AccessToken oAuth2AccessToken, OAuth2RefreshToken oAuth2RefreshToken) {
//        oAuthUser.setAccessToken(oAuth2AccessToken.getTokenValue());
//        oAuthUser.setAccessTokenIssuedAt(oAuth2AccessToken.getIssuedAt());
//        oAuthUser.setAccessTokenExpiration(oAuth2AccessToken.getExpiresAt());
//
//        if(oAuth2RefreshToken != null) {
//            oAuthUser.setRefreshToken(oAuth2RefreshToken.getTokenValue());
//            oAuthUser.setRefreshTokenIssuedAt(oAuth2RefreshToken.getIssuedAt());
//            oAuthUser.setRefreshTokenExpiration(oAuth2RefreshToken.getExpiresAt());
//        }
    }

}
