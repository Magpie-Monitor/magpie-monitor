package notifications.email

import pl.pwr.zpi.notifications.email.service.EmailReceiverService
import pl.pwr.zpi.notifications.email.dto.EmailReceiverDTO
import pl.pwr.zpi.notifications.email.entity.EmailReceiver
import pl.pwr.zpi.notifications.email.repository.EmailRepository
import spock.lang.Ignore
import spock.lang.Specification
import spock.lang.Subject

@Ignore
class EmailReceiverServiceTest extends Specification {

    def emailRepository = Mock(EmailRepository)

    @Subject
    def emailReceiverService = new EmailReceiverService(emailRepository)

    def "should get all email receivers"() {
        given:
        def emailReceiverList = [buildEmailReceiver(1L, "Receiver 1", "receiver1@example.com")]

        emailRepository.findAll() >> emailReceiverList

        when:
        def result = emailReceiverService.getAllEmails()

        then:
        result == emailReceiverList
        1 * emailRepository.findAll()
    }

    def "should add a new email receiver"() {
        given:
        def emailReceiverDTO = buildEmailReceiverDTO("Receiver 1", "receiver1@example.com")
        def newReceiver = buildEmailReceiver(1L, emailReceiverDTO.getName(), emailReceiverDTO.getEmail())

        emailRepository.existsByReceiverEmail(emailReceiverDTO.getEmail()) >> false
        emailRepository.save(_) >> newReceiver

        when:
        emailReceiverService.addNewEmail(emailReceiverDTO)

        then:
        1 * emailRepository.existsByReceiverEmail(emailReceiverDTO.getEmail())
        1 * emailRepository.save(newReceiver)
    }

    def "should throw exception when adding email that already exists"() {
        given:
        def emailReceiverDTO = buildEmailReceiverDTO("Receiver 1", "receiver1@example.com")

        emailRepository.existsByReceiverEmail(emailReceiverDTO.getEmail()) >> true

        when:
        emailReceiverService.addNewEmail(emailReceiverDTO)

        then:
        thrown(IllegalArgumentException)
        1 * emailRepository.existsByReceiverEmail(emailReceiverDTO.getEmail())
    }

    def "should update email receiver successfully"() {
        given:
        def id = 1L
        def emailReceiverDTO = buildEmailReceiverDTO("Updated Receiver", "updatedreceiver@example.com")
        def existingReceiver = buildEmailReceiver(id, "Receiver 1", "receiver1@example.com")

        emailRepository.findById(id) >> Optional.of(existingReceiver)
        emailRepository.existsByReceiverEmail(emailReceiverDTO.getEmail()) >> false
        emailRepository.save(_) >> existingReceiver

        when:
        def updatedReceiver = emailReceiverService.updateEmail(id, emailReceiverDTO)

        then:
        updatedReceiver.receiverName == "Updated Receiver"
        updatedReceiver.receiverEmail == "updatedreceiver@example.com"
        1 * emailRepository.findById(id)
        1 * emailRepository.existsByReceiverEmail(emailReceiverDTO.getEmail()) >> false
        1 * emailRepository.save(_ as EmailReceiver)
    }

    def "should throw exception when trying to update email with already existing email"() {
        given:
        def id = 1L
        def emailReceiverDTO = buildEmailReceiverDTO("Receiver 1", "receiver1@example.com")
        def existingReceiver = buildEmailReceiver(id, "Receiver 1", "receiver1@example.com")

        emailRepository.findById(id) >> Optional.of(existingReceiver)
        emailRepository.existsByReceiverEmail(emailReceiverDTO.getEmail()) >> true

        when:
        emailReceiverService.updateEmail(id, emailReceiverDTO)

        then:
        thrown(IllegalArgumentException)
        1 * emailRepository.findById(id)
        1 * emailRepository.existsByReceiverEmail(emailReceiverDTO.getEmail())
    }

    def "should throw exception when email receiver not found for update"() {
        given:
        def id = 1L
        def emailReceiverDTO = buildEmailReceiverDTO("Receiver 1", "receiver1@example.com")

        emailRepository.findById(id) >> Optional.empty()

        when:
        emailReceiverService.updateEmail(id, emailReceiverDTO)

        then:
        thrown(IllegalArgumentException)
        1 * emailRepository.findById(id)
    }

    private EmailReceiverDTO buildEmailReceiverDTO(String name, String email) {
        return EmailReceiverDTO.builder()
                .name(name)
                .email(email)
                .build()
    }

    private EmailReceiver buildEmailReceiver(Long id, String name, String email) {
        return EmailReceiver.builder()
                .id(id)
                .receiverName(name)
                .receiverEmail(email)
                .createdAt(System.currentTimeMillis())
                .updatedAt(System.currentTimeMillis())
                .build()
    }
}
