package openai

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/envs"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/jsonl"
	scheduledjobs "github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/scheduled_jobs"
	"go.uber.org/zap"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	OpenAiBatchStatus__Validating = "validating"
	OpenAiBatchStatus__Failed     = "failed"
	OpenAiBatchStatus__InProgress = "in_progress"
	OpenAiBatchStatus__Finalizing = "finalizing"
	OpenAiBatchStatus__Completed  = "completed"
	OpenAiBatchStatus__Expired    = "expired"
	OpenAiBatchStatus__Cancelling = "cancelling"
	OpenAiBatchStatus__Cancelled  = "cancelled"
)
const (
	CHARS_PER_OPENAI_TOKEN = 5
)

const (
	BATCH_AWAITING_INTERVAL_SECONDS_KEY     = "REPORTS_BATCH_AWAITING_INTERVAL_SECONDS"
	REPORTS_MAX_IN_PRORESS_TOKENS_KEY       = "REPORTS_MAX_IN_PRORESS_TOKENS"
	MAX_OPENAI_OUTPUT_COMPLETION_TOKENS_KEY = "REPORTS_MAX_OPENAI_OUTPUT_COMPLETION_TOKENS"
)

type BatchPoller struct {
	batches                      chan *Batch
	client                       *Client
	scheduledJobsRepository      scheduledjobs.ScheduledJobRepository[*OpenAiJob]
	pollingIntervalSeconds       int
	batchAwaitingIntervalSeconds int
	maxInProgressTokens          int
	maxCompletionOutputTokens    int
}

func NewBatchPoller(client *Client, scheduledJobsRepository scheduledjobs.ScheduledJobRepository[*OpenAiJob]) *BatchPoller {

	envs.ValidateEnvs("Missing envs for openai batch poller",
		[]string{POLLING_INTERVAL_SECONDS_KEY,
			BATCH_AWAITING_INTERVAL_SECONDS_KEY,
			REPORTS_MAX_IN_PRORESS_TOKENS_KEY,
			MAX_OPENAI_OUTPUT_COMPLETION_TOKENS_KEY,
		})

	pollingIntervalSeconds := os.Getenv(POLLING_INTERVAL_SECONDS_KEY)
	pollingIntervalSecondsInt, err := strconv.Atoi(pollingIntervalSeconds)
	if err != nil {
		panic(fmt.Sprintf("%s is not a number", POLLING_INTERVAL_SECONDS_KEY))
	}

	batchAwaitingIntervalSeconds := os.Getenv(BATCH_AWAITING_INTERVAL_SECONDS_KEY)
	batchAwaitingIntervalSecondsInt, err := strconv.Atoi(batchAwaitingIntervalSeconds)
	if err != nil {
		panic(fmt.Sprintf("%s is not a number", BATCH_AWAITING_INTERVAL_SECONDS_KEY))
	}

	maxOutputTokens := os.Getenv(MAX_OPENAI_OUTPUT_COMPLETION_TOKENS_KEY)
	maxOutputTokensInt, err := strconv.Atoi(maxOutputTokens)
	if err != nil {
		panic(fmt.Sprintf("%s is not a number", MAX_OPENAI_OUTPUT_COMPLETION_TOKENS_KEY))
	}

	maxInProgressTokens := os.Getenv(REPORTS_MAX_IN_PRORESS_TOKENS_KEY)
	maxInProgressTokensInt, err := strconv.Atoi(maxInProgressTokens)

	if err != nil {
		panic(fmt.Sprintf("%s is not a number", MAX_OPENAI_OUTPUT_COMPLETION_TOKENS_KEY))
	}

	return &BatchPoller{
		batches:                      make(chan *Batch),
		client:                       client,
		scheduledJobsRepository:      scheduledJobsRepository,
		pollingIntervalSeconds:       pollingIntervalSecondsInt,
		batchAwaitingIntervalSeconds: batchAwaitingIntervalSecondsInt,
		maxCompletionOutputTokens:    maxOutputTokensInt,
		maxInProgressTokens:          maxInProgressTokensInt,
	}
}

func (p *BatchPoller) tokensFromJobs(jobs []*OpenAiJob) (int64, error) {

	completionTokens := 0

	for _, job := range jobs {
		jobsCompletionTokens, err := p.tokensFromJob(job)
		if err != nil {
			return 0, err
		}

		completionTokens += int(jobsCompletionTokens)
	}

	return int64(completionTokens), nil
}

// Aproximate the tokens for a job based on an average of chars per token and maximum tokens
// that a model might output for each request
func (p *BatchPoller) tokensFromJob(job *OpenAiJob) (int64, error) {

	completionTokens := p.maxCompletionOutputTokens * len(job.CompletionRequests)

	batchFile := bytes.NewBufferString("")
	err := jsonl.NewJsonLinesEncoder(batchFile).Encode(job.CompletionRequests)
	if err != nil {
		p.client.logger.Error("Failed to encode scheduled job", zap.Error(err), zap.Any("job", job.Id))
		return 0, err
	}

	completionTokens += batchFile.Len() / CHARS_PER_OPENAI_TOKEN

	return int64(completionTokens), nil
}

func (p *BatchPoller) dequeScheduledJob(enqueuedJobs []*OpenAiJob, pendingJobs []*OpenAiJob) error {
	if len(enqueuedJobs) == 0 {
		return nil
	}

	inProgressTokens, err := p.tokensFromJobs(pendingJobs)
	if err != nil {
		p.client.logger.Error("Failed to calculate tokens from pending jobs", zap.Error(err), zap.Any("pendingJobs", pendingJobs))
		return err
	}

	lastEnqueuedJobTokens, err := p.tokensFromJob(enqueuedJobs[0])
	if err != nil {
		p.client.logger.Error("Failed to calculate tokens from enqueued job", zap.Error(err), zap.Any("enqueuedJob", enqueuedJobs[0]))
		return err
	}

	if lastEnqueuedJobTokens+inProgressTokens >= int64(p.maxInProgressTokens) {
		p.client.logger.Info("Waiting for jobs to complete before enqueuing next one", zap.Any("newJob", enqueuedJobs[0]))
		return nil
	}

	batch, err := p.client.UploadAndCreateBatch(enqueuedJobs[0].CompletionRequests)
	if err != nil {
		p.client.logger.Error("Failed to upload enqueued job as a batch", zap.Error(err), zap.Any("enqueuedJob", enqueuedJobs[0]))
		return err
	}

	enqueuedJobs[0].Status = OpenAiJobStatus__InProgress
	enqueuedJobs[0].BatchId = &batch.Id

	updateErr := p.scheduledJobsRepository.UpdateScheduledJob(context.Background(), enqueuedJobs[0])
	if updateErr != nil {
		p.client.logger.Error("Failed to update enqueued job", zap.Error(err), zap.Any("job", enqueuedJobs[0]))
		return err
	}

	return nil
}
func (p *BatchPoller) CompleteScheduledJob(scheduledJob *OpenAiJob) error {
	scheduledJob.Status = OpenAiJobStatus__Completed
	return p.scheduledJobsRepository.UpdateScheduledJob(context.Background(), scheduledJob)
}

func (p *BatchPoller) FailScheduledJob(scheduledJob *OpenAiJob) error {
	scheduledJob.Status = OpenAiJobStatus__Failed
	return p.scheduledJobsRepository.UpdateScheduledJob(context.Background(), scheduledJob)
}

func (p *BatchPoller) Start() {

	for {
		enqueuedJobs, err := p.scheduledJobsRepository.GetScheduledJobsByStatus(context.Background(), OpenAiJobStatus__Enqueued)
		pendingJobs, err := p.scheduledJobsRepository.GetScheduledJobsByStatus(context.Background(), OpenAiJobStatus__InProgress)

		p.client.logger.Debug("Enqueued jobs", zap.Any("jobs", len(enqueuedJobs)))
		p.client.logger.Debug("Pending jobs", zap.Any("jobs", len(pendingJobs)))

		scheduleErr := p.dequeScheduledJob(enqueuedJobs, pendingJobs)
		if scheduleErr != nil {
			p.client.logger.Error("Failed to dequeue scheduled job", zap.Error(err))
			continue
		}

		if err != nil {
			p.client.logger.Error("Failed to get pending batches", zap.Error(err))
			continue
		}

		for _, pendingJob := range pendingJobs {
			batch, err := p.client.Batch(*pendingJob.BatchId)

			if err != nil {
				p.client.logger.Error("Failed to getch batch from OpenAI", zap.Error(err), zap.Any("batch", batch))
				continue
			}

			if batch.isCompleted() {
				p.client.logger.Debug("Batch has been completed", zap.Any("batch", batch))
				p.CompleteScheduledJob(pendingJob)
			}

			if batch.isFailed() {
				p.client.logger.Error("Batch has been failed", zap.Any("batch", batch))
				p.FailScheduledJob(pendingJob)
			}

			if batch.isExpired() {
				p.client.logger.Error("Batch has expired", zap.Any("batch", batch))
				p.FailScheduledJob(pendingJob)
			}

			p.client.logger.Info("Currenly pending batch", zap.Any("batch", pendingJob.BatchId))
		}

		time.Sleep(time.Second * time.Duration(p.pollingIntervalSeconds))
	}
}

// func (p *BatchPoller) scheduledJob(jobId string) (*OpenAiJob, error) {
//
// 	job, err := p.scheduledJobsRepository.GetScheduledJob(context.Background(), jobId)
//
// 	// if (*job).IsEqueued() {
// 	// 	return job
// 	// }
//
// 	// batchId := (**job).BatchId
// 	if err != nil {
// 		// p.client.logger.Error("Failed to check if batch is pending", zap.String("batchId", *batchId), zap.Error(err))
// 		p.client.logger.Error("Failed to get scheduled job by id", )
// 		return nil, err
// 	}
// 	//
// 	// batch, clientErr := p.client.Batch(*batchId)
// 	//
// 	// if clientErr != nil {
// 	// 	p.client.logger.Error("Failed to fetch completed batch", zap.Error(err), zap.Any("batch", batchId))
// 	// 	return nil, err
// 	// }
//
// 	return *job, nil
// }

// func (p *BatchPoller) ManyBatches(batchIds []string) (map[string]*Batch, error) {
//
// 	batches := make(map[string]*Batch, 0)
//
// 	for _, batchId := range batchIds {
// 		batch, err := p.batch(batchId)
// 		if err != nil {
// 			p.client.logger.Error("Failed to get batch from poller", zap.Error(err), zap.Any("batch", batchId))
// 			return nil, err
// 		}
// 		batches[batchId] = batch
// 	}
//
// 	return batches, nil
//
// }

func (p *BatchPoller) BatchesFromJobs(jobs []*OpenAiJob) ([]*Batch, error) {

	batches := make([]*Batch, 0, len(jobs))

	for _, job := range jobs {
		batch, err := p.client.Batch(*job.BatchId)
		if err != nil {
			p.client.logger.Error("Failed to get batch from poller", zap.Error(err), zap.Any("job", job))
			return nil, err
		}
		batches = append(batches, batch)
	}

	return batches, nil
}

// Returns (completed, failed, errors)
func (p *BatchPoller) AwaitPendingJobs(jobIds []string) ([]*OpenAiJob, []*OpenAiJob, error) {
	completedJobsChannel := make(chan *OpenAiJob, len(jobIds))
	failedJobsChannel := make(chan *OpenAiJob, len(jobIds))
	errorsChannel := make(chan error, len(jobIds))
	completedJobs := make([]*OpenAiJob, 0, len(jobIds))
	failedJobs := make([]*OpenAiJob, 0, len(jobIds))

	var wg sync.WaitGroup

	for _, jobId := range jobIds {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				job, err := p.scheduledJobsRepository.GetScheduledJob(context.Background(), jobId)

				p.client.logger.Info("Waiting for job", zap.Any("job", (*job).Id))

				if err != nil {
					p.client.logger.Error("Failed to await an openAi batch", zap.Error(err), zap.Any("jobId", (*job).Id))
					errorsChannel <- err
					return
				}

				if (*job).IsCompleted() {
					p.client.logger.Info("Job was finished", zap.Any("job", (*job).Id))
					completedJobsChannel <- (*job)
					return
				}

				if (*job).IsFailed() {
					p.client.logger.Error("Batch was failed", zap.Any("job", (*job).Id))
					failedJobsChannel <- (*job)
					return
				}

				time.Sleep(time.Second * time.Duration(p.batchAwaitingIntervalSeconds))
			}
		}()
	}

	wg.Wait()

	close(completedJobsChannel)
	close(failedJobsChannel)
	close(errorsChannel)

	for err := range errorsChannel {
		return nil, nil, err
	}

	for job := range completedJobsChannel {
		completedJobs = append(completedJobs, job)
	}

	for batch := range failedJobsChannel {
		failedJobs = append(failedJobs, batch)
	}

	return completedJobs, failedJobs, nil

}

// func (p *BatchPoller) InsertPendingBatch(batch *Batch) error {
// 	if err := p.pendingBatchRepository.AddPendingBatch(batch); err != nil {
// 		p.client.logger.Error("Failed to set pending batch", zap.Error(err), zap.Any("batch", batch))
// 		return err
// 	}
// 	return nil
// }
//
// func (p *BatchPoller) InsertPendingBatches(batches []*Batch) error {
// 	p.client.logger.Info("Batches", zap.Any("batches", batches))
// 	if err := p.pendingBatchRepository.AddPendingBatches(batches); err != nil {
// 		p.client.logger.Error("Failed to set pending batch", zap.Error(err), zap.Any("batches", batches))
// 		return err
// 	}
// 	return nil
// }
