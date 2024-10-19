package pl.pwr.zpi;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.scheduling.annotation.EnableScheduling;

@SpringBootApplication
@EnableScheduling
public class MagpieMonitorApplication {

    public static void main(String[] args) {
        SpringApplication.run(MagpieMonitorApplication.class, args);
    }
}
