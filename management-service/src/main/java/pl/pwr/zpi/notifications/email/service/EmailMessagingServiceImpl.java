package pl.pwr.zpi.notifications.email.service;

import jakarta.mail.internet.MimeMessage;
import lombok.RequiredArgsConstructor;
import lombok.SneakyThrows;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.core.io.ClassPathResource;
import org.springframework.mail.javamail.JavaMailSender;
import org.springframework.mail.javamail.MimeMessageHelper;
import org.springframework.scheduling.annotation.Async;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.notifications.common.ResourceLoaderUtils;

@Service
@RequiredArgsConstructor
public class EmailMessagingServiceImpl implements EmailMessagingService {

    @Value("${spring.mail.username}")
    private String EMAIL_FROM;

    private final JavaMailSender mailSender;

    @Async
    @Override
    @SneakyThrows
    public void sendMessage(String receiver, String subject, String body, boolean isHtml) {
        MimeMessage mimeMessage = mailSender.createMimeMessage();
        MimeMessageHelper helper = new MimeMessageHelper(mimeMessage, true, "UTF-8");

        helper.setFrom(EMAIL_FROM);
        helper.setTo(receiver);
        helper.setSubject(subject);
        helper.setText(body, isHtml);
//        helper.addInline("magpie-monitor-icon", ResourceLoaderUtils.loadResourceByteArray("templates/email/assets/magpie-monitor-icon.png"));
//        helper.addInline("github-logo-icon", ResourceLoaderUtils.loadResourceByteArray("templates/email/assets/github-logo-icon.png"));
        ClassPathResource magpieMonitorIcon = new ClassPathResource("templates/email/assets/magpie-monitor-icon.png");
        ClassPathResource githubLogoIcon = new ClassPathResource("templates/email/assets/github-logo-icon.png");

        helper.addInline("magpie-monitor-icon", magpieMonitorIcon);
        helper.addInline("github-logo-icon", githubLogoIcon);

        mailSender.send(mimeMessage);
    }
}
