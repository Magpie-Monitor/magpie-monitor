package pl.pwr.zpi.notifications.common;

import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

import javax.crypto.BadPaddingException;
import javax.crypto.Cipher;
import javax.crypto.IllegalBlockSizeException;
import javax.crypto.SecretKey;
import javax.crypto.spec.SecretKeySpec;
import java.nio.charset.StandardCharsets;
import java.security.InvalidKeyException;
import java.util.Base64;

@Slf4j
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

    public String encrypt(String message) {
        initCipher();
        byte[] encryptedMessage = null;
        try {
            encryptedMessage = cipher.doFinal(message.getBytes(StandardCharsets.UTF_8));
        } catch (IllegalBlockSizeException | BadPaddingException e) {
            log.error("Encryption error: {}", e.getMessage());
            throw new RuntimeException(e);
        }
        return Base64.getEncoder().encodeToString(encryptedMessage);
    }

    private void initCipher() {
        try {
            cipher.init(Cipher.ENCRYPT_MODE, secretKey);
        } catch (InvalidKeyException e) {
            log.error("Cipher init error: {}", e.getMessage());
            throw new RuntimeException(e);
        }
    }

    public String decrypt(String encryptedMessage) {
        initCipher();
        byte[] decryptedMessage = null;
        try {
            decryptedMessage = cipher.doFinal(Base64.getDecoder().decode(encryptedMessage));
        } catch (IllegalBlockSizeException | BadPaddingException e) {
            log.error("Decryption error: {}", e.getMessage());
            throw new RuntimeException(e);
        }
        return new String(decryptedMessage, StandardCharsets.UTF_8);
    }
}
