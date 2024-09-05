package pl.pwr.zpi.security;

import lombok.extern.slf4j.Slf4j;
import org.springframework.security.oauth2.client.userinfo.OAuth2UserRequest;
import org.springframework.security.oauth2.core.user.OAuth2User;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.user.data.User;
import pl.pwr.zpi.user.dto.Provider;

@Service
public class GoogleOAuth2UserInfoExtractor {

    public User extractUserInfo(OAuth2User oAuth2User) {
        User customUserDetails = new User();
        customUserDetails.setEmail(retrieveAttr("email", oAuth2User));
        customUserDetails.setNickname(retrieveAttr("name", oAuth2User));
        customUserDetails.setProvider(Provider.GOOGLE);
        return customUserDetails;
    }

    public boolean accepts(OAuth2UserRequest userRequest) {
        return Provider.GOOGLE.name().equalsIgnoreCase(userRequest.getClientRegistration().getRegistrationId());
    }

    private String retrieveAttr(String attr, OAuth2User oAuth2User) {
        Object attribute = oAuth2User.getAttributes().get(attr);
        return attribute == null ? "" : attribute.toString();
    }
}