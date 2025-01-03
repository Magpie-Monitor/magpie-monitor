package openai

import (
	"bytes"
	"encoding"
	"encoding/json"
	"errors"
	"fmt"
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
)

const (
	API_URL_KEY                  = "REPORTS_OPENAI_API_URL"
	MODEL_KEY                    = "REPORTS_OPENAI_API_MODEL"
	API_KEY                      = "REPORTS_OPENAI_API_KEY"
	BATCH_SIZE_KEY               = "REPORTS_OPENAI_BATCH_SIZE_BYTES"
	CONTEXT_SIZE_KEY             = "REPORTS_OPENAI_CONTEXT_SIZE_BYTES"
	MODEL_TEMPERATURE_KEY        = "REPORTS_OPENAI_MODEL_TEMPERATURE"
	POLLING_INTERVAL_SECONDS_KEY = "REPORTS_OPENAI_POLLING_INTERVAL_SECONDS"
)

type Client struct {
	model            string
	authKey          string
	apiUrl           string
	logger           *zap.Logger
	BatchSizeBytes   int
	ContextSizeBytes int
	Temperature      float32
}

type CompletionRequest struct {
	Model          string     `json:"model" bson:"model"`
	Messages       []*Message `json:"messages" bson:"messages"`
	Temperature    float32    `json:"temperature" bson:"temperature"`
	ResponseFormat any        `json:"response_format" bson:"response_format"`
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
	Role    string `json:"role" bson:"role"`
	Content string `json:"content" bson:"content"`
}

type BatchFileCompletionRequestEntry struct {
	CustomId string             `json:"custom_id" bson:"custom_id"`
	Method   string             `json:"method" bson:"method"`
	Url      string             `json:"url" bson:"url"`
	Body     *CompletionRequest `json:"body" bson:"body"`
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
		API_URL_KEY, MODEL_KEY, API_KEY, BATCH_SIZE_KEY, CONTEXT_SIZE_KEY, MODEL_TEMPERATURE_KEY,
	})

	apiUrl := os.Getenv(API_URL_KEY)
	model := os.Getenv(MODEL_KEY)
	authKey := os.Getenv(API_KEY)
	batchSize := os.Getenv(BATCH_SIZE_KEY)
	contextSize := os.Getenv(CONTEXT_SIZE_KEY)
	temperature := os.Getenv(MODEL_TEMPERATURE_KEY)

	batchSizeInt, err := strconv.Atoi(batchSize)
	if err != nil {
		panic("OpenAI batch size is not a number")
	}

	contextSizeInt, err := strconv.Atoi(contextSize)
	if err != nil {
		panic("OpenAI context size is not a number")
	}

	temperatureFloat, err := strconv.ParseFloat(temperature, 2)
	if err != nil {
		panic("OpenAI model temperature is not a float")
	}
	return &Client{
		model:            model,
		authKey:          authKey,
		apiUrl:           apiUrl,
		logger:           logger,
		BatchSizeBytes:   batchSizeInt,
		ContextSizeBytes: contextSizeInt,
		Temperature:      float32(temperatureFloat),
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
		Temperature:    c.Temperature,
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
	Id               string             `json:"id" redis:"id"`
	InputFileId      string             `json:"input_field_id" redis:"input_field_id"`
	CompletionWindow string             `json:"completion_window" redis:"completion_window"`
	Status           string             `json:"status" redis:"status"`
	Errors           BatchErrorsWrapper `json:"-" redis:"-"`
	OutputFileId     string             `json:"output_file_id" redis:"output_file_id"`
	CreatedAt        int64              `json:"created_at" redis:"created_at"`
	ExpiresAt        int64              `json:"expires_at" redis:"expires_at"`
	CompletedAt      int64              `json:"completed_at" redis:"completed_at"`
	FailedAt         int64              `json:"failed_at" redis:"failed_at"`
	ExpiredAt        int64              `json:"expired_at" redis:"expired_at"`
	RequestCounts    BatchRequestCount  `json:"-" redis:"-"`
}

type BatchErrorsWrapper struct {
	Object string       `json:"object" redis:"object"`
	Data   []BatchError `json:"data" redis:"data"`
}

type BatchError struct {
	Code    string `json:"code" redis:"code"`
	Message string `json:"message" redis:"message"`
}

func (b *Batch) isCompleted() bool {
	return b.Status == OpenAiBatchStatus__Completed
}

func (b *Batch) isFailed() bool {
	return b.Status == OpenAiBatchStatus__Failed
}

func (b *Batch) isExpired() bool {
	return b.Status == OpenAiBatchStatus__Expired
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
		c.logger.Error("Failed to decode batch response", zap.Error(err), zap.Any("body", respBody))
		return nil, err
	}

	if decodedResponse.Id == "" {
		c.logger.Error("Failed to recieve id from a batch", zap.Any("response", respBody))
		return nil, errors.New("Failed to recieve id from a batch")
	}

	return &decodedResponse, nil
}

func (c *Client) UploadAndCreateBatches(completionRequests map[string]*CompletionRequest) ([]*Batch, error) {

	// Split requests to batches of size not greater than OpenAi's maximum batch size (~100MB)
	requestsByBatch, err := c.SplitCompletionReqestsByBatchSize(completionRequests)
	if err != nil {
		c.logger.Error("Failed to split completion requests by batch", zap.Error(err))
		return nil, err
	}

	batches, batchErrors := parallelRequest(requestsByBatch, func(requests map[string]*CompletionRequest) (*Batch, error) {

		batch, err := c.UploadAndCreateBatch(requests)
		if err != nil {
			c.logger.Error("Failed to create and upload a completion request batch", zap.Error(err))
			return nil, err
		}
		return batch, nil
	})

	if len(batchErrors) != 0 {
		c.logger.Error("Failed to upload batches", zap.Error(errors.Join(batchErrors...)))
		return nil, errors.Join(batchErrors...)
	}

	c.logger.Debug("Creating batches", zap.Any("batches", batches))

	return batches, nil
}

func (c *Client) SplitCompletionReqestsByBatchSize(completionRequests map[string]*CompletionRequest) ([]map[string]*CompletionRequest, error) {

	requestPackets := make([]map[string]*CompletionRequest, 0)
	lastPacket := make(map[string]*CompletionRequest, 0)
	var lastPacketSize = 0

	if len(completionRequests) == 0 {
		return requestPackets, nil
	}

	for customId, request := range completionRequests {

		encodedPacket, err := json.Marshal(request)

		if err != nil {
			c.logger.Error("Failed to enode completion request packet", zap.Error(err))
			return nil, err
		}

		if lastPacketSize+len(encodedPacket) > c.BatchSizeBytes {
			requestPackets = append(requestPackets, lastPacket)
			lastPacket = map[string]*CompletionRequest{customId: request}
			lastPacketSize = len(encodedPacket)

		} else {
			lastPacket[customId] = request
			lastPacketSize += len(encodedPacket)
		}

	}

	requestPackets = append(requestPackets, lastPacket)

	return requestPackets, nil
}

func (c *Client) UploadAndCreateBatch(completionRequests map[string]*CompletionRequest) (*Batch, error) {

	batchCompletionEntries := make([]*BatchFileCompletionRequestEntry, 0, len(completionRequests))
	for customId, request := range completionRequests {
		batchCompletionEntries = append(batchCompletionEntries, NewBatchEntryFromCompletionRequest(request, customId))
	}

	batchFile := bytes.NewBufferString("")
	err := jsonl.NewJsonLinesEncoder(batchFile).Encode(batchCompletionEntries)

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
	c.logger.Info("batch response", zap.Any("batch", batchCreateResponse))

	return batchCreateResponse, nil
}

func NewBatchEntryFromCompletionRequest(completionRequest *CompletionRequest, customId string) *BatchFileCompletionRequestEntry {
	return &BatchFileCompletionRequestEntry{
		CustomId: customId,
		Method:   "POST",
		Url:      COMPLETION_PATH,
		Body:     completionRequest,
	}
}

func (c *Client) Batches(ids []string) ([]*Batch, error) {

	batches, err := parallelRequest(ids, func(id string) (*Batch, error) {
		batch, err := c.Batch(id)
		if err != nil {
			c.logger.Error("Failed to retrieve batch", zap.Error(err))
			return nil, err
		}
		return batch, nil
	})

	if len(err) != 0 {
		c.logger.Error("Failed to get OpenAI batches by ids", zap.Error(errors.Join(err...)))
		return nil, errors.Join(err...)
	}

	return batches, nil

}

func parallelRequest[Key any, Value any](keys []Key, f func(Key) (Value, error)) ([]Value, []error) {

	var valueChannel = make(chan Value, len(keys))
	var values = make([]Value, 0, len(keys))

	var errorChannel = make(chan error, len(keys))
	var errors = make([]error, 0, len(keys))

	var wg sync.WaitGroup
	for _, key := range keys {
		wg.Add(1)

		go func() {
			defer wg.Done()

			value, err := f(key)
			if err != nil {
				errorChannel <- err
				return
			}

			valueChannel <- value
		}()
	}

	wg.Wait()
	close(valueChannel)
	close(errorChannel)

	for value := range valueChannel {
		values = append(values, value)
	}

	for err := range errorChannel {
		errors = append(errors, err)
	}
	if len(errors) != 0 {
		return nil, errors
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
		c.logger.Error("Failed to decode batch response", zap.Error(err), zap.Any("body", respBody))
		return nil, err
	}

	return &decodedResponse, nil
}

func (c *Client) CompletionResponseEntriesFromBatch(batch *Batch) ([]*BatchFileCompletionResponseEntry, error) {

	responses := make([]*BatchFileCompletionResponseEntry, 0)

	if batch.OutputFileId == "" {
		c.logger.Error("Batch doesn't have an output file, skipping", zap.Any("batch", batch))
		return responses, nil
	}

	outputFile, err := c.File(batch.OutputFileId)
	if err != nil {
		c.logger.Error("Failed to get batch output file", zap.Error(err))
		return nil, err
	}

	err = jsonl.NewJsonLinesDecoder(bytes.NewReader(outputFile)).Decode(&responses)
	if err != nil {
		c.logger.Error("Failed to decode completion batch response ", zap.Error(err))
		return nil, err
	}

	return responses, nil
}

func (c *Client) CompletionResponseEntriesFromBatchById(batchId string) ([]*BatchFileCompletionResponseEntry, error) {
	batch, err := c.Batch(batchId)
	if err != nil {
		c.logger.Error("Failed to get batch by id", zap.Error(err))
		return nil, err
	}

	return c.CompletionResponseEntriesFromBatch(batch)
}

func (c *Client) OutputFileFromBatch(batchId string) ([]byte, error) {
	batch, err := c.Batch(batchId)
	if err != nil {
		c.logger.Error("Failed to fetch a batch by id", zap.Error(err), zap.Any("id", batchId))
		return nil, err
	}

	if batch.OutputFileId == "" {
		c.logger.Error("No output file from batch", zap.Any("batch", batch))
		return nil, errors.New("No output file from batch")
	}

	outputFile, err := c.File(batch.OutputFileId)
	if err != nil {
		c.logger.Error("Failed get output batch file", zap.Error(err), zap.Any("id", batchId))
		return nil, err
	}

	return outputFile, nil
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

	if len(err) != 0 {
		c.logger.Error("Failed to get completions requests form batches")
		return nil, errors.Join(err...)
	}

	var flattenedResponses []*BatchFileCompletionResponseEntry
	for _, batchResponses := range allResponses {

		// Skip batch in case of an error
		if batchResponses == nil {
			continue
		}

		flattenedResponses = append(flattenedResponses, batchResponses...)
	}

	return flattenedResponses, nil
}

func (c *Client) BatchOutputFiles(batches []*Batch) ([][]byte, error) {

	files, err := parallelRequest(batches, func(batch *Batch) ([]byte, error) {
		file, err := c.File(batch.OutputFileId)
		if err != nil {
			c.logger.Error("Failed to retrieve batch output file", zap.Error(err))
			return nil, err
		}
		return file, nil
	})

	if len(err) != 0 {
		c.logger.Error("Failed to get output files from batches", zap.Error(errors.Join(err...)))
		return nil, errors.Join(err...)
	}

	return files, nil
}

func (c *Client) Files(ids []string) ([][]byte, error) {

	files, err := parallelRequest(ids, func(id string) ([]byte, error) {
		file, err := c.File(id)
		if err != nil {
			c.logger.Error("Failed to retrieve file", zap.Error(err))
			return nil, err
		}
		return file, nil
	})

	if len(err) != 0 {
		c.logger.Error("Failed to get files from OpenAI Files API")
		return nil, errors.Join(err...)
	}

	return files, nil

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

// Implement the BinaryMarshaller interface
func (b *Batch) MarshalBinary() ([]byte, error) {
	return json.Marshal(b)
}

var _ encoding.BinaryMarshaler = &Batch{}
