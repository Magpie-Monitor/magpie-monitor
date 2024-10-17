package pl.pwr.zpi.utils.client;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.ObjectMapper;
import lombok.extern.slf4j.Slf4j;
import okhttp3.OkHttpClient;
import okhttp3.Request;
import okhttp3.Response;
import org.springframework.stereotype.Component;

import java.io.IOException;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;

@Component
@Slf4j
public class Client implements HttpClient {

    private final OkHttpClient httpClient;
    private final ObjectMapper objectMapper;

    public Client() {
        this.httpClient = new OkHttpClient();
        this.objectMapper = new ObjectMapper();
        this.objectMapper.configure(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES, false);
    }

    @Override
    public <T> T get(String url, Map<String, String> params, Class<T> clazz) {
        String responseBody = sendGetRequest(getUrl(url, params));
        try {
            return objectMapper.readValue(responseBody, clazz);
        } catch (JsonProcessingException e) {
            throw new RuntimeException(e);
        }
    }

    @Override
    public <T> List<T> getList(String url, Map<String, String> params, Class<T> clazz) {
        String responseBody = sendGetRequest(getUrl(url, params));
        try {
            TypeReference<List<T>> typeReference = new TypeReference<>() {
            };
            return objectMapper.readValue(responseBody, typeReference);
        } catch (JsonProcessingException e) {
            throw new RuntimeException(e);
        }
    }

    private String getUrl(String baseUrl, Map<String, String> params) {
        String queryParams = params.entrySet().stream()
                .map(entry -> entry.getKey() + "=" + entry.getValue())
                .collect(Collectors.joining("&"));
        return baseUrl + "?" + queryParams;
    }

    private String sendGetRequest(String url) {
        Request request = new Request.Builder().url(url).build();

        log.info("Sending GET request to: {}", url);

        try (Response response = httpClient.newCall(request).execute()) {
            if (!response.isSuccessful()) {
                log.error("Failed to fetch the resource. Status: {}", response.code());
                throw new RuntimeException("Failed to fetch the resource");
            }

            String responseBody = response.body().string();
            log.info("Response received: {}", responseBody);
            return responseBody;
        } catch (IOException e) {
            log.error("Error fetching resource: {}", e.getMessage(), e);
            throw new RuntimeException("Error fetching resource", e);
        }
    }
}
