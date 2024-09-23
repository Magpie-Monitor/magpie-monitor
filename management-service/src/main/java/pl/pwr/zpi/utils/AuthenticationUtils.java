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


@Component
@AllArgsConstructor
public class AuthenticationUtils {

    private final UserService userService;
    public User getOAuthUserFromAuthentication(Authentication authentication) {
        String email = ((OAuth2User)authentication.getPrincipal()).getAttribute("email");
        return userService.findByEmail(email);
    }

//    public Long getLoggedInUserId(Authentication authentication) {
//        User user;
//        CustomerLoginInfo customerLoginInfo;
//        if (authentication instanceof UsernamePasswordAuthenticationToken) {
//            UserDetailsService authenticatedUserDetailsService = getAuthenticatedUserDetailsService(authentication);
//            String userName = authentication.getName();
//            if (authenticatedUserDetailsService == crmUserDetails) {
//                user = userService.findByUsername(userName).get(0);
//                if (user == null) {
//                    return -1;
//                }
//                return user.getId();
//            } else if (authenticatedUserDetailsService == customerUserDetails) {
//                customerLoginInfo = customerLoginInfoService.findByEmail(userName);
//                if (customerLoginInfo == null) {
//                    return -1;
//                }
//                return customerLoginInfo.getId();
//            }
//        } else {
//            OAuthUser oAuthUser = getOAuthUserFromAuthentication(authentication);
//            if (oAuthUser == null) {
//                return -1;
//            }
//            user = oAuthUser.getUser();
//            return user.getId();
//        }
//        return -1;
//    }
//    public boolean checkIfAppHasAccess(String serviceAccessUrl, OAuthUser oAuthUser) {
//        return oAuthUser.getGrantedScopes().contains(serviceAccessUrl);
//    }
//
//    public UserDetailsService getAuthenticatedUserDetailsService(Authentication authentication) {
//        UserDetailsService authenticatedUserDetailsService = null;
//
//        if (authentication != null && authentication.getPrincipal() instanceof org.springframework.security.core.userdetails.User authenticatedUser) {
//            String authenticatedUsername = authenticatedUser.getUsername();
//
//            if (authenticatedUsername != null) {
//                try {
//                    if (crmUserDetails != null && authenticatedUsername.equals(crmUserDetails.loadUserByUsername(authenticatedUsername).getUsername())) {
//                        authenticatedUserDetailsService = crmUserDetails;
//                    }
//                } catch (UsernameNotFoundException e) {
//                    // Swallow the exception and continue to the next condition
//                }
//
//                if (authenticatedUserDetailsService == null && customerUserDetails != null) {
//                    try {
//                        if (authenticatedUsername.equals(customerUserDetails.loadUserByUsername(authenticatedUsername).getUsername())) {
//                            authenticatedUserDetailsService = customerUserDetails;
//                        }
//                    } catch (UsernameNotFoundException e) {
//                        // Swallow the exception and continue to the next steps
//                    }
//                }
//            }
//        }
//
//        return authenticatedUserDetailsService;
//    }
}