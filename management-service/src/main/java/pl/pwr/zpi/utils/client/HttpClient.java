package pl.pwr.zpi.utils.client;

import java.util.List;

public interface HttpClient {
    <T> T get(String url, Class<T> clazz);
    <T> List<T> getList(String url, Class<T> clazz);
}
