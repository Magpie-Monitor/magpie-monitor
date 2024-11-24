package notifications.common

import pl.pwr.zpi.notifications.common.ConfidentialTextEncoder
import spock.lang.Specification
import spock.lang.Subject

class ConfidentialTextEncoderTest extends Specification {

    @Subject
    ConfidentialTextEncoder confidentialTextEncoder

    def setup() {
        String exampleEncryptionKey = "O8JvErGt84wzZzPPeFg4tQ=="
        String secretKeyCode = Base64.getEncoder().encodeToString(exampleEncryptionKey.getBytes("UTF-8"))
        String cipherAlgorithm = "AES"

        confidentialTextEncoder = new ConfidentialTextEncoder(secretKeyCode, cipherAlgorithm)
    }

    private String processMessage(String plainText) {
        String encryptedText = confidentialTextEncoder.encrypt(plainText)
        String decryptedText = confidentialTextEncoder.decrypt(encryptedText)
        return decryptedText
    }

    def "should encrypt and decrypt messages correctly"() {
        expect:
        processMessage(plainText) == plainText

        where:
        plainText << [
                "Example confidential information",
                "123454321",
                "https://discord.com/api/webhooks/1234554321/xKh5vF0Som55bSex4q9slwOApmB0VXjcUoVS5Z9v9vu89snl-XeedfHj",
                ""
        ]
    }
}
