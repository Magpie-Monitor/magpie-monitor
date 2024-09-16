package pl.pwr.zpi.auth.controller;

import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import pl.pwr.zpi.auth.dto.ErrorMessageDTO;
import pl.pwr.zpi.auth.service.AuthenticationService;

import java.security.Principal;


@RestController
@RequiredArgsConstructor
public class AuthenticationController {

    private final AuthenticationService service;

    @GetMapping("/public/api/v1/user-details")
    public ResponseEntity<?> getUser(Principal principal) {
        try {
            return ResponseEntity.ok().body(service.getUserDetails(principal.getName()));
        } catch (Exception e) {
            return ResponseEntity.badRequest().body(new ErrorMessageDTO(e.getMessage()));
        }
    }

    @GetMapping("/api/vi/auth/token-validation-time")
    public ResponseEntity<?> getTokenValidationTime(@CookieValue("authToken") String authToken) {
        try {
            return ResponseEntity.ok().body(service.getTokenValidationTime(authToken));
        } catch (Exception e) {
            return ResponseEntity.badRequest().body(new ErrorMessageDTO(e.getMessage()));
        }
    }
}
