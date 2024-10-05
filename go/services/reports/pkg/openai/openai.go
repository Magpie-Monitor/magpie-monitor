package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/IBM/fp-go/array"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/jsonl"
	"go.uber.org/zap"
)

type Client struct {
	model   string
	authKey string
	apiUrl  string
	logger  *zap.Logger
}

type CompletionRequest struct {
	Model          string     `json:"model"`
	Messages       []*Message `json:"messages"`
	Temperature    float32    `json:"temperature"`
	ResponseFormat any        `json:"response_format"`
}

type CompletionResponse struct {
	Id      string
	Object  string
	Model   string
	Usage   UsageStats `json:"usage"`
	Choices []Choice   `json:"choices"`
}

type FileApiReponse struct {
	Id        string `json:"id"`
	Object    string `json:"object"`
	Bytes     int64  `json:"bytes"`
	CreatedAt int64  `json:"created_at"`
	Filename  string `json:"filename"`
	Purpose   string `json:"purpose"`
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

type BatchFileCompletionEntry struct {
	CustomId string             `json:"custom_id"`
	Method   string             `json:"method"`
	Url      string             `json:"url"`
	Body     *CompletionRequest `json:"body"`
}

const (
	COMPLETION_PATH string = "/v1/chat/completions"
)

const (
	FILES_API_PATH string = "/v1/files"
)

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

func (c *Client) Model() string {
	return c.model
}

func (c *Client) addAuthHeaders(r *http.Request) *http.Request {
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.authKey))
	return r
}

func (c *Client) addEncodingHeaders(r *http.Request) *http.Request {
	r.Header.Add("Content-Type", "application/json")
	return r
}

func (c *Client) completionUrl() string {
	return fmt.Sprintf("%s/%s", c.apiUrl, COMPLETION_PATH)
}

func (c *Client) filesApiUrl() string {
	return fmt.Sprintf("%s/%s", c.apiUrl, FILES_API_PATH)
}

func (c *Client) Complete(messages []*Message, responseFormat any) (*CompletionResponse, error) {
	httpClient := http.Client{}

	completionRequest := CompletionRequest{
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
	messageReader := strings.NewReader(string(encodedMessages))
	req, err := http.NewRequest("POST", c.completionUrl(), messageReader)
	if err != nil {
		c.logger.Error("Failed to create request to openai completions api", zap.Error(err))
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

	var decodedResponse CompletionResponse

	body, err := io.ReadAll(res.Body)
	if err != nil {
		c.logger.Error("Failed to read body", zap.Error(err))
		return nil, err
	}

	err = json.Unmarshal(body, &decodedResponse)
	if err != nil {
		c.logger.Error("Failed to decode response from openai api", zap.Error(err))
		return nil, err
	}

	return &decodedResponse, nil
}

func (c *Client) uploadBatch(batchInputFile string) (*FileApiReponse, error) {

	httpClient := http.Client{}
	var fileBuf bytes.Buffer
	mpw := multipart.NewWriter(&fileBuf)
	fileWriter, err := mpw.CreateFormFile("file", "batchinput.jsonl")
	if err != nil {
		c.logger.Error("Failed to add batch file", zap.Error(err))
	}

	_, err = fileWriter.Write([]byte(batchInputFile))
	if err != nil {
		c.logger.Error("Failed to set file header in openai batch")
		return nil, err
	}

	fieldWriter, err := mpw.CreateFormField("purpose")
	if err != nil {
		c.logger.Error("Failed to set purpose header in openai batch")
		return nil, err
	}

	_, err = fieldWriter.Write([]byte("batch"))
	if err != nil {
		c.logger.Error("Failed to set purpose form field")
		return nil, err
	}

	mpw.Close()

	req, err := http.NewRequest("POST", c.filesApiUrl(), &fileBuf)
	if err != nil {
		c.logger.Error("Failed to create request to openai files api", zap.Error(err))
		return nil, err
	}

	c.addAuthHeaders(req)
	req.Header.Add("Content-Type", mpw.FormDataContentType())

	res, err := httpClient.Do(req)
	if err != nil {
		c.logger.Error("Failed to create batch", zap.Error(err))
		return nil, err
	}

	defer res.Body.Close()

	var decodedResponse FileApiReponse
	body, err := io.ReadAll(res.Body)
	if err != nil {
		c.logger.Error("Failed to read files api response body", zap.Error(err))
		return nil, err
	}

	err = json.Unmarshal(body, &decodedResponse)
	if err != nil {
		c.logger.Error("Failed to decode files api response")
		return nil, err
	}

	return &decodedResponse, nil
}

func (c *Client) CreateBatch(completionRequests []*CompletionRequest) (*FileApiReponse, error) {

	batchRequestEntries := array.Map(func(req *CompletionRequest) *BatchFileCompletionEntry {
		return &BatchFileCompletionEntry{
			CustomId: fmt.Sprintf("report-%s", time.Now().String()),
			Method:   "POST",
			Url:      COMPLETION_PATH,
			Body:     req,
		}
	})(completionRequests)

	batchFile := bytes.NewBufferString("")
	err := jsonl.NewJsonLinesEncoder(batchFile).Encode(batchRequestEntries)

	if err != nil {
		c.logger.Sugar().Errorln("Failed to encode messages to complete", zap.Error(err))
		return nil, err
	}

	batchResponse, err := c.uploadBatch(batchFile.String())
	if err != nil {
		c.logger.Error("Failed to create batch", zap.Error(err))
		return nil, err
	}

	return batchResponse, nil
}
