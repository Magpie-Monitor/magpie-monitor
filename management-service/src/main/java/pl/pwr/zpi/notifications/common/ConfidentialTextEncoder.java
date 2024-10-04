package pl.pwr.zpi.notifications.common;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

import javax.crypto.Cipher;
import javax.crypto.SecretKey;
import javax.crypto.spec.SecretKeySpec;
import java.nio.charset.StandardCharsets;
import java.util.Base64;

@Component
public class ConfidentialTextEncoder {
    private final SecretKey secretKey;
    private final Cipher cipher;

    public ConfidentialTextEncoder(@Value("${encoding.cipher.secret.key}") String secretKeyCode,
                                   @Value("${encoding.cipher.algorithm}") String cipherAlgorithm) throws Exception {
        byte[] keyBytes = Base64.getDecoder().decode(secretKeyCode);
        this.secretKey = new SecretKeySpec(keyBytes, 0, keyBytes.length, cipherAlgorithm);
        this.cipher = Cipher.getInstance(cipherAlgorithm);
    }

    public String encrypt(String message) throws Exception {
        cipher.init(Cipher.ENCRYPT_MODE, secretKey);
        byte[] encryptedMessage = cipher.doFinal(message.getBytes(StandardCharsets.UTF_8));
        return Base64.getEncoder().encodeToString(encryptedMessage);
    }

    public String decrypt(String encryptedMessage) throws Exception {
        cipher.init(Cipher.DECRYPT_MODE, secretKey);
        byte[] decryptedMessage = cipher.doFinal(Base64.getDecoder().decode(encryptedMessage));
        return new String(decryptedMessage, StandardCharsets.UTF_8);
    }
}
