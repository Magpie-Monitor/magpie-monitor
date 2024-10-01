package pl.pwr.zpi;

import lombok.RequiredArgsConstructor;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.scheduling.annotation.EnableAsync;
import org.springframework.security.core.userdetails.UserDetailsService;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.password.PasswordEncoder;
import pl.pwr.zpi.user.repository.UserRepository;
import pl.pwr.zpi.user.service.UserService;

@Configuration
@RequiredArgsConstructor
@EnableAsync
public class ApplicationConfig {

    private final UserService userService;

    @Bean
    public UserDetailsService userDetailsService() {
        return userService::findByEmailNoOptional;
    }

    @Bean
    public PasswordEncoder passwordEncoder() {
        return new BCryptPasswordEncoder();
    }
}
