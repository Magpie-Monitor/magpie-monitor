package pl.pwr.zpi.auth.controller;

import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpHeaders;
import org.springframework.http.ResponseCookie;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.Authentication;
import org.springframework.web.bind.annotation.*;
import pl.pwr.zpi.auth.oauth2.OauthRefreshTokenService;
import pl.pwr.zpi.auth.service.AuthenticationService;

import java.security.Principal;

@RestController
@RequiredArgsConstructor
public class AuthenticationController {

    private final AuthenticationService service;
    private final OauthRefreshTokenService oauthRefreshTokenService;

    @GetMapping("/api/v1/auth/user-details")
    public ResponseEntity<?> getUser(Authentication authentication) {
        return ResponseEntity.ok().body(service.getUserDetails(authentication));
    }

    @GetMapping("/api/vi/auth/auth-token/validation-time")
    public ResponseEntity<?> getTokenValidationTime(Authentication authentication) {
        return ResponseEntity.ok().body(service.getTokenValidationTime(authentication));
    }

    @PostMapping("/api/v1/auth/refreshToken")
    public ResponseEntity<?> refreshToken(Authentication authentication) {
        ResponseCookie updatedToken = oauthRefreshTokenService.updateRefreshToken(authentication);

        return ResponseEntity
                .ok()
                .header(HttpHeaders.SET_COOKIE, updatedToken.toString())
                .build();
    }

}