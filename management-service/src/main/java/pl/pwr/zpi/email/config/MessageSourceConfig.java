package pl.pwr.zpi.email.config;

import org.springframework.context.MessageSource;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.support.ReloadableResourceBundleMessageSource;

@Configuration
public class MessageSourceConfig {

    @Bean
    public MessageSource messageSource() {
        return getMessageSource("messages/messages");
    }

    @Bean
    public MessageSource testMailSource() {
        return getMessageSource("emails/test/test");
    }

    @Bean
    public MessageSource newReportMailSource() {
        return getMessageSource("emails/newReport/new_report");
    }

    private MessageSource getMessageSource(String path) {
        var messageSource = new ReloadableResourceBundleMessageSource();
        messageSource.setBasenames("classpath:" + path, path);
        messageSource.setDefaultEncoding("UTF-8");
        return messageSource;
    }
}
