package batchcache

import (
	"context"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/envs"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"os"
)

type BatchCache interface {
	SetReportBatchId(ctx context.Context, reportId string, batchId string) error
	GetBatchId(ctx context.Context, reportId string) (*string, error)
}

type RedisBatchCache struct {
	logger      *zap.Logger
	redisClient *redis.Client
}

const CLOUD_REDIS_PASSWORD_KEY = "CLOUD_REDIS_PASSWORD"
const CLOUD_REDIS_HOST_KEY = "CLOUD_REDIS_HOST"

func NewRedisBatchCache(logger *zap.Logger) *RedisBatchCache {

	envs.ValidateEnvs("RedisBatchCache envs are not set",
		[]string{CLOUD_REDIS_PASSWORD_KEY, CLOUD_REDIS_HOST_KEY})

	host, _ := os.LookupEnv(CLOUD_REDIS_HOST_KEY)
	password, _ := os.LookupEnv(CLOUD_REDIS_PASSWORD_KEY)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       0,
	})

	return &RedisBatchCache{
		logger:      logger,
		redisClient: redisClient,
	}
}

func (c *RedisBatchCache) SetReportBatchId(ctx context.Context, reportId string, batchId string) error {
	err := c.redisClient.Set(ctx, reportId, batchId, 0).Err()
	if err != nil {
		c.logger.Error("Failed to set batchId in reportBatch cache")
		return err
	}

	return nil
}

func (c *RedisBatchCache) GetBatchId(ctx context.Context, reportId string) (*string, error) {
	cmd := c.redisClient.Get(ctx, reportId)

	err := cmd.Err()
	if err != nil {
		c.logger.Error("Failed to get batchId from reportBatch cache")
		return nil, err
	}
	val := cmd.Val()

	return &val, nil
}

var _ BatchCache = &RedisBatchCache{}
