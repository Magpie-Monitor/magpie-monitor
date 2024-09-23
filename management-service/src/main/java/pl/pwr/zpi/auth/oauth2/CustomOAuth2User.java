package pl.pwr.zpi.auth.oauth2;

import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.authority.SimpleGrantedAuthority;
import org.springframework.security.oauth2.core.oidc.user.DefaultOidcUser;
import org.springframework.security.oauth2.core.user.OAuth2User;
import pl.pwr.zpi.user.data.User;

import java.util.Collection;
import java.util.Collections;
import java.util.Map;
import java.util.stream.Collectors;


public class CustomOAuth2User implements OAuth2User {
    private OAuth2User oauth2User;
    private final User user;
    private final DefaultOidcUser defaultOidcUser;
    public CustomOAuth2User(DefaultOidcUser defaultOidcUser, User user) {
        this.defaultOidcUser = defaultOidcUser;
        this.user = user;
    }

    @Override
    public <A> A getAttribute(String name) {
        return OAuth2User.super.getAttribute(name);
    }

    @Override
    public Map<String, Object> getAttributes() {
        return defaultOidcUser.getAttributes();
    }

    @Override
    public Collection<? extends GrantedAuthority> getAuthorities() {
        return Collections.emptyList();
    }

    @Override
    public String getName() {
        return defaultOidcUser.<String>getAttribute("name");
    }
    public String getEmail() {
        return defaultOidcUser.<String>getAttribute("email");
    }

}
