package pl.pwr.zpi.reports;

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

@Component
@Slf4j
public class ReportsClient {

    private final OkHttpClient httpClient;
    private final ObjectMapper objectMapper;

    public ReportsClient() {
        this.httpClient = new OkHttpClient();
        this.objectMapper = new ObjectMapper();
        this.objectMapper.configure(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES, false);
    }

    public <T> T sendGetRequest(String urlToCall, Class<T> responseType) throws Exception {
        String responseBody = executeHttpGet(urlToCall);
        return objectMapper.readValue(responseBody, responseType);
    }

    public <T> List<T> sendGetRequestForList(String urlToCall, TypeReference<List<T>> typeReference) throws Exception {
        String responseBody = executeHttpGet(urlToCall);
        return objectMapper.readValue(responseBody, typeReference);
    }

    private String executeHttpGet(String urlToCall) throws Exception {
        Request request = new Request.Builder()
                .url(urlToCall)
                .build();

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
