package pl.pwr.zpi.security;

import com.google.api.client.auth.oauth2.TokenResponseException;
import com.google.api.client.googleapis.auth.oauth2.GoogleAuthorizationCodeFlow;
import com.google.api.client.googleapis.auth.oauth2.GoogleRefreshTokenRequest;
import com.google.api.client.googleapis.auth.oauth2.GoogleTokenResponse;
import com.google.api.client.http.*;
import com.google.api.client.http.javanet.NetHttpTransport;
import com.google.api.client.json.JsonObjectParser;
import com.google.api.client.json.gson.GsonFactory;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpSession;
import lombok.AllArgsConstructor;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.security.core.Authentication;
import org.springframework.security.oauth2.core.oidc.user.DefaultOidcUser;
import org.springframework.stereotype.Service;
import org.springframework.web.servlet.view.RedirectView;
import pl.pwr.zpi.google.config.GoogleApiProperties;
import pl.pwr.zpi.google.config.GoogleAuthorizationCodeFlowWrapper;
import pl.pwr.zpi.user.data.User;
import pl.pwr.zpi.user.service.UserService;
import pl.pwr.zpi.utils.AuthenticationUtils;
import pl.pwr.zpi.utils.StringSetWrapper;

import java.io.IOException;
import java.security.GeneralSecurityException;
import java.time.Instant;
import java.util.*;

@Service
@AllArgsConstructor
@Deprecated
public class GoogleAccessServiceImpl implements GoogleAccessService {

//    @Value("${app.base-url}")
    private final String domain = "http://localhost:8080/";
    private final UserService userService;
    private final GoogleApiProperties googleApiProperties;
    private final GoogleAuthorizationCodeFlowWrapper googleAuthorizationCodeFlowWrapper;
    private final AuthenticationUtils authenticationUtils;

    @Override
    public RedirectView grantGoogleAccess(Authentication authentication, HttpSession session, boolean grantCalendarAccess, boolean grantGmailAccess, boolean grantDriveAccess, HttpServletRequest request) {
        String state = "YOUR_STATE_VALUE";
        String accessType = "offline";

        String email = ((DefaultOidcUser) authentication.getPrincipal()).getEmail();

        Set<String> newGrantedScopes = new HashSet<>();

        session.setAttribute("updatedScopes", new StringSetWrapper(newGrantedScopes));
        String[] scopes = googleApiProperties.getScope().split(",");
        Set<String> SetOfRequiredScopes = new HashSet<>(Arrays.asList(scopes));
        SetOfRequiredScopes.addAll(newGrantedScopes);

        String mainDomain = domain;
        while (mainDomain.endsWith("/") || mainDomain.endsWith(" ")) {
            mainDomain = mainDomain.substring(0, mainDomain.length() - 1);
        }

        String homeLink = request.getContextPath().isEmpty() ? "/" : request.getContextPath() + "/";
        mainDomain = mainDomain+homeLink+REDIRECT_URI;

        List<String> requiredScopes = new ArrayList<>(SetOfRequiredScopes);
        GoogleAuthorizationCodeFlow authorizationCodeFlow = googleAuthorizationCodeFlowWrapper.build(requiredScopes);
        String authorizationRequestUrl = googleApiProperties.buildAuthorizationUri(mainDomain, state, accessType, email, requiredScopes, authorizationCodeFlow);

        return new RedirectView(authorizationRequestUrl);
    }

    @Override
    public String handleGrantedAccess(HttpSession session, String error, String authCode, String state, Authentication authentication,  HttpServletRequest request) throws IOException {

        if ("access_denied".equals(error)) {
            // The user has canceled the consent page, so update the granted scopes
            return "redirect:/register";
        }

//        Long userId = authenticationUtils.getLoggedInUserId(authentication);
//        User user = userService.findById(userId);
//        User oAuthUser = authenticationUtils.getOAuthUserFromAuthentication(authentication);

//        handleScopeChanges(session, oAuthUser);
        List<String> requiredScopes = new ArrayList<>();

        GoogleAuthorizationCodeFlow flow = googleAuthorizationCodeFlowWrapper.build(requiredScopes);

        String mainDomain = domain;
        while (mainDomain.endsWith("/") || mainDomain.endsWith(" ")) {
            mainDomain = mainDomain.substring(0, mainDomain.length() - 1);
        }

        String homeLink = request.getContextPath().isEmpty() ? "/" : request.getContextPath() + "/";
        mainDomain = mainDomain+homeLink+REDIRECT_URI;

        // Exchange the authorization code for an access token and a refresh token.
        GoogleTokenResponse tokenResponse = flow.newTokenRequest(authCode).setRedirectUri(mainDomain).execute();

        // Obtain the OAuthUser object that represents the authenticated user. This assumes you have a method to retrieve the OAuthUser.

        // Update the OAuthUser object with the new access token, refresh token, and access token expiration.
//        oAuthUser.setAccessToken(tokenResponse.getAccessToken());
//        oAuthUser.setRefreshToken(tokenResponse.getRefreshToken());
//        oAuthUser.setAccessTokenExpiration(Instant.ofEpochMilli(tokenResponse.getExpiresInSeconds() * 1000L));

//        oAuthUserService.refreshAccessTokenIfNeeded(oAuthUser);

        Set<String> actualGrantedScopes = extractActualGrantedScopes(tokenResponse);
//        oAuthUser.setGrantedScopes(actualGrantedScopes);
//        oAuthUserService.save(oAuthUser, user);
//        if (actualGrantedScopes.contains(SCOPE_DRIVE)) {
//            try {
//                googleDriveApiService.findOrCreateTemplateFolder(oAuthUser, "Templates");
//            } catch (IOException | GeneralSecurityException e) {
////                throw new RuntimeException(e);
//            }
//        }


        return "redirect:/employee/settings/google-services";
    }

    public void verifyAccessAndHandleRevokedToken(User user, List<String> scopesToCheck) throws IOException {
        try {
            // Request a new access token using the refresh token
            HttpTransport httpTransport = new NetHttpTransport();
            GsonFactory jsonFactory = new GsonFactory();

            GoogleRefreshTokenRequest refreshTokenRequest = new GoogleRefreshTokenRequest(
                    httpTransport,
                    jsonFactory,
                    "",
                    googleApiProperties.getClientId(),
                    googleApiProperties.getClientSecret()
            );
            GoogleTokenResponse tokenResponse = refreshTokenRequest.execute();
            String newAccessToken = tokenResponse.getAccessToken();

            HttpRequestFactory requestFactory = httpTransport.createRequestFactory(request -> request.setParser(new JsonObjectParser(jsonFactory)));
            GenericUrl url = new GenericUrl("https://www.googleapis.com/oauth2/v3/tokeninfo");
            url.set("access_token", newAccessToken);

            HttpRequest request = requestFactory.buildGetRequest(url);
            HttpResponse response = request.execute();

        } catch (TokenResponseException e) {
            if (e.getDetails() != null && "invalid_grant".equals(e.getDetails().getError())) {
                // Handle the specific case when the token has been expired or revoked
                String[] scopes = googleApiProperties.getScope().split(",");
                Set<String> setOfRequiredScopes = new HashSet<>(Arrays.asList(scopes));
//                oAuthUserService.save(oAuthUser, user);
            }
//            else {
//                // Log and handle other exceptions as needed
//            }
            return;
        } catch (IOException e) {
            // Log and handle other IOExceptions as needed
            return;
        }
    }

    private void updateAccess(Set<String> grantedScopes, boolean grantAccess, String scope) {
        if (grantAccess) {
            grantedScopes.add(scope);
        } else {
            grantedScopes.remove(scope);
        }
    }

//    private void handleScopeChanges(HttpSession session) {
//        StringSetWrapper stringSetWrapper = SessionUtils.getSessionAttribute(session, "updatedScopes", StringSetWrapper.class);
//        Set<String> scopesToRemove = new HashSet<>();
//        Set<String> scopesToAdd = new HashSet<>();
//        boolean scopeChanges = false;
//        String[] scopes = googleApiProperties.getScope().split(",");
//        Set<String> grantedScopes = new HashSet<>(Arrays.asList(scopes));
//
//        if (stringSetWrapper != null) {
//            Set<String> updatedScopes = stringSetWrapper.getStringSet();
//            for (String scope : updatedScopes) {
//                if (!oAuthUser.getGrantedScopes().contains(scope)) {
//                    scopesToAdd.add(scope);
//                    scopeChanges = true;
//                }
//            }
//        }
//        for (String scope : oAuthUser.getGrantedScopes()) {
//            if (grantedScopes.contains(scope)) {
//                continue;
//            }
//            if (stringSetWrapper == null || !stringSetWrapper.getStringSet().contains(scope)) {
//                scopesToRemove.add(scope);
//                scopeChanges = true;
//            }
//        }
//
//        oAuthUser.getGrantedScopes().addAll(scopesToAdd);
//        oAuthUser.getGrantedScopes().removeAll(scopesToRemove);
//
//        if (scopeChanges) {
//            oAuthUserService.revokeAccess(oAuthUser);
//        }
//    }

    private Set<String> extractActualGrantedScopes(GoogleTokenResponse tokenResponse) {
        String grantedScopesStr = tokenResponse.get("scope").toString();
        Set<String> actualGrantedScopes = new HashSet<>(Arrays.asList(grantedScopesStr.split(" ")));
        Set<String> userGrantedScopes = new HashSet<>();
        for (String scope : actualGrantedScopes) {
            if (scope.contains("userinfo.")) {
                if (scope.contains("email")) {
                    userGrantedScopes.add("email");
                }
                if (scope.contains("profile")) {
                    userGrantedScopes.add("profile");
                }
            } else {
                userGrantedScopes.add(scope);
            }
        }
        return userGrantedScopes;
    }
}
