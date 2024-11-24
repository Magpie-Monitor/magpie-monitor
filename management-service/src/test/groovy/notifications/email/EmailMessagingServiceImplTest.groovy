package notifications.email

import jakarta.mail.internet.MimeMessage
import org.springframework.mail.javamail.JavaMailSender
import pl.pwr.zpi.notifications.email.service.EmailMessagingServiceImpl
import spock.lang.Specification
import spock.lang.Subject

class EmailMessagingServiceImplTest extends Specification {

    JavaMailSender mailSender = Mock(JavaMailSender)
    MimeMessage mimeMessage = Mock(MimeMessage)

    @Subject
    EmailMessagingServiceImpl emailMessagingService = new EmailMessagingServiceImpl(mailSender)

    def setup() {
        emailMessagingService.EMAIL_FROM = "test@example.com"
    }
}
