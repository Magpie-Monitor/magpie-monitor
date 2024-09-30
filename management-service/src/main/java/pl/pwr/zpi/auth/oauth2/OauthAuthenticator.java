package pl.pwr.zpi.auth.oauth2;

import com.google.api.client.auth.oauth2.TokenResponseException;
import com.google.api.client.http.*;
import com.google.api.client.http.javanet.NetHttpTransport;
import com.google.api.client.json.JsonObjectParser;
import com.google.api.client.json.gson.GsonFactory;
import jakarta.servlet.*;
import lombok.extern.slf4j.Slf4j;
import org.springframework.security.core.Authentication;
import org.springframework.security.oauth2.client.OAuth2AuthorizedClient;
import org.springframework.security.oauth2.client.OAuth2AuthorizedClientService;
import org.springframework.security.oauth2.client.authentication.OAuth2AuthenticationToken;
import org.springframework.security.oauth2.core.OAuth2AccessToken;
import org.springframework.stereotype.Component;

import java.io.IOException;

@Component
@Slf4j
public class OauthAuthenticator implements Filter {
    private final OAuth2AuthorizedClientService authorizedClientService;
    private final HttpTransport httpTransport;
    private final GsonFactory jsonFactory;

    public OauthAuthenticator(OAuth2AuthorizedClientService authorizedClientService) {
        this.authorizedClientService = authorizedClientService;
        this.httpTransport = new NetHttpTransport();
        this.jsonFactory = new GsonFactory();
    }

    @Override
    public void doFilter(ServletRequest servletRequest, ServletResponse servletResponse, FilterChain filterChain) throws IOException {
        try {
            Authentication authentication = (Authentication) servletRequest.getAttribute("authentication");

            if (authentication instanceof OAuth2AuthenticationToken) {
                verifyUser(authentication);
            }

            filterChain.doFilter(servletRequest, servletResponse);

        } catch (ServletException | IOException e) {
            servletResponse.setContentType("application/json");
            servletResponse.getWriter().write("{\"error\":\"Authentication failed\"}");
            servletResponse.getWriter().flush();
        }
    }

    private void verifyUser(Authentication authentication) throws ServletException, IOException {
        try {
            OAuth2AuthenticationToken oauthToken = (OAuth2AuthenticationToken) authentication;
            String registrationId = oauthToken.getAuthorizedClientRegistrationId();

            if (registrationId == null) {
                throw new ServletException("Failed to find the registrationId from the authorities");
            }

            OAuth2AuthorizedClient authorizedClient = authorizedClientService.loadAuthorizedClient(registrationId, authentication.getName());

            OAuth2AccessToken oAuth2AccessToken = authorizedClient.getAccessToken();

            HttpRequestFactory requestFactory = httpTransport.createRequestFactory(request -> request.setParser(new JsonObjectParser(jsonFactory)));
            GenericUrl url = new GenericUrl("https://www.googleapis.com/oauth2/v3/tokeninfo");
            url.set("access_token", oAuth2AccessToken.getTokenValue());

            HttpRequest request = requestFactory.buildGetRequest(url);
            HttpResponse response = request.execute();

            if (response.getStatusCode() != 200) {
                throw new ServletException("Failed to verify OAuth token, status: " + response.getStatusCode());
            }

            log.info("User verified successfully");

        } catch (TokenResponseException e) {
            if (e.getDetails() != null && "invalid_grant".equals(e.getDetails().getError())) {
                throw new ServletException("Token has been expired or revoked", e);
            } else {
                throw new ServletException("Token response error", e);
            }
        } catch (IOException e) {
            throw new IOException("IO error while verifying OAuth token", e);
        }
    }
}
