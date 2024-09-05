package pl.pwr.zpi.auth.service;

import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.user.dto.UserDTO;
import pl.pwr.zpi.user.repository.UserRepository;

@RequiredArgsConstructor
@Service
public class AuthenticationService {

    private final UserRepository userRepository;

    public UserDTO getUserDetails(String username) {
        var user = userRepository.findByEmail(username).orElseThrow();
        return UserDTO.toUserDTO(user);
    }
}
