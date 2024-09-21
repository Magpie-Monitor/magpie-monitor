package pl.pwr.zpi.healthcheck;

import io.swagger.v3.oas.annotations.Operation;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

@RequiredArgsConstructor
@RestController
public class HealthCheckController {
    @Operation(summary = "Check if the application is running")
    @GetMapping("/public/api/v1/healthcheck")
    public ResponseEntity<?> checkHealth() {
        return ResponseEntity.ok().build();
    }
}
