package openai

import (
	"fmt"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/envs"
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
	BATCH_AWAITING_INTERVAL_SECONDS_KEY = "REPORTS_BATCH_AWAITING_INTERVAL_SECONDS"
)

type BatchPoller struct {
	batches                      chan *Batch
	client                       *Client
	pendingBatchRepository       PendingBatchRepository
	pollingIntervalSeconds       int
	batchAwaitingIntervalSeconds int
}

func NewBatchPoller(client *Client, pendingBatchRepository PendingBatchRepository) *BatchPoller {

	envs.ValidateEnvs("Missing envs for openai batch poller",
		[]string{POLLING_INTERVAL_SECONDS_KEY,
			BATCH_AWAITING_INTERVAL_SECONDS_KEY})

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

	return &BatchPoller{
		batches:                      make(chan *Batch),
		client:                       client,
		pendingBatchRepository:       pendingBatchRepository,
		pollingIntervalSeconds:       pollingIntervalSecondsInt,
		batchAwaitingIntervalSeconds: batchAwaitingIntervalSecondsInt,
	}
}

func (p *BatchPoller) Start() {
	for {
		batchIds, err := p.pendingBatchRepository.GetAllPending()
		p.client.logger.Debug("Got batchIds from repository", zap.Any("ids", batchIds))

		if err != nil {
			p.client.logger.Error("Failed to get pending batches", zap.Error(err))
			continue
		}

		for _, batchId := range batchIds {
			batch, err := p.client.Batch(batchId)

			if err != nil {
				p.client.logger.Error("Failed to getch batch from OpenAI", zap.Error(err), zap.Any("batch", batch))
				continue
			}

			if batch.isCompleted() {
				p.client.logger.Debug("Batch has been completed", zap.Any("batch", batch))
				p.pendingBatchRepository.CompleteBatch(batchId)
			}

			if batch.isFailed() {
				p.client.logger.Error("Batch has been failed", zap.Any("batch", batch))
				p.pendingBatchRepository.FailBatch(batchId)
			}

			if batch.isExpired() {
				p.client.logger.Error("Batch has expired", zap.Any("batch", batch))
				p.pendingBatchRepository.FailBatch(batchId)
			}
		}

		p.client.logger.Info("Currenly pending batches", zap.Any("batchIds", batchIds))
		time.Sleep(time.Second * time.Duration(p.pollingIntervalSeconds))
	}
}

func (p *BatchPoller) Batch(batchId string) (*Batch, error) {

	batch, err := p.pendingBatchRepository.GetPendingBatch(batchId)
	if err != nil {
		p.client.logger.Error("Failed to check if batch is pending", zap.String("batchId", batchId), zap.Error(err))
		return nil, err
	}

	// If the batch is not in the repository then a nil will be returned
	if batch != nil {
		return batch, nil
	}

	// If the batch is not pending, then fetch it from OpenAi
	batch, err = p.client.Batch(batchId)

	if err != nil {
		p.client.logger.Error("Failed to fetch completed batch", zap.Error(err), zap.Any("batch", batchId))
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

// Returns (completed, failed, errors)
func (p *BatchPoller) AwaitPendingBatches(batchIds []string) ([]*Batch, []*Batch, error) {
	completedBatchesChannel := make(chan *Batch, len(batchIds))
	failedBatchesChannel := make(chan *Batch, len(batchIds))
	errorsChannel := make(chan error, len(batchIds))
	completedBatches := make([]*Batch, 0, len(batchIds))
	failedBatches := make([]*Batch, 0, len(batchIds))

	var wg sync.WaitGroup

	for _, batchId := range batchIds {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				batch, err := p.Batch(batchId)
				p.client.logger.Info("Waiting for batch", zap.Any("batch", batch))
				if err != nil {
					p.client.logger.Error("Failed to await an openAi batch", zap.Error(err), zap.Any("batchId", batchId))
					errorsChannel <- err
					return
				}

				if batch.isCompleted() {
					p.client.logger.Info("Batch was finished", zap.Any("batch", batch))
					completedBatchesChannel <- batch
					return
				}

				if batch.isFailed() {
					p.client.logger.Error("Batch was failed", zap.Any("batch", batch))
					failedBatchesChannel <- batch
					return
				}

				if batch.isExpired() {
					p.client.logger.Error("Batch was expired", zap.Any("batch", batch))
					failedBatchesChannel <- batch
					return
				}

				time.Sleep(time.Second * time.Duration(p.batchAwaitingIntervalSeconds))
			}
		}()
	}

	wg.Wait()

	close(completedBatchesChannel)
	close(failedBatchesChannel)
	close(errorsChannel)

	for err := range errorsChannel {
		return nil, nil, err
	}

	for batch := range completedBatchesChannel {
		completedBatches = append(completedBatches, batch)
	}

	for batch := range failedBatchesChannel {
		failedBatches = append(failedBatches, batch)
	}

	return completedBatches, failedBatches, nil

}

func (p *BatchPoller) InsertPendingBatch(batch *Batch) error {
	if err := p.pendingBatchRepository.AddPendingBatch(batch); err != nil {
		p.client.logger.Error("Failed to set pending batch", zap.Error(err), zap.Any("batch", batch))
		return err
	}
	return nil
}

func (p *BatchPoller) InsertPendingBatches(batches []*Batch) error {
	p.client.logger.Info("Batches", zap.Any("batches", batches))
	if err := p.pendingBatchRepository.AddPendingBatches(batches); err != nil {
		p.client.logger.Error("Failed to set pending batch", zap.Error(err), zap.Any("batches", batches))
		return err
	}
	return nil
}
