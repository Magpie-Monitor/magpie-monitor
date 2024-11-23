package pl.pwr.zpi.notifications.common;

import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

import javax.crypto.*;
import javax.crypto.spec.SecretKeySpec;
import java.nio.charset.StandardCharsets;
import java.security.InvalidKeyException;
import java.security.NoSuchAlgorithmException;
import java.util.Base64;

@Slf4j
@Component
public class ConfidentialTextEncoder {
    private final SecretKey secretKey;
    private final Cipher cipher;

    public ConfidentialTextEncoder(@Value("${encoding.cipher.secret.key}") String secretKeyCode,
                                   @Value("${encoding.cipher.algorithm}") String cipherAlgorithm) {
        byte[] keyBytes = Base64.getDecoder().decode(secretKeyCode);
        this.secretKey = new SecretKeySpec(keyBytes, 0, keyBytes.length, cipherAlgorithm);
        this.cipher = getCipher(cipherAlgorithm);
    }

    private Cipher getCipher(String cipherAlgorithm) {
        try {
            return Cipher.getInstance(cipherAlgorithm);
        } catch (NoSuchAlgorithmException | NoSuchPaddingException e) {
            log.error("Error creating Cipher instance: {}", e.getMessage());
            throw new RuntimeException(e);
        }
    }

    public String encrypt(String message) {
        cipherInit(Cipher.ENCRYPT_MODE);
        return Base64.getEncoder()
                .encodeToString(
                        cipherDoFinal(message.getBytes(StandardCharsets.UTF_8))
                );
    }

    public String decrypt(String encryptedMessage) {
        cipherInit(Cipher.DECRYPT_MODE);
        return new String(
                cipherDoFinal(Base64.getDecoder().decode(encryptedMessage)),
                StandardCharsets.UTF_8
        );
    }

    private byte[] cipherDoFinal(byte[] input) {
        try {
            return cipher.doFinal(input);
        } catch (IllegalBlockSizeException | BadPaddingException e) {
            log.error("Cipher doFinal error: {}", e.getMessage());
            throw new RuntimeException(e);
        }
    }

    private void cipherInit(int mode) {
        try {
            cipher.init(mode, secretKey);
        } catch (InvalidKeyException e) {
            log.error("Cipher init error: {}", e.getMessage());
            throw new RuntimeException(e);
        }
    }
}
