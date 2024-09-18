package pl.pwr.zpi.auth.dto;

import pl.pwr.zpi.security.jwt.JwtToken;

public record AuthenticationResponse(JwtToken token) {
}