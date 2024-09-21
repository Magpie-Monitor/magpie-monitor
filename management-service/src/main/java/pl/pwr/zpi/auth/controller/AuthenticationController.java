package pl.pwr.zpi.auth.controller;

import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import pl.pwr.zpi.auth.service.AuthenticationService;

import java.security.Principal;

@RestController
@RequiredArgsConstructor
public class AuthenticationController {

    private final AuthenticationService service;

    @GetMapping("/api/v1/auth/user-details")
    public ResponseEntity<?> getUser(Principal principal) {
        return ResponseEntity.ok().body(service.getUserDetails(principal.getName()));
    }

    @GetMapping("/api/vi/auth/token-validation-time")
    public ResponseEntity<?> getTokenValidationTime(@CookieValue("authToken") String authToken) {
        return ResponseEntity.ok().body(service.getTokenValidationTime(authToken));
    }
}