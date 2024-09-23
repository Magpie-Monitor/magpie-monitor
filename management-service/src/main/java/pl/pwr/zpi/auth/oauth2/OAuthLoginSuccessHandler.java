package pl.pwr.zpi.auth.oauth2;

import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import jakarta.servlet.http.HttpSession;
import lombok.AllArgsConstructor;
import org.springframework.core.env.Environment;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.authority.SimpleGrantedAuthority;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.security.oauth2.client.OAuth2AuthorizedClient;
import org.springframework.security.oauth2.client.OAuth2AuthorizedClientService;
import org.springframework.security.oauth2.client.authentication.OAuth2AuthenticationToken;
import org.springframework.security.oauth2.core.OAuth2AccessToken;
import org.springframework.security.oauth2.core.OAuth2RefreshToken;
import org.springframework.security.oauth2.core.oidc.user.DefaultOidcUser;
import org.springframework.security.oauth2.core.user.OAuth2User;
import org.springframework.security.web.authentication.SimpleUrlAuthenticationSuccessHandler;
import org.springframework.stereotype.Component;
import pl.pwr.zpi.security.cookie.CookieService;
import pl.pwr.zpi.user.data.User;
import pl.pwr.zpi.user.service.OAuthUserService;
import pl.pwr.zpi.user.service.UserService;
import pl.pwr.zpi.utils.AuthenticationUtils;
import pl.pwr.zpi.utils.StringUtils;

import java.io.IOException;
import java.time.LocalDateTime;
import java.util.ArrayList;
import java.util.List;
import java.util.stream.Collectors;

@Component
@AllArgsConstructor
public class OAuthLoginSuccessHandler extends SimpleUrlAuthenticationSuccessHandler {

    public final OAuthUserService oAuthUserService;
    public final UserService userService;
    private final OAuth2AuthorizedClientService authorizedClientService;
    public final AuthenticationUtils authenticationUtils;
    private final Environment environment;
    private final CookieService cookieService;


    @Override
    public void onAuthenticationSuccess(HttpServletRequest request, HttpServletResponse response, Authentication authentication) throws IOException, ServletException {
        // Get the registration ID of the OAuth2 provider
        String googleClientId = environment.getProperty("spring.security.oauth2.client.registration.google.client-id");
        String googleClientSecret = environment.getProperty("spring.security.oauth2.client.registration.google.client-secret");
        boolean x = true;
        if (StringUtils.isEmpty(googleClientId) || StringUtils.isEmpty(googleClientSecret)) {
            response.sendRedirect("/error-page");
            return;
        }
        OAuth2AuthenticationToken oauthToken = (OAuth2AuthenticationToken) authentication;
        String registrationId = oauthToken.getAuthorizedClientRegistrationId();

        if (registrationId == null) {
            // Handle the case when the registrationId is not found
            throw new ServletException("Failed to find the registrationId from the authorities");
        }
        // Obtain the OAuth2AuthorizedClient
        OAuth2AuthorizedClient authorizedClient = authorizedClientService.loadAuthorizedClient(registrationId, authentication.getName());


        // Get the access and the refresh token from the OAuth2AuthorizedClient

        OAuth2AccessToken oAuth2AccessToken = authorizedClient.getAccessToken();
        OAuth2RefreshToken oAuth2RefreshToken = authorizedClient.getRefreshToken();
        System.out.println(oAuth2AccessToken.getTokenValue());
        System.out.println(oAuth2RefreshToken.getTokenValue());

        cookieService.createAuthCookie(oAuth2AccessToken.getTokenValue());
        cookieService.createRefreshCookie(oAuth2RefreshToken.getTokenValue());

        HttpSession session = request.getSession();
        boolean previouslyUsedRegularAccount = session.getAttribute("loggedInUserId") != null;
        Long userId = (previouslyUsedRegularAccount) ? (Long) session.getAttribute("loggedInUserId") : -1;
        User loggedUser = null;
        if (userId != -1) {
            loggedUser = userService.findById(userId);
        }
        User oAuthUser = authenticationUtils.getOAuthUserFromAuthentication(authentication);
        if (loggedUser != null && oAuthUser == null) {
            oAuthUser = new User();
            String email = ((DefaultOidcUser) authentication.getPrincipal()).getEmail();
            oAuthUser.setEmail(email);
            oAuthUserService.updateOAuthUserTokens(oAuthUser, oAuth2AccessToken, oAuth2RefreshToken);
            oAuthUserService.save(oAuthUser);
            response.sendRedirect("/connect-accounts");
        } else {

            String email = ((DefaultOidcUser) authentication.getPrincipal()).getEmail();
            String img = ((DefaultOidcUser) authentication.getPrincipal()).getPicture();
            String firstName = ((DefaultOidcUser) authentication.getPrincipal()).getGivenName();
            String lastName = ((DefaultOidcUser) authentication.getPrincipal()).getFamilyName();
            String username = email.split("@")[0];


//            Long currUserId = authenticationUtils.getLoggedInUserId(authentication);
//            User user = userService.findById(currUserId);
//            OAuthUser loggedOAuthUser;
//
//            if (user == null) {
//                user = new User();
//                UserProfile userProfile = new UserProfile();
//                userProfile.setFirstName(firstName);
//                userProfile.setLastName(lastName);
//                userProfile.setOathUserImageLink(img);
//                user.setEmail(email);
//                user.setUsername(username);
//                user.setPasswordSet(true);
//
//                long countUsers = userService.countAllUsers();
//                if (countUsers == 0) {
//                    role = roleService.findByName("ROLE_MANAGER");
//                    user.setStatus("active");
//                    userProfile.setStatus("active");
//                } else {
//                    role = roleService.findByName("ROLE_EMPLOYEE");
//                    user.setStatus("inactive");
//                    userProfile.setStatus("inactive");
//                }
//
//                user.setRoles(List.of(role));
//                user.setCreatedAt(LocalDateTime.now());
//                User createdUser = userService.save(user);
//                userProfile.setUser(createdUser);
//                userProfileService.save(userProfile);
//
//                loggedOAuthUser = new OAuthUser();
//                loggedOAuthUser.setEmail(email);
//                loggedOAuthUser.getGrantedScopes().addAll(List.of("openid", "email", "profile"));
//                oAuthUserService.updateOAuthUserTokens(loggedOAuthUser, oAuth2AccessToken, oAuth2RefreshToken);
//            } else {
//                loggedOAuthUser = user.getOauthUser();
//            }


            List<GrantedAuthority> updatedAuthorities = new ArrayList<>(authentication.getAuthorities());

            OAuth2User oauthUser = (OAuth2User) authentication.getPrincipal();

            Authentication updatedAuthentication = new OAuth2AuthenticationToken(
                    oauthUser,
                    updatedAuthorities,
                    registrationId
            );


            SecurityContextHolder.getContext().setAuthentication(updatedAuthentication);
        }
    }
}