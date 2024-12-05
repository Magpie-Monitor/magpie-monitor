package pl.pwr.zpi.utils

import pl.pwr.zpi.utils.exception.JsonMappingException
import pl.pwr.zpi.utils.mapper.JsonMapper
import spock.lang.Specification

class JsonMapperTest extends Specification {

    JsonMapper jsonMapper

    def setup() {
        jsonMapper = new JsonMapper()
    }

    def "fromJson should correctly deserialize a valid JSON string into an object"() {
        given:
        String json = '{"name": "John Doe", "age": 30}'

        when:
        Person person = jsonMapper.fromJson(json, Person)

        then:
        person.name == "John Doe"
        person.age == 30
    }

    def "fromJson should throw JsonMappingException when JSON string is invalid"() {
        given:
        String invalidJson = '{name: "John Doe", age: }'

        when:
        jsonMapper.fromJson(invalidJson, Person)

        then:
        JsonMappingException ex = thrown()
        ex.message.contains("Unexpected character")
    }

    def "fromJson should throw JsonMappingException for incompatible JSON and target class"() {
        given:
        String json = '{"name": "John Doe", "age": 30}'

        when:
        jsonMapper.fromJson(json, DifferentClass)

        then:
        JsonMappingException ex = thrown()
        ex.message.contains("Unrecognized field")
    }


    static class Person {
        String name
        int age
    }

    static class DifferentClass {
        String firstName
    }
}