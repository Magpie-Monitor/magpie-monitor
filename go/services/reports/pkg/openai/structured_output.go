package openai

import (
	"reflect"
	"strings"
)

type ResponseFormat struct {
	Type       string     `json:"type" bson:"type"`
	JsonSchema JsonSchema `json:"json_schema" bson:"json_schema"`
}

type JsonSchema struct {
	Name   string `json:"name" bson:"name"`
	Schema Schema `json:"schema" bson:"schema"`
	Strict bool   `json:"strict" bson:"strict"`
}

type Schema struct {
	Type                 string             `json:"type" bson:"type"`
	Items                *Schema            `json:"items,omitempty" bson:"items,omitempty"`           // Use omitempty
	Properties           map[string]*Schema `json:"properties,omitempty" bson:"properties,omitempty"` // Use omitempty
	Required             []string           `json:"required,omitempty" bson:"required,omitempty"`     // Use omitempty
	AdditionalProperties bool               `json:"additionalProperties" bson:"additionalProperties"`
}

func getSchemaFromStruct(obj interface{}) *Schema {
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if isPrimitive(t.Kind()) {
		schemaType := mapKindToSchema(t.Kind())
		return &Schema{
			Type: schemaType,
		}
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

		// Add field schema to properties
		schema.Properties[fieldName] = createSchemaType(fieldType)

		// Add to required fields if the field does not have the 'omitempty' tag
		if !strings.Contains(field.Tag.Get("json"), "omitempty") {
			schema.Required = append(schema.Required, fieldName)
		}
	}

	return schema
}

func createSchemaType(fieldType reflect.Type) *Schema {
	kind := fieldType.Kind()
	if isPrimitive(kind) {
		return handlePrimitiveSchema(fieldType)
	} else if kind == reflect.Struct {
		return handleObjectSchema(fieldType)
	} else if kind == reflect.Slice {
		return handleArraySchema(fieldType)
	} else {
		return &Schema{Type: "string"} // Default to string for unknown types
	}
}

func handleObjectSchema(fieldType reflect.Type) *Schema {
	return getSchemaFromStruct(reflect.New(fieldType).Interface())
}

func handleArraySchema(fieldType reflect.Type) *Schema {
	return &Schema{
		Type:  "array",
		Items: getSchemaFromStruct(reflect.New(fieldType.Elem()).Interface()),
	}
}

func handlePrimitiveSchema(fieldType reflect.Type) *Schema {
	return &Schema{Type: mapKindToSchema(fieldType.Kind())}
}

func CreateJsonReponseFormat(schemaName string, responseShape any) ResponseFormat {

	return ResponseFormat{
		Type: "json_schema",
		JsonSchema: JsonSchema{
			Name:   "incident_report",
			Strict: true,
			Schema: *getSchemaFromStruct(responseShape),
		},
	}
}

func mapKindToSchema(kind reflect.Kind) string {
	switch kind {
	case reflect.String:
		return "string"
	case reflect.Int, reflect.Int32, reflect.Int64:
		return "integer"
	case reflect.Float32, reflect.Float64:
		return "number"
	case reflect.Bool:
		return "boolean"
	default:
		panic("Non-primitive type type passed to kind mapper")
	}
}

func isPrimitive(kind reflect.Kind) bool {
	return kind != reflect.Slice && kind != reflect.Struct
}
