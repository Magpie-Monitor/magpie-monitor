package pl.pwr.zpi.user.service;

import lombok.RequiredArgsConstructor;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.user.data.User;
import pl.pwr.zpi.user.repository.UserRepository;

import java.time.Instant;

@RequiredArgsConstructor
@Service
public class UserService {

    private final UserRepository userRepository;

    public User getCurrentUser() {
        var userDetails = (UserDetails) SecurityContextHolder.getContext().getAuthentication().getPrincipal();
        var userEmail = userDetails.getUsername();
        return getUserByEmail(userEmail);
    }

    public User getUserByEmail(String email) {
        return userRepository.findByEmail(email);
    }

    public User saveUser(User user) {
        return userRepository.save(user);
    }

    public User findById(Long id) {
        return userRepository.findById(id).orElseThrow();
    }

    public User findByEmail(String email) {
        return userRepository.findByEmail(email);
    }

    public void updateUserToken(String userEmail) {
        var user = userRepository.findByEmail(userEmail);
        user.setAuthTokenExpDate(Instant.now());
        userRepository.save(user);
    }
}
