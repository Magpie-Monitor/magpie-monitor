package notifications.slack

import pl.pwr.zpi.notifications.slack.service.SlackMessagingServiceImpl
import spock.lang.Specification
import spock.lang.Subject

class SlackMessagingServiceImplTest extends Specification {

    @Subject
    SlackMessagingServiceImpl slackMessagingService = new SlackMessagingServiceImpl()

    def "should send message to Slack successfully"() {
    }

    def "should handle IOException when sending message"() {
    }
}