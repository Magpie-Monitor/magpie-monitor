package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/IBM/fp-go/array"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/envs"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/jsonl"
	"go.uber.org/zap"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	API_URL_KEY      = "REPORTS_OPENAI_API_URL"
	MODEL_KEY        = "REPORTS_OPENAI_API_MODEL"
	API_KEY          = "REPORTS_OPENAI_API_KEY"
	BATCH_SIZE_KEY   = "REPORTS_OPENAI_BATCH_SIZE_BYTES"
	CONTEXT_SIZE_KEY = "REPORTS_OPENAI_CONTEXT_SIZE_BYTES"
)

type Client struct {
	model            string
	authKey          string
	apiUrl           string
	logger           *zap.Logger
	BatchSizeBytes   int
	ContextSizeBytes int
}

type CompletionRequest struct {
	Model          string     `json:"model"`
	Messages       []*Message `json:"messages"`
	Temperature    float32    `json:"temperature"`
	ResponseFormat any        `json:"response_format"`
}

type CompletionResponse struct {
	Id      string     `json:"id"`
	Object  string     `json:"object"`
	Model   string     `json:"model"`
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

type BatchFileCompletionRequestEntry struct {
	CustomId string             `json:"custom_id"`
	Method   string             `json:"method"`
	Url      string             `json:"url"`
	Body     *CompletionRequest `json:"body"`
}

type BatchFileCompletionResponseEntry struct {
	Id       string `json:"id"`
	CustomId string `json:"custom_id"`
	Response struct {
		StatusCode int                 `json:"status_code"`
		RequestId  string              `json:"request_id"`
		Body       *CompletionResponse `json:"body"`
	} `json:"response"`
}

const (
	COMPLETION_PATH string = "/v1/chat/completions"
	BATCHES_PATH    string = "/v1/batches"
	FILES_API_PATH  string = "/v1/files"
)

func NewOpenAiClient(logger *zap.Logger) *Client {

	envs.ValidateEnvs("OpenAI client parameters are not set", []string{
		API_URL_KEY, MODEL_KEY, API_KEY, BATCH_SIZE_KEY, CONTEXT_SIZE_KEY,
	})

	apiUrl := os.Getenv(API_URL_KEY)
	model := os.Getenv(MODEL_KEY)
	authKey := os.Getenv(API_KEY)
	batchSize := os.Getenv(BATCH_SIZE_KEY)
	contextSize := os.Getenv(CONTEXT_SIZE_KEY)

	batchSizeInt, err := strconv.Atoi(batchSize)
	if err != nil {
		panic("OpenAI batch size is not a number")
	}

	contextSizeInt, err := strconv.Atoi(contextSize)
	if err != nil {
		panic("OpenAI context size is not a number")
	}
	return &Client{
		model:            model,
		authKey:          authKey,
		apiUrl:           apiUrl,
		logger:           logger,
		BatchSizeBytes:   batchSizeInt,
		ContextSizeBytes: contextSizeInt,
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

func (c *Client) batchesUrl() string {
	return fmt.Sprintf("%s/%s", c.apiUrl, BATCHES_PATH)
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

func (c *Client) uploadFileForBatch(inputFile string) (*FileApiReponse, error) {

	httpClient := http.Client{}
	var fileBuf bytes.Buffer
	mpw := multipart.NewWriter(&fileBuf)
	fileWriter, err := mpw.CreateFormFile("file", "batchinput.jsonl")
	if err != nil {
		c.logger.Error("Failed to add batch file", zap.Error(err))
	}

	_, err = fileWriter.Write([]byte(inputFile))
	if err != nil {
		c.logger.Error("Failed to set file header in openai files api request")
		return nil, err
	}

	fieldWriter, err := mpw.CreateFormField("purpose")
	if err != nil {
		c.logger.Error("Failed to set purpose header in openai files api request")
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

type CreateBatchRequest struct {
	InputFileId      string `json:"input_file_id"`
	CompletionWindow string `json:"completion_window"`
	Endpoint         string `json:"endpoint"`
}

type BatchRequestCount struct {
	Total     int `json:"total"`
	Completed int `json:"completed"`
	Failed    int `json:"failed"`
}

type Batch struct {
	Id               string            `json:"id"`
	InputFileId      string            `json:"input_field_id"`
	CompletionWindow string            `json:"completion_window"`
	Status           string            `json:"status"`
	OutputFileId     string            `json:"output_file_id"`
	CreatedAt        int64             `json:"created_at"`
	ExpiresAt        int64             `json:"expires_at"`
	CompletedAt      int64             `json:"completed_at"`
	FailedAt         int64             `json:"failed_at"`
	ExpiredAt        int64             `json:"expired_at"`
	RequestCounts    BatchRequestCount `json:"request_counts"`
}

func (c *Client) createBatch(batchParams CreateBatchRequest) (*Batch, error) {

	httpClient := http.Client{}

	body, err := json.Marshal(batchParams)
	if err != nil {
		c.logger.Error("Failed to encode batch params", zap.Error(err))
		return nil, err
	}

	req, err := http.NewRequest(
		"POST",
		c.batchesUrl(),
		strings.NewReader(string(body)),
	)

	if err != nil {
		c.logger.Error("Failed to create batch request", zap.Error(err))
		return nil, err
	}

	c.addAuthHeaders(req)
	c.addEncodingHeaders(req)

	res, err := httpClient.Do(req)
	if err != nil {
		c.logger.Error("Failed to make batch request", zap.Error(err))
		return nil, err
	}

	defer res.Body.Close()

	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		c.logger.Error("Failed to read batch response", zap.Error(err))
		return nil, err
	}

	var decodedResponse Batch
	err = json.Unmarshal(respBody, &decodedResponse)
	if err != nil {
		c.logger.Error("Failed to decode batch response", zap.Error(err))
		return nil, err
	}

	return &decodedResponse, nil
}

func (c *Client) UploadAndCreateBatches(completionRequests []*CompletionRequest) ([]*Batch, error) {

	// Split requests to batches of size not greater than OpenAi's maximum batch size (~100MB)
	requestsByBatch, err := c.splitCompletionReqestsByBatchSize(completionRequests)
	if err != nil {
		c.logger.Error("Failed to split completion requests by batch", zap.Error(err))
		return nil, err
	}

	return parallelRequest(requestsByBatch, func(requests []*CompletionRequest) (*Batch, error) {

		batch, err := c.UploadAndCreateBatch(requests)
		if err != nil {
			c.logger.Error("Failed to create and upload a completion request batch", zap.Error(err))
			return nil, err
		}
		return batch, nil
	})
}

func (c *Client) splitCompletionReqestsByBatchSize(completionRequests []*CompletionRequest) ([][]*CompletionRequest, error) {

	var requestPackets [][]*CompletionRequest
	var lastPacket []*CompletionRequest
	var lastPacketSize = 0

	for _, request := range completionRequests {

		encodedPacket, err := json.Marshal(request)

		if err != nil {
			c.logger.Error("Failed to enode completion request packet", zap.Error(err))
			return nil, err
		}

		if lastPacketSize+len(encodedPacket) > c.BatchSizeBytes {
			requestPackets = append(requestPackets, lastPacket)
			lastPacket = []*CompletionRequest{request}
			lastPacketSize = len(encodedPacket)

		} else {
			lastPacket = append(lastPacket, request)
			lastPacketSize += len(encodedPacket)
		}

	}
	requestPackets = append(requestPackets, lastPacket)
	return requestPackets, nil
}

func (c *Client) UploadAndCreateBatch(completionRequests []*CompletionRequest) (*Batch, error) {

	batchRequestEntries := array.Map(NewBatchEntryFromCompletionRequet)(completionRequests)

	batchFile := bytes.NewBufferString("")
	err := jsonl.NewJsonLinesEncoder(batchFile).Encode(batchRequestEntries)

	if err != nil {
		c.logger.Sugar().Errorln("Failed to encode messages to complete", zap.Error(err))
		return nil, err
	}

	batchUploadResponse, err := c.uploadFileForBatch(batchFile.String())
	if err != nil {
		c.logger.Error("Failed to upload files for a batch", zap.Error(err))
		return nil, err
	}

	batchCreateResponse, err := c.createBatch(CreateBatchRequest{
		InputFileId: batchUploadResponse.Id,

		// For now this is the only possible value for completion window xd.
		CompletionWindow: "24h",
		Endpoint:         COMPLETION_PATH,
	})
	if err != nil {
		c.logger.Error("Failed to execute batch", zap.Error(err))
		return nil, err
	}

	return batchCreateResponse, nil
}

func NewBatchEntryFromCompletionRequet(completionRequest *CompletionRequest) *BatchFileCompletionRequestEntry {
	return &BatchFileCompletionRequestEntry{
		CustomId: fmt.Sprintf("report-%s", time.Now().String()),
		Method:   "POST",
		Url:      COMPLETION_PATH,
		Body:     completionRequest,
	}
}

func (c *Client) Batches(ids []string) ([]*Batch, error) {

	return parallelRequest(ids, func(id string) (*Batch, error) {
		batch, err := c.Batch(id)
		if err != nil {
			c.logger.Error("Failed to retrieve batch", zap.Error(err))
			return nil, err
		}
		return batch, nil
	})

}

func parallelRequest[Key any, Value any](keys []Key, f func(Key) (Value, error)) ([]Value, error) {

	var valueChannel = make(chan Value, len(keys))
	var values = make([]Value, 0, len(keys))

	var wg sync.WaitGroup
	for _, key := range keys {
		wg.Add(1)

		go func() {
			defer wg.Done()

			value, err := f(key)
			if err != nil {
				return
			}

			valueChannel <- value
		}()
	}

	wg.Wait()
	close(valueChannel)

	for batch := range valueChannel {
		values = append(values, batch)
	}

	return values, nil
}

func (c *Client) Batch(id string) (*Batch, error) {

	httpClient := http.Client{}

	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/%s", c.batchesUrl(), id),
		strings.NewReader(""),
	)

	if err != nil {
		c.logger.Error("Failed to create batch request", zap.Error(err))
		return nil, err
	}

	c.addAuthHeaders(req)
	c.addEncodingHeaders(req)

	res, err := httpClient.Do(req)
	if err != nil {
		c.logger.Error("Failed to make batch request", zap.Error(err))
		return nil, err
	}

	defer res.Body.Close()

	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		c.logger.Error("Failed to read batch response", zap.Error(err))
		return nil, err
	}

	var decodedResponse Batch
	err = json.Unmarshal(respBody, &decodedResponse)
	if err != nil {
		c.logger.Error("Failed to decode batch response", zap.Error(err))
		return nil, err
	}

	return &decodedResponse, nil
}

func (c *Client) CompletionResponseEntriesFromBatch(batch *Batch) ([]*BatchFileCompletionResponseEntry, error) {

	outputFile, err := c.File(batch.OutputFileId)
	if err != nil {
		c.logger.Error("Failed to get batch output file", zap.Error(err))
		return nil, err
	}

	var responses []*BatchFileCompletionResponseEntry
	err = jsonl.NewJsonLinesDecoder(bytes.NewReader(outputFile)).Decode(&responses)
	if err != nil {
		c.logger.Error("Failed to decode completion batch response ", zap.Error(err))
		return nil, err
	}

	return responses, nil
}
func (c *Client) CompletionResponseEntriesFromBatches(batches []*Batch) ([]*BatchFileCompletionResponseEntry, error) {
	allResponses, err := parallelRequest(batches, func(batch *Batch) ([]*BatchFileCompletionResponseEntry, error) {
		responseEntries, err := c.CompletionResponseEntriesFromBatch(batch)
		if err != nil {
			c.logger.Error("Failed to get response entries from completion batch")
			return nil, err
		}
		return responseEntries, nil
	})

	if err != nil {
		c.logger.Error("Failed to get completions requests form batches")
		return nil, err
	}

	var flattenedResponses []*BatchFileCompletionResponseEntry
	for _, batchResponses := range allResponses {
		flattenedResponses = append(flattenedResponses, batchResponses...)
	}

	return flattenedResponses, nil
}

func (c *Client) BatchOutputFiles(batches []*Batch) ([][]byte, error) {

	return parallelRequest(batches, func(batch *Batch) ([]byte, error) {
		file, err := c.File(batch.OutputFileId)
		if err != nil {
			c.logger.Error("Failed to retrieve batch output file", zap.Error(err))
			return nil, err
		}
		return file, nil
	})

}

func (c *Client) Files(ids []string) ([][]byte, error) {

	return parallelRequest(ids, func(id string) ([]byte, error) {
		file, err := c.File(id)
		if err != nil {
			c.logger.Error("Failed to retrieve file", zap.Error(err))
			return nil, err
		}
		return file, nil
	})

}

func (c *Client) File(id string) ([]byte, error) {

	httpClient := http.Client{}

	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/%s/content", c.filesApiUrl(), id),
		strings.NewReader(""),
	)
	if err != nil {
		c.logger.Error("Failed to create batch request", zap.Error(err))
		return nil, err
	}

	c.addAuthHeaders(req)

	res, err := httpClient.Do(req)
	if err != nil {
		c.logger.Error("Failed to make batch request", zap.Error(err))
		return nil, err
	}

	defer res.Body.Close()

	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		c.logger.Error("Failed to read batch response", zap.Error(err))
		return nil, err
	}

	return respBody, nil
}
