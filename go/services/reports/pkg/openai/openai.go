package openai

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"strings"
)

type Client struct {
	model   string
	authKey string
	apiUrl  string
	logger  *zap.Logger
}

type Request struct {
	Model          string     `json:"model"`
	Messages       []*Message `json:"messages"`
	Temperature    float32    `json:"temperature"`
	ResponseFormat any        `json:"response_format"`
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

func (c *Client) Complete(messages []*Message, responseFormat any) (*Response, error) {
	httpClient := http.Client{}

	completionRequest := Request{
		Model:          c.model,
		Messages:       messages,
		Temperature:    0.6,
		ResponseFormat: responseFormat,
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
