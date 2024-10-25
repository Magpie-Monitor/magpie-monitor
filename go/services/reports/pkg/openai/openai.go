package openai

import (
	"bytes"
	"encoding/json"
	"errors"
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
	API_URL_KEY           = "REPORTS_OPENAI_API_URL"
	MODEL_KEY             = "REPORTS_OPENAI_API_MODEL"
	API_KEY               = "REPORTS_OPENAI_API_KEY"
	BATCH_SIZE_KEY        = "REPORTS_OPENAI_BATCH_SIZE_BYTES"
	CONTEXT_SIZE_KEY      = "REPORTS_OPENAI_CONTEXT_SIZE_BYTES"
	MODEL_TEMPERATURE_KEY = "REPORTS_OPENAI_MODEL_TEMPERATURE"
)

type BatchPoller struct {
	batches                chan *Batch
	client                 *Client
	listeners              []BatchUpdateListener
	pendingBatchRepository PendingBatchsRepository
}

type BatchUpdateListener = func(batch *Batch) error

func NewBatchPoller(client *Client, pendingBatchRepository PendingBatchsRepository) *BatchPoller {

	poller := &BatchPoller{
		batches:                make(chan *Batch),
		client:                 client,
		listeners:              []BatchUpdateListener{},
		pendingBatchRepository: pendingBatchRepository,
	}

	go poller.Start()

	return poller
}

func (p *BatchPoller) Start() {
	for {
		batchIds, err := p.pendingBatchRepository.GetAllPending()
		p.client.logger.Debug("Got batchIds from repository %+v", zap.Any("ids", batchIds))

		if err != nil {
			p.client.logger.Error("Failed to get pending batches", zap.Error(err))
			continue
		}

		for _, batchId := range batchIds {
			batch, err := p.client.Batch(batchId)
			p.client.logger.Debug("Got Batch from OpenAi %+v", zap.Any("batch", batch))
			if err != nil {
				p.client.logger.Error("Failed to getch batch from OpenAI", zap.Error(err))
				continue
			}

			if batch.Status == "completed" {
				p.client.logger.Debug("Batch from OpenAi has been completed %+v", zap.Any("batch", batch))
				p.pendingBatchRepository.CompleteBatch(batchId)
			}
		}

		time.Sleep(100000)
	}
}

func (p *BatchPoller) Batch(batchId string) (*Batch, error) {

	batch, err := p.pendingBatchRepository.GetPendingBatch(batchId)
	if err != nil {
		p.client.logger.Error("Batch is not pending", zap.String("batch", batchId), zap.Error(err))
		// return nil, err
	}

	if batch != nil {
		return batch, nil
	}

	// If the batch is not pending, then fetch it from OpenAi
	batch, err = p.client.Batch(batchId)

	if err != nil {
		p.client.logger.Error("Failed to fetch completed batch from openai", zap.Error(err), zap.Any("batch", batchId))
		// p.client.Batch(batchId)
		return nil, err
	}

	return batch, nil
}

func (p *BatchPoller) ManyBatches(batchIds []string) (map[string]*Batch, error) {

	batches := make(map[string]*Batch, 0)

	for _, batchId := range batchIds {
		batch, err := p.Batch(batchId)
		if err != nil {
			p.client.logger.Error("Failed to get batch from poller", zap.Error(err), zap.Any("batch", batchId))
			return nil, err
		}
		batches[batchId] = batch
	}

	return batches, nil

}

func (p *BatchPoller) InsertPendingBatch(batch *Batch) error {
	if err := p.pendingBatchRepository.AddPendingBatch(batch); err != nil {
		p.client.logger.Error("Failed to set pending batch", zap.Error(err), zap.Any("batch", batch))
		return err
	}
	return nil
}

func (p *BatchPoller) InsertPendingBatches(batches []*Batch) error {
	if err := p.pendingBatchRepository.AddPendingBatches(batches); err != nil {
		p.client.logger.Error("Failed to set pending batch", zap.Error(err), zap.Any("batches", batches))
		return err
	}
	return nil
}

// func (p *BatchPoller) GetBatches(batchIds []string) ([]*Batch, error) {
// 	rawBatches, err := p.pendingBatchRepository.GetMany(batchIds)
// 	if err != nil {
// 		p.client.logger.Error("Failed to fetch pending batches", zap.Error(err), zap.Any("batchIds", batchIds))
// 	}
//
// 	batches := make([]string, 0, len(batchIds))
// 	for _, rawBatch := range rawBatches {
// 		batches = append(batches, rawBatch.(string))
// 	}
//
// 	for _, batch := range batches {
// 		if batch == "error" {
// 			p.client.logger.Error("Failed fetch a batch", zap.Any("batch", batch))
// 			return nil, errors.New(fmt.Sprintf("Failed batch %s", batch))
// 		}
// 		if batch == "pending" {
//
// 		}
// 	}
//
// 	return []*Batch{}, nil
// }

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
	Id               string            `json:"id" redis:"id"`
	InputFileId      string            `json:"input_field_id" redis:"input_field_id"`
	CompletionWindow string            `json:"completion_window" redis:"completion_window"`
	Status           string            `json:"status" redis:"status"`
	OutputFileId     string            `json:"output_file_id" redis:"output_file_id"`
	CreatedAt        int64             `json:"created_at" redis:"created_at"`
	ExpiresAt        int64             `json:"expires_at" redis:"expires_at"`
	CompletedAt      int64             `json:"completed_at" redis:"completed_at"`
	FailedAt         int64             `json:"failed_at" redis:"failed_at"`
	ExpiredAt        int64             `json:"expired_at" redis:"expired_at"`
	RequestCounts    BatchRequestCount `json:"request_counts" redis:"request_counts"`
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

	batches, batchErrors := parallelRequest(requestsByBatch, func(requests []*CompletionRequest) (*Batch, error) {

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

	return batches, nil
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

	if len(err) != 0 {
		c.logger.Error("Failed to get completions requests form batches")
		return nil, errors.Join(err...)
	}

	var flattenedResponses []*BatchFileCompletionResponseEntry
	for _, batchResponses := range allResponses {
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
