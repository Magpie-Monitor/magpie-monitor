package pl.pwr.zpi.utils.mapper;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.springframework.stereotype.Component;
import pl.pwr.zpi.utils.exception.JsonMappingException;

@Component
public class JsonMapper {

    private final ObjectMapper objectMapper;

    public JsonMapper() {
        this.objectMapper = new ObjectMapper();
    }

    public <T> T fromJson(String json, Class<T> clazz) {
        try {
            return objectMapper.readValue(json, clazz);
        } catch (JsonProcessingException e) {
            throw new JsonMappingException(e.getMessage());
        }
    }

}
