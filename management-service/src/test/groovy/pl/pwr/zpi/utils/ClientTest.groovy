package pl.pwr.zpi.utils

import com.fasterxml.jackson.core.type.TypeReference
import com.fasterxml.jackson.databind.ObjectMapper
import okhttp3.Call
import okhttp3.MediaType
import okhttp3.OkHttpClient
import okhttp3.Response
import okhttp3.ResponseBody
import pl.pwr.zpi.utils.client.Client
import pl.pwr.zpi.utils.exception.JsonMappingException
import spock.lang.Specification

import java.lang.reflect.Field

class ClientTest extends Specification {

    Client client
    OkHttpClient mockHttpClient
    ObjectMapper objectMapper

    def setup() {
        mockHttpClient = Mock(OkHttpClient)
        objectMapper = new ObjectMapper()
        client = new Client()

        Field httpClientField = Client.getDeclaredField("httpClient")
        httpClientField.accessible = true
        httpClientField.set(client, mockHttpClient)

        Field objectMapperField = Client.getDeclaredField("objectMapper")
        objectMapperField.accessible = true
        objectMapperField.set(client, objectMapper)
    }

    def "get should return mapped object on valid response"() {
        given:
        String url = "https://example.com/api/resource"
        Map<String, String> params = [param1: "value1", param2: "value2"]
        String jsonResponse = '{"key":"value"}'
        def responseMock = Mock(Response) {
            isSuccessful() >> true
            body() >> ResponseBody.create(MediaType.get("application/json"), jsonResponse)
        }
        mockHttpClient.newCall(_) >> Mock(Call) {
            execute() >> responseMock
        }

        when:
        Map result = client.get(url, params, Map)

        then:
        result.key == "value"
    }

    def "getList should return list of mapped objects on valid response"() {
        given:
        String url = "https://example.com/api/resource"
        Map<String, String> params = [:]
        String jsonResponse = '[{"key":"value1"}, {"key":"value2"}]'
        def responseMock = Mock(Response) {
            isSuccessful() >> true
            body() >> ResponseBody.create(MediaType.get("application/json"), jsonResponse)
        }
        mockHttpClient.newCall(_) >> Mock(Call) {
            execute() >> responseMock
        }

        when:
        List<Map> result = client.getList(url, params, new TypeReference<List<Map>>() {})

        then:
        result.size() == 2
        result[0].key == "value1"
        result[1].key == "value2"
    }

    def "get should throw JsonMappingException on invalid JSON"() {
        given:
        String url = "https://example.com/api/resource"
        Map<String, String> params = [:]
        String invalidJsonResponse = 'invalid-json'
        def responseMock = Mock(Response) {
            isSuccessful() >> true
            body() >> ResponseBody.create(MediaType.get("application/json"), invalidJsonResponse)
        }
        mockHttpClient.newCall(_) >> Mock(Call) {
            execute() >> responseMock
        }

        when:
        client.get(url, params, Map)

        then:
        thrown(JsonMappingException)
    }

    def "sendGetRequest should throw RuntimeException on unsuccessful response"() {
        given:
        String url = "https://example.com/api/resource"
        def responseMock = Mock(Response) {
            isSuccessful() >> false
            code() >> 500
        }
        mockHttpClient.newCall(_) >> Mock(Call) {
            execute() >> responseMock
        }

        when:
        client.get(url, [:], Map)

        then:
        RuntimeException ex = thrown()
        ex.message == "Failed to fetch the resource"
    }

    def "sendGetRequest should throw RuntimeException on IOException"() {
        given:
        String url = "https://example.com/api/resource"
        mockHttpClient.newCall(_) >> Mock(Call) {
            execute() >> { throw new IOException("Test IO Exception") }
        }

        when:
        client.get(url, [:], Map)

        then:
        RuntimeException ex = thrown()
        ex.message == "Error fetching resource"
    }
}