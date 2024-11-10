package pl.pwr.zpi.auth.service;

import com.google.api.client.http.*;
import com.google.api.client.http.javanet.NetHttpTransport;
import com.google.api.client.json.JsonObjectParser;
import com.google.api.client.json.gson.GsonFactory;
import com.google.api.client.util.GenericData;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.auth.dto.TokenExpTime;
import pl.pwr.zpi.user.dto.UserDTO;

import java.io.IOException;
import java.time.Instant;

@RequiredArgsConstructor
@Service
public class AuthenticationService {

    private static final String TOKEN_INFO_URL = "https://www.googleapis.com/oauth2/v3/tokeninfo";
    private final HttpTransport httpTransport = new NetHttpTransport();
    private final GsonFactory jsonFactory = new GsonFactory();

    public UserDTO getUserDetails(String authToken) throws IOException {
        String accessToken = extractAccessToken(authToken);
        GenericData tokenInfo = fetchTokenInfo(accessToken);

        String email = (String) tokenInfo.get("email");
        String nickname = extractNicknameFromEmail(email);

        return UserDTO.builder()
                .nickname(nickname)
                .email(email)
                .build();
    }

    public TokenExpTime getTokenValidationTime(String authToken) {
        long expiresAtMillis = extractExpiryTime(authToken);
        long timeRemainingMillis = expiresAtMillis - Instant.now().toEpochMilli();
        return new TokenExpTime(timeRemainingMillis);
    }

    private GenericData fetchTokenInfo(String accessToken) throws IOException {
        HttpRequestFactory requestFactory = httpTransport.createRequestFactory(request ->
                request.setParser(new JsonObjectParser(jsonFactory))
        );

        GenericUrl url = new GenericUrl(TOKEN_INFO_URL);
        url.set("access_token", accessToken);

        HttpRequest request = requestFactory.buildGetRequest(url);
        HttpResponse response = request.execute();

        return response.parseAs(GenericData.class);
    }

    private String extractAccessToken(String authToken) {
        String[] parts = authToken.split("\\|");
        if (parts.length != 2) {
            throw new IllegalArgumentException("Invalid authToken format");
        }
        return parts[0];
    }

    private long extractExpiryTime(String authToken) {
        String[] parts = authToken.split("\\|");
        if (parts.length != 2) {
            throw new IllegalArgumentException("Invalid authToken format");
        }
        try {
            return Long.parseLong(parts[1]);
        } catch (NumberFormatException e) {
            throw new IllegalArgumentException("Invalid expiry time format in authToken", e);
        }
    }

    private String extractNicknameFromEmail(String email) {
        return (email != null && email.contains("@")) ? email.split("@")[0] : "unknown";
    }
}
