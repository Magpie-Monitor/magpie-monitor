package openai

import (
	"reflect"
	"strings"
)

type ResponseFormat struct {
	Type       string     `json:"type"`
	JsonSchema JsonSchema `json:"json_schema"`
}

type JsonSchema struct {
	Name   string `json:"name"`
	Schema Schema `json:"schema"`
	Strict bool   `json:"strict"`
}

type Schema struct {
	Type                 string             `json:"type"`
	Items                *Schema            `json:"items,omitempty"`      // Use omitempty
	Properties           map[string]*Schema `json:"properties,omitempty"` // Use omitempty
	Required             []string           `json:"required,omitempty"`   // Use omitempty
	AdditionalProperties bool               `json:"additionalProperties"`
}

func getSchemaFromStruct(obj interface{}) *Schema {
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	schema := &Schema{
		Type:                 "object",
		Properties:           map[string]*Schema{},
		AdditionalProperties: false,
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldType := field.Type
		fieldName := strings.Split(field.Tag.Get("json"), ",")[0]
		if fieldName == "" {
			fieldName = field.Name
		}

		// Determine the schema type of the field
		var fieldSchema *Schema
		switch fieldType.Kind() {
		case reflect.String:
			fieldSchema = &Schema{Type: "string"}
		case reflect.Int, reflect.Int32, reflect.Int64:
			fieldSchema = &Schema{Type: "integer"}
		case reflect.Float32, reflect.Float64:
			fieldSchema = &Schema{Type: "number"}
		case reflect.Bool:
			fieldSchema = &Schema{Type: "boolean"}
		case reflect.Struct:
			// Recursively get the schema of the nested struct
			fieldSchema = getSchemaFromStruct(reflect.New(fieldType).Interface())
		case reflect.Slice:
			// Handle slice (arrays) fields
			fieldSchema = &Schema{
				Type:  "array",
				Items: getSchemaFromStruct(reflect.New(fieldType.Elem()).Interface()),
			}
		default:
			fieldSchema = &Schema{Type: "string"} // Default to string for unknown types
		}

		// Add field schema to properties
		schema.Properties[fieldName] = fieldSchema

		// Add to required fields if the field does not have the 'omitempty' tag
		if !strings.Contains(field.Tag.Get("json"), "omitempty") {
			schema.Required = append(schema.Required, fieldName)
		}
	}

	return schema
}

func CreateJsonReponseFormat(schemaName string, responseShape any) ResponseFormat {

	return ResponseFormat{
		Type: "json_schema",
		JsonSchema: JsonSchema{
			Name:   "incident_report",
			Schema: *getSchemaFromStruct(responseShape),
		},
	}
}
