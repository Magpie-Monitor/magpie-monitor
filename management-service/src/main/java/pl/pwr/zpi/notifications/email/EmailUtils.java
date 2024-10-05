package pl.pwr.zpi.notifications.email;

import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Component;
import pl.pwr.zpi.notifications.common.ResourceLoaderUtils;
import pl.pwr.zpi.notifications.email.html.HtmlBuilder;
import pl.pwr.zpi.notifications.email.html.service.MarkdownService;
import pl.pwr.zpi.notifications.email.iternalization.service.LocalizedMessageService;

import java.time.LocalDateTime;
import java.util.Map;

@RequiredArgsConstructor
@Component
public class EmailUtils {

    private final MarkdownService markdownService;
    private final LocalizedMessageService localizedTestMailServiceImpl;
    private final LocalizedMessageService localizedNewReportMailServiceImpl;

    private static final String HTML_HEADER_PATH = "templates/email/banner.html";
    private static final String HTML_FOOTER_PATH = "templates/email/footer.html";
    private static final String HTML_WRAPPER_PATH = "templates/email/wrapper.html";

    public String createTextMailTemplate(String markdownText) {
        return createTemplate(markdownService.toHtmlWithMarkdowns(markdownText), null, null);
    }

    public String createTestEmailTemplate() {
        return createLocalizedTemplate(localizedTestMailServiceImpl, "test.body");
    }

    public String createNewReportEmailTemplate(String urlRedirect) {
        return createLocalizedTemplate(localizedNewReportMailServiceImpl,"new-report.body", "new-report.button", urlRedirect);
    }

    private String createLocalizedTemplate(LocalizedMessageService messageService, String bodyKey, String buttonKey, String urlRedirect) {
        String text = messageService.getMessage(
                bodyKey,
                messageService.getLanguageFromContextOrDefault()
        );
        String markdownHtml = markdownService.toHtmlWithMarkdowns(text);
        String buttonText = messageService.getMessage(
                buttonKey,
                messageService.getLanguageFromContextOrDefault()
        );
        return createTemplate(markdownHtml, buttonText, urlRedirect);
    }

    private String createLocalizedTemplate(LocalizedMessageService messageService, String bodyKey) {
        String text = messageService.getMessage(
                bodyKey,
                messageService.getLanguageFromContextOrDefault()
        );
        String markdownHtml = markdownService.toHtmlWithMarkdowns(text);

        return createTemplate(markdownHtml, null, null);
    }

    private String createTemplate(String markdownHtml, String buttonText, String urlRedirect) {
        EmailTextBuilder builder = new EmailTextBuilder().withText(markdownHtml);
        if (buttonText != null && urlRedirect != null) {
            builder.withButton(buttonText, urlRedirect);
        }
        return wrapIntoEmail(builder.build());
    }

    private String wrapContent(String htmlTemplate) {
        HtmlBuilder htmlBuilder = new HtmlBuilder();
        Map<String, String> textKeys = Map.of("%content%", htmlTemplate);
        return htmlBuilder.appendAndReplace(ResourceLoaderUtils.loadResourceToString(HTML_WRAPPER_PATH), textKeys).build();
    }

    private String wrapIntoEmail(String htmlBody) {
        HtmlBuilder htmlBuilder = new HtmlBuilder();
        String htmlHeader = ResourceLoaderUtils.loadResourceToString(HTML_HEADER_PATH);
        String htmlFooter = createFooter();
        String content = htmlBuilder.append(htmlHeader).append(htmlBody).append(htmlFooter).build();
        return wrapContent(content);
    }

    private String createFooter() {
        String footer = ResourceLoaderUtils.loadResourceToString(HTML_FOOTER_PATH);
        return HtmlBuilder.replace(footer, Map.of("%random_text%", LocalDateTime.now().toString()));
    }
}
