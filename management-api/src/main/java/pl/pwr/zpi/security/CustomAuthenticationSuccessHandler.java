package pl.pwr.zpi.security;

import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import lombok.RequiredArgsConstructor;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.ResponseCookie;
import org.springframework.security.core.Authentication;
import org.springframework.security.oauth2.core.user.OAuth2User;
import org.springframework.security.web.authentication.SimpleUrlAuthenticationSuccessHandler;
import org.springframework.stereotype.Component;
import pl.pwr.zpi.security.cookie.CookieService;
import pl.pwr.zpi.security.jwt.JwtService;
import pl.pwr.zpi.security.jwt.JwtToken;
import pl.pwr.zpi.user.data.User;
import pl.pwr.zpi.user.dto.Provider;
import pl.pwr.zpi.user.repository.UserRepository;

import java.io.IOException;
import java.util.Optional;

@RequiredArgsConstructor
@Component
public class CustomAuthenticationSuccessHandler extends SimpleUrlAuthenticationSuccessHandler {

    private final JwtService tokenProvider;
    private final UserRepository userRepository;
    private final CookieService cookieService;

    @Value("${oauth2.google.redirect-uri}")
    private String REDIRECT_URI;

    @Override
    public void onAuthenticationSuccess(HttpServletRequest request, HttpServletResponse response, Authentication authentication) throws IOException {
        handle(request, response, authentication);
        super.clearAuthenticationAttributes(request);
    }

    @Override
    protected void handle(HttpServletRequest request, HttpServletResponse response, Authentication authentication) throws IOException {
        String targetUrl = REDIRECT_URI.isEmpty() ?
                determineTargetUrl(request, response, authentication) : REDIRECT_URI;

        Optional<String> optionalEmail = getEmail(authentication);
        Optional<String> optionalName = getUsername(authentication);

        if (optionalEmail.isEmpty()) {
            throw new IllegalStateException("Email not found in authentication");
        }

        String email = optionalEmail.get();
        String name = optionalName.orElse("Unknown User");
        User user = userRepository.findByEmail(email)
                .orElseGet(() -> persistNewUser(email, name));

        JwtToken token = tokenProvider.generateToken(user);
        ResponseCookie authCookie = cookieService.createAuthCookie(token.token());
        response.addHeader("Set-Cookie", authCookie.toString());

        getRedirectStrategy().sendRedirect(request, response, targetUrl);
    }


    private Optional<String> getEmail(Authentication authentication) {
        Object principal = authentication.getPrincipal();
        if (principal instanceof OAuth2User oauth2User) {
            return Optional.ofNullable(oauth2User.getAttribute("email"));
        }
        return Optional.empty();
    }

    private Optional<String> getUsername(Authentication authentication) {
        String details = authentication.getPrincipal().toString();
        String namePrefix = "name=";
        int startIndex = details.indexOf(namePrefix);
        if (startIndex != -1) {
            startIndex += namePrefix.length();
            int endIndex = details.indexOf(',', startIndex);
            if (endIndex == -1) {
                endIndex = details.indexOf('}', startIndex);
            }
            if (endIndex != -1) {
                return Optional.of(details.substring(startIndex, endIndex).trim());
            }
        }
        return Optional.empty();
    }

    private User persistNewUser(String email, String username) {
        Optional<User> existingUser = userRepository.findByEmail(email);
        if (existingUser.isPresent()) {
            return existingUser.get();
        }
        User newUser = User.builder()
                .email(email)
                .provider(Provider.GOOGLE)
                .nickname(username)
                .build();
        userRepository.save(newUser);
        return newUser;
    }
}