package pl.pwr.zpi.utils.client;

import com.fasterxml.jackson.core.type.TypeReference;

import java.util.List;
import java.util.Map;

public interface HttpClient {
    <T> T get(String url, Map<String, String> params, Class<T> clazz);
    <T> List<T> getList(String url, Map<String, String> params, TypeReference<List<T>> clazz);
}
