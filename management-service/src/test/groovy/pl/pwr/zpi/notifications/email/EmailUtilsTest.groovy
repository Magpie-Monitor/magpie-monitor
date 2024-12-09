package pl.pwr.zpi.notifications.email

import pl.pwr.zpi.notifications.common.ResourceLoaderUtils
import pl.pwr.zpi.notifications.email.html.service.MarkdownService
import pl.pwr.zpi.notifications.email.internalization.service.LocalizedMessageService
import pl.pwr.zpi.notifications.email.utils.EmailUtils
import spock.lang.Specification
import spock.lang.Subject

class EmailUtilsTest extends Specification {

    def markdownService = Mock(MarkdownService)
    def localizedTestMailServiceImpl = Mock(LocalizedMessageService)
    def localizedNewReportMailServiceImpl = Mock(LocalizedMessageService)

    @Subject
    def emailUtils

    def setup() {
        emailUtils = new EmailUtils(markdownService, localizedTestMailServiceImpl, localizedNewReportMailServiceImpl)
    }

    def "createNewReportEmailTemplate should create a valid new report email template"() {
        given:
        def bodyText = "New report available."
        def buttonText = "View Report"
        def urlRedirect = "https://magpie-monitor.rolo-labs.xyz"
        localizedNewReportMailServiceImpl.getMessage("new-report.body", _) >> bodyText
        localizedNewReportMailServiceImpl.getMessage("new-report.button", _) >> buttonText
        markdownService.toHtmlWithMarkdowns(bodyText) >> "<html><body>${bodyText}</body></html>"

        when:
        def result = emailUtils.createNewReportEmailTemplate(urlRedirect)

        then:
        result.contains("<html><body>${bodyText}</body></html>")
        result.contains(buttonText)
        result.contains(urlRedirect)
    }

    def "createTestEmailTemplate should generate a valid test email template with localized content"() {
        given:
        def bodyText = "Test email body"
        localizedTestMailServiceImpl.getMessage("test.body", _) >> bodyText
        markdownService.toHtmlWithMarkdowns(bodyText) >> "<html><body>${bodyText}</body></html>"

        when:
        def result = emailUtils.createTestEmailTemplate()

        then:
        result.contains("<html><body>${bodyText}</body></html>")
        !result.contains("Click Here")
    }

    def "wrapContent should wrap the content using the wrapper template"() {
        given:
        def htmlTemplate = "<html><body>Original Content</body></html>"
        def wrapperTemplate = "<html><body>%content%</body></html>"
        ResourceLoaderUtils.loadResourceToString(EmailUtils.HTML_WRAPPER_PATH) >> wrapperTemplate

        when:
        def result = emailUtils.wrapContent(htmlTemplate)

        then:
        result.contains("<html><body>Original Content</body></html>")
    }
}
