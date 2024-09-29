package pl.pwr.zpi.user.service;

import org.springframework.security.oauth2.core.OAuth2AccessToken;
import org.springframework.security.oauth2.core.OAuth2RefreshToken;
import pl.pwr.zpi.user.data.User;


public interface OAuthUserService {

    User findById(Long id);

    User findByEmail(String email);

    String refreshAccessTokenIfNeeded(User oauthUser);

    void revokeAccess(User oAuthUser);

    void save(User user);

    void deleteById(int id);

    void updateOAuthUserTokens(User oAuthUser, OAuth2AccessToken oAuth2AccessToken, OAuth2RefreshToken oAuth2RefreshToken);
}
