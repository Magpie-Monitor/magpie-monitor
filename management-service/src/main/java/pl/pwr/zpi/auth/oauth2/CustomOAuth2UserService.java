package pl.pwr.zpi.auth.oauth2;

import lombok.RequiredArgsConstructor;
import org.springframework.security.oauth2.client.userinfo.DefaultOAuth2UserService;
import org.springframework.security.oauth2.client.userinfo.OAuth2UserRequest;
import org.springframework.security.oauth2.core.OAuth2AuthenticationException;
import org.springframework.security.oauth2.core.oidc.user.DefaultOidcUser;
import org.springframework.security.oauth2.core.user.OAuth2User;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.user.data.User;
import pl.pwr.zpi.user.repository.UserRepository;

@Service
@RequiredArgsConstructor
public class CustomOAuth2UserService extends DefaultOAuth2UserService {

    private final UserRepository userRepository;

    @Override
    public OAuth2User loadUser(OAuth2UserRequest userRequest) throws OAuth2AuthenticationException {
        DefaultOidcUser defaultOidcUser = (DefaultOidcUser) super.loadUser(userRequest);

        String email = defaultOidcUser.getEmail();
        User user = userRepository.findByEmail(email);

        return new CustomOAuth2User(defaultOidcUser, user);
    }

}