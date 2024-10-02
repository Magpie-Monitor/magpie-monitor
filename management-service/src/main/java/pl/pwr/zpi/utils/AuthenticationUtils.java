package pl.pwr.zpi.utils;


import lombok.AllArgsConstructor;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.userdetails.UserDetailsService;
import org.springframework.security.core.userdetails.UsernameNotFoundException;
import org.springframework.security.oauth2.core.user.OAuth2User;
import org.springframework.stereotype.Component;
import pl.pwr.zpi.user.data.User;
import pl.pwr.zpi.user.service.UserService;

import java.util.Optional;


@Component
@AllArgsConstructor
public class AuthenticationUtils {

    private final UserService userService;
    public Optional<User> getOAuthUserFromAuthentication(Authentication authentication) {
        String email = ((OAuth2User)authentication.getPrincipal()).getAttribute("email");
        return userService.findByEmail(email);
    }
}