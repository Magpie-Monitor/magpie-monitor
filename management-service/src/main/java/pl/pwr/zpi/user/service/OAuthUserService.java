package pl.pwr.zpi.user.service;

import org.springframework.security.oauth2.core.OAuth2AccessToken;
import org.springframework.security.oauth2.core.OAuth2RefreshToken;
import pl.pwr.zpi.user.data.User;


public interface OAuthUserService {

    public User findById(Long id);

    public User findByEmail(String email);

    public String refreshAccessTokenIfNeeded(User oauthUser);

    public void revokeAccess(User oAuthUser);

    public void save(User user);

    public void deleteById(int id);

    public void updateOAuthUserTokens(User oAuthUser, OAuth2AccessToken oAuth2AccessToken, OAuth2RefreshToken oAuth2RefreshToken);


}
