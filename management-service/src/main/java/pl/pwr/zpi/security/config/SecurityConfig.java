package pl.pwr.zpi.security.config;

import lombok.RequiredArgsConstructor;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.http.HttpStatus;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity;
import org.springframework.security.config.annotation.web.configurers.AbstractHttpConfigurer;
import org.springframework.security.oauth2.client.OAuth2AuthorizedClientService;
import org.springframework.security.web.SecurityFilterChain;
import org.springframework.security.web.access.AccessDeniedHandler;
import org.springframework.security.web.access.intercept.AuthorizationFilter;
import org.springframework.security.web.authentication.logout.HttpStatusReturningLogoutSuccessHandler;
import org.springframework.web.cors.CorsConfiguration;
import org.springframework.web.cors.CorsConfigurationSource;
import org.springframework.web.cors.UrlBasedCorsConfigurationSource;
import pl.pwr.zpi.auth.CustomAccessDeniedHandler;
import pl.pwr.zpi.auth.oauth2.CustomOAuth2UserService;
import pl.pwr.zpi.auth.oauth2.OAuthLoginSuccessHandler;
import pl.pwr.zpi.auth.oauth2.OauthAuthenticator;

import java.util.Arrays;

@Configuration
@EnableWebSecurity
@RequiredArgsConstructor
public class SecurityConfig {

    private final OAuthLoginSuccessHandler oAuth2LoginSuccessHandler;
    private final CustomOAuth2UserService oauthUserService;
    private final OAuth2AuthorizedClientService authorizedClientService;

    @Bean
    public SecurityFilterChain securityFilterChain(HttpSecurity http, CorsConfigurationSource corsConfigurationSource) throws Exception {
        http
                .csrf(AbstractHttpConfigurer::disable)
                .cors(cors -> cors.configurationSource(corsConfigurationSource))
                .authorizeHttpRequests(request -> {
                    request.requestMatchers(
                            "/public/api/**",
                            "/v3/api-docs/**",
                            "/swagger-ui/**",
                            "/api/v1/auth/login",
                            "/login/oauth2/code/**",
                            "/",
                            "/login**",
                            "/error",
                            "/api/**"
                            ).permitAll();
                    request.anyRequest().authenticated();
                })
                .oauth2Login(oauth2 -> oauth2
                        .userInfoEndpoint(userInfo -> userInfo
                                .userService(oauthUserService))
                        .successHandler(oAuth2LoginSuccessHandler))
                .logout((logout) -> logout
                        .logoutUrl("/api/v1/auth/logout")
                        .clearAuthentication(true)
                        .invalidateHttpSession(true)
                        .deleteCookies("authToken")
                        .deleteCookies("refreshToken")
                        .deleteCookies("JSESSIONID")
                        .logoutSuccessHandler(new HttpStatusReturningLogoutSuccessHandler(HttpStatus.OK)))
                .addFilterBefore(new OauthAuthenticator(authorizedClientService), AuthorizationFilter.class)
                .exceptionHandling(exception ->
                    exception.accessDeniedHandler(accessDeniedHandler())
                );
        return http.build();
    }

    @Bean
    public CorsConfigurationSource corsConfigurationSource(@Value("${frontend.client.url}") String frontendUrl) {
        CorsConfiguration configuration = new CorsConfiguration();
        configuration.addAllowedOrigin(frontendUrl);
        configuration.addAllowedOrigin("http://localhost"); //TODO: Remove before release
        configuration.setAllowedMethods(Arrays.asList("GET", "POST", "PATCH", "DELETE", "HEAD", "OPTIONS"));
        configuration.addAllowedHeader("*");
        configuration.setAllowCredentials(true);
        configuration.setMaxAge(3600L);
        UrlBasedCorsConfigurationSource source = new UrlBasedCorsConfigurationSource();
        source.registerCorsConfiguration("/**", configuration);
        return source;
    }

    @Bean
    public AccessDeniedHandler accessDeniedHandler() {
        return new CustomAccessDeniedHandler();
    }
}
