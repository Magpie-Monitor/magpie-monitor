package openai

import (
	"errors"
	"fmt"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/envs"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/redis"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"os"
	"strconv"
	"strings"
)

const (
	REDIS_URL_KEY      = "REPORTS_PENDING_BATCH_REDIS_URL"
	REDIS_PASSWORD_KEY = "REPORTS_PENDING_BATCH_REDIS_PASSWORD"
	REDIS_DB_KEY       = "REPORTS_PENDING_BATCH_REDIS_DB"
)

type PendingBatchsRepository interface {
	AddPendingBatch(batch *Batch) error
	AddPendingBatches(batch []*Batch) error
	CompleteBatch(batchId string) error
	FailBatch(batchId string) error
	GetAllPending() ([]string, error)
	GetPendingBatch(id string) (*Batch, error)
}

type RedisPendingBatchRepository struct {
	redisClient redis.Redis
	logger      *zap.Logger
}

type RedisPendingBatchRepositoryParams struct {
	fx.In
	Logger *zap.Logger
}

const (
	RedisPendingBatchNamespace__Pending = "pending"
	RedisPendingBatchNamespace__Failed  = "failed"
)

func NewRedisPendingBatchRepository(params RedisPendingBatchRepositoryParams) *RedisPendingBatchRepository {

	envs.ValidateEnvs("Failed to build RedisPendingBatchRepository. Missing env",
		[]string{REDIS_URL_KEY, REDIS_PASSWORD_KEY, REDIS_DB_KEY})

	redisUrl := os.Getenv(REDIS_URL_KEY)
	redisPassword := os.Getenv(REDIS_PASSWORD_KEY)
	redisDb := os.Getenv(REDIS_DB_KEY)

	redisDbInt, err := strconv.Atoi(redisDb)
	if err != nil {
		panic("PendingBatchRepository redisDb is not a number")
	}

	redisClient := redis.NewRedis(redisUrl, redisPassword, redisDbInt)

	return &RedisPendingBatchRepository{
		logger:      params.Logger,
		redisClient: redisClient,
	}
}

func (r *RedisPendingBatchRepository) AddPendingBatch(batch *Batch) error {

	if err := r.redisClient.HSet(r.getPendingKeyFromBatchId(batch.Id), batch); err != nil {
		r.logger.Error("Failed to add pending batch to repository")
		return err
	}

	return nil
}

func (r *RedisPendingBatchRepository) AddPendingBatches(batches []*Batch) error {

	for _, batch := range batches {
		if err := r.redisClient.HSet(r.getPendingKeyFromBatchId(batch.Id), batch); err != nil {
			r.logger.Error("Failed to add pending batch to repository")
			return err
		}
	}

	return nil
}

func (r *RedisPendingBatchRepository) FailBatch(batchId string) error {
	return r.CompleteBatch(batchId)
}

func (r *RedisPendingBatchRepository) CompleteBatch(batchId string) error {

	if err := r.redisClient.Del(fmt.Sprintf("in_progress:%s", batchId)); err != nil {
		r.logger.Error("Failed to remove pending batch to repository")
		return err
	}

	return nil
}

func (r *RedisPendingBatchRepository) GetPendingBatch(batchId string) (*Batch, error) {

	var resultBatch Batch
	if err := r.redisClient.HGetAll(r.getPendingKeyFromBatchId(batchId), &resultBatch); err != nil {
		r.logger.Error("Failed to get a pending batch")
		return nil, err
	}

	// RedisClient leaves object unpopulated if no entry is found
	if resultBatch.Id == "" {
		return nil, nil
	}

	return &resultBatch, nil
}

func (r *RedisPendingBatchRepository) getBatchIdFromKey(key string) (string, error) {
	parts := strings.Split(key, ":")
	if len(parts) < 2 {
		return "", errors.New("Incorrect key (missing key type or separator :)")
	}

	return parts[1], nil
}

func (r *RedisPendingBatchRepository) getPendingKeyFromBatchId(batchId string) string {
	return fmt.Sprintf("%s:%s", RedisPendingBatchNamespace__Pending, batchId)
}

func (r *RedisPendingBatchRepository) GetAllPending() ([]string, error) {

	keys, err := r.redisClient.Keys(fmt.Sprintf("%s:*", RedisPendingBatchNamespace__Pending))
	if err != nil {
		r.logger.Error("Failed to remove pending batch to repository")
		return nil, err
	}

	batchIds := make([]string, 0, len(keys))

	for _, key := range keys {
		batchId, err := r.getBatchIdFromKey(key)
		if err != nil {
			r.logger.Error("Failed to get pending batch key", zap.Error(err), zap.Any("batchId", batchId))
			return nil, err
		}

		batchIds = append(batchIds, batchId)
	}

	return batchIds, nil
}

func ProvideAsPendingBatchRepository(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(PendingBatchsRepository)),
	)
}

var _ PendingBatchsRepository = &RedisPendingBatchRepository{}
