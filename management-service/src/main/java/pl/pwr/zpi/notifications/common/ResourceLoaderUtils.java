package pl.pwr.zpi.notifications.common;

import lombok.SneakyThrows;
import org.springframework.core.io.*;

import java.nio.charset.Charset;
import java.nio.charset.StandardCharsets;

public final class ResourceLoaderUtils {

    private ResourceLoaderUtils() {}

    private static final ResourceLoader resourceLoader = new DefaultResourceLoader();

    @SneakyThrows
    public static String loadResourceToString(String path) {
        return loadResourceToString(path, StandardCharsets.UTF_8);
    }

    @SneakyThrows
    public static String loadResourceToString(String path, Charset charset) {
        return new String(loadResource(path), charset);
    }

    @SneakyThrows
    public static byte[] loadResource(String path) {
        return resourceLoader.getResource(path).getInputStream().readAllBytes();
    }

    @SneakyThrows
    public static ByteArrayResource loadResourceByteArray(String path) {
        return new ByteArrayResource(loadResource(path));
    }
}
