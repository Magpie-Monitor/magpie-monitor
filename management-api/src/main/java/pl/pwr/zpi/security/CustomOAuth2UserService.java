package pl.pwr.zpi.security;

import lombok.RequiredArgsConstructor;
import org.springframework.security.oauth2.client.userinfo.DefaultOAuth2UserService;
import org.springframework.security.oauth2.client.userinfo.OAuth2UserRequest;
import org.springframework.security.oauth2.core.user.OAuth2User;
import org.springframework.stereotype.Component;
import pl.pwr.zpi.security.jwt.JwtService;
import pl.pwr.zpi.security.jwt.JwtToken;
import pl.pwr.zpi.user.data.User;
import pl.pwr.zpi.user.dto.Provider;
import pl.pwr.zpi.user.repository.UserRepository;


@Component
@RequiredArgsConstructor
public class CustomOAuth2UserService extends DefaultOAuth2UserService {

    private final GoogleOAuth2UserInfoExtractor oAuth2UserInfoExtractors;
    private final UserRepository userRepository;
    private final JwtService jwtService;


    @Override
    public OAuth2User loadUser(OAuth2UserRequest userRequest) {
        OAuth2User oAuth2User = super.loadUser(userRequest);
        User customUserDetails = oAuth2UserInfoExtractors.extractUserInfo(oAuth2User);
        processOAuthPostLogin(customUserDetails.getEmail(), customUserDetails.getNickname());
        System.out.println("CustomUserDetails: " + customUserDetails);
        return oAuth2User;
    }

    public JwtToken processOAuthPostLogin(String username, String nickname) {
        if (userRepository.existsByEmail(username))
            throw new RuntimeException("User with email " + username + " already exists");
        User newUser = User.builder()
                .email(username)
                .provider(Provider.GOOGLE)
                .nickname(nickname)
                .build();
        userRepository.save(newUser);

        return jwtService.generateToken(newUser);
    }

}