package openai

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/repositories"
	"go.uber.org/zap"
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

func CreateIncidentReportSchema() ResponseFormat {
	reports := repositories.Report{}

	return ResponseFormat{
		Type: "json_schema",
		JsonSchema: JsonSchema{
			Name:   "incident_report",
			Schema: *getSchemaFromStruct(reports),
		},
	}
}

type Client struct {
	model   string
	authKey string
	apiUrl  string
	logger  *zap.Logger
}

type Request struct {
	Model          string         `json:"model"`
	Messages       []*Message     `json:"messages"`
	Temperature    float32        `json:"temperature"`
	ResponseFormat ResponseFormat `json:"response_format"`
}

type Response struct {
	Id      string
	Object  string
	Model   string
	Usage   UsageStats `json:"usage"`
	Choices []Choice   `json:"choices"`
}

type UsageStats struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type Choice struct {
	Message      *Message `json:"message"`
	FinishReason string   `json:"finish_reason"`
	Index        int      `json:"index"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type BatchFileEntry struct {
	CustomId string   `json:"custom_id"`
	Method   string   `json:"method"`
	Url      string   `json:"url"`
	Body     *Request `json:"body"`
}

func NewOpenAiClient(logger *zap.Logger) *Client {

	apiUrl := os.Getenv("OPENAI_API_URL")
	model := os.Getenv("OPENAI_API_MODEL")
	authKey := os.Getenv("OPENAI_API_KEY")

	return &Client{
		model:   model,
		authKey: authKey,
		apiUrl:  apiUrl,
		logger:  logger,
	}
}

func (c *Client) addAuthHeaders(r *http.Request) *http.Request {
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.authKey))
	return r
}

func (c *Client) addEncodingHeaders(r *http.Request) *http.Request {
	r.Header.Add("Content-Type", "application/json")
	return r
}

func (c *Client) Complete(messages []*Message) (*Response, error) {
	httpClient := http.Client{}

	completionRequest := Request{
		Model:          c.model,
		Messages:       messages,
		Temperature:    0.6,
		ResponseFormat: CreateIncidentReportSchema(),
	}

	encodedMessages, err := json.Marshal(completionRequest)
	if err != nil {
		c.logger.Sugar().Errorln("Failed to encode messages to complete", zap.Error(err))
		return nil, err
	}
	c.logger.Debug("Encoded messages", zap.Any("msg", string(encodedMessages)))
	messageReader := strings.NewReader(string(encodedMessages))

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/chat/completions", c.apiUrl), messageReader)
	if err != nil {
		c.logger.Error("Failed to create request to /chat/completions", zap.Error(err))
		return nil, err
	}

	c.addAuthHeaders(req)
	c.addEncodingHeaders(req)

	res, err := httpClient.Do(req)
	if err != nil {
		c.logger.Error("Failed to make a request to openai api", zap.Error(err))
		return nil, err
	}

	defer res.Body.Close()

	var decodedResponse Response

	body, err := io.ReadAll(res.Body)
	if err != nil {
		c.logger.Error("Failed to read body", zap.Error(err))
		return nil, err
	}

	c.logger.Debug("Response", zap.Any("response", string(body)))

	err = json.Unmarshal(body, &decodedResponse)
	if err != nil {
		c.logger.Error("Failed to decode response from openai api", zap.Error(err))
		return nil, err
	}

	return &decodedResponse, nil

}
