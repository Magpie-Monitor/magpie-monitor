package openai

import (
	"github.com/Magpie-Monitor/magpie-monitor/pkg/envs"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/redis"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"os"
	"strconv"
)

const (
	REDIS_URL_KEY      = "REPORTS_PENDING_BATCH_REDIS_URL"
	REDIS_PASSWORD_KEY = "REPORTS_PENDING_BATCH_REDIS_PASSWORD"
	REDIS_DB_KEY       = "REPORTS_PENDING_BATCH_REDIS_DB"
)

type PendingBatchsRepository interface {
	Add(batchId string) error
	Remove(batchId string) error
	GetAll() ([]string, error)
	GetMany(keys []string) ([]interface{}, error)
}

type RedisPendingBatchRepository struct {
	redisClient redis.Redis
	logger      *zap.Logger
}

type RedisPendingBatchRepositoryParams struct {
	fx.In
	Logger *zap.Logger
}

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

func (r *RedisPendingBatchRepository) Add(batchId string) error {

	if err := r.redisClient.Set(batchId, "pending", 1000); err != nil {
		r.logger.Error("Failed to add pending batch to repository")
		return err
	}

	return nil
}

func (r *RedisPendingBatchRepository) Remove(batchId string) error {

	if err := r.redisClient.Set(batchId, "completed", 1000); err != nil {
		r.logger.Error("Failed to remove pending batch to repository")
		return err
	}

	return nil
}

func (r *RedisPendingBatchRepository) GetAll() ([]string, error) {

	res, err := r.redisClient.GetKeys()
	if err != nil {
		r.logger.Error("Failed to remove pending batch to repository")
		return nil, err
	}

	return res, nil
}

func (r *RedisPendingBatchRepository) GetMany(keys []string) ([]interface{}, error) {

	res, err := r.redisClient.MGet(keys)
	r.logger.Debug("Got many entries from pending batch repository %+v", zap.Any("batches", res))
	if err != nil {
		r.logger.Error("Failed to remove pending batch to repository")
		return nil, err
	}

	return res, nil
}

func ProvideAsPendingBatchRepository(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(PendingBatchsRepository)),
	)
}

var _ PendingBatchsRepository = &RedisPendingBatchRepository{}
