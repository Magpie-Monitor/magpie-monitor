package pl.pwr.zpi.notifications.email.internalization;

import lombok.Getter;

import java.util.Arrays;
import java.util.Locale;
import java.util.Optional;

@Getter
public enum SupportedLanguage {

    EN(Locale.of("en"));

    private final Locale locale;

    SupportedLanguage(Locale locale) {
        this.locale = locale;
    }

    public static Optional<SupportedLanguage> fromLocale(Locale locale) {
        return Arrays.stream(SupportedLanguage.values())
                .filter(language -> language.getLocale().getLanguage().equals(locale.getLanguage()))
                .findFirst();
    }
}
