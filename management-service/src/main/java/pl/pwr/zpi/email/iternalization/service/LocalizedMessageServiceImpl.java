package pl.pwr.zpi.email.iternalization.service;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.MessageSource;
import org.springframework.context.NoSuchMessageException;
import org.springframework.context.i18n.LocaleContextHolder;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.email.iternalization.SupportedLanguage;

@Service
@RequiredArgsConstructor
@Slf4j
public class LocalizedMessageServiceImpl implements LocalizedMessageService {

    @Value("${language.default}")
    private SupportedLanguage DEFAULT_LANGUAGE;

    private final MessageSource messageSource;

    @Override
    public String getMessage(String key, SupportedLanguage language) {
        try {
            return messageSource.getMessage(key, null, language.getLocale());
        } catch (NoSuchMessageException e) {
            log.error("Message {} not found, check if correct key is in resource bundles", key);
            return key;
        }
    }

    @Override
    public String getMessageWithArgs(String key, SupportedLanguage language, Object... args) {
        try {
            return messageSource.getMessage(key, args, language.getLocale());
        } catch (NoSuchMessageException e) {
            log.error("Message {} not found, check if correct key is in resource bundles", key);
            return key;
        }
    }

    public String getMessageFromContext(String key) {
        return getMessage(key, getLanguageFromContextOrDefault());
    }

    public String getMessageWithArgsFromContext(String key, Object... args) {
        return getMessageWithArgs(key, getLanguageFromContextOrDefault(), args);
    }

    @Override
    public SupportedLanguage getLanguageFromContextOrDefault() {
        var language = SupportedLanguage.fromLocale(LocaleContextHolder.getLocale());
        return language.orElseGet(() -> DEFAULT_LANGUAGE);
    }
}
