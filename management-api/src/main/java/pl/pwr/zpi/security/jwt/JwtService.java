package pl.pwr.zpi.security.jwt;

import io.jsonwebtoken.Claims;
import io.jsonwebtoken.Jwts;
import io.jsonwebtoken.io.Decoders;
import io.jsonwebtoken.security.Keys;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.stereotype.Service;

import javax.crypto.SecretKey;
import java.util.Date;
import java.util.HashMap;
import java.util.Map;
import java.util.function.Function;

@Service
public class JwtService {
    @Value("${jwt.secret-key}")
    private String SECRET_KEY;

    @Value("${jwt.expiration-time-seconds}")
    private long EXPIRATION_TIME_SECONDS;

    public String extractUsername(JwtToken token) {
        return extractClaim(token, Claims::getSubject);
    }

    public <T> T extractClaim(JwtToken token, Function<Claims, T> claimsResolver) {
        final Claims claims = extractAllClaims(token);
        return claimsResolver.apply(claims);
    }

    public JwtToken generateToken(UserDetails userDetails) {
        return generateToken(new HashMap<>(), userDetails);
    }

    public JwtToken generateToken(Map<String, Object> extraClaims, UserDetails userDetails) {
        return new JwtToken(Jwts
                .builder()
                .claims(extraClaims)
                .subject(userDetails.getUsername())
                .issuedAt(new Date(System.currentTimeMillis()))
                .expiration(new Date(System.currentTimeMillis() + 1000 * EXPIRATION_TIME_SECONDS))
                .signWith(getPublicSigningKey())
                .compact());
    }

    public boolean isTokenValid(JwtToken token, UserDetails userDetails) {
        final String username = extractUsername(token);
        return (username.equals(userDetails.getUsername())) && !isTokenExpired(token);
    }

    public boolean isTokenExpired(JwtToken token) {
        return extractExpiration(token).before(new Date());
    }

    private Date extractExpiration(JwtToken token) {
        return extractClaim(token, Claims::getExpiration);
    }

    private Claims extractAllClaims(JwtToken token) {
        return Jwts.parser()
                .verifyWith(getPublicSigningKey())
                .build()
                .parseSignedClaims(token.token())
                .getPayload();
    }

    private SecretKey getPublicSigningKey() {
        byte[] keyBytes = Decoders.BASE64.decode(SECRET_KEY);
        return Keys.hmacShaKeyFor(keyBytes);
    }

}
