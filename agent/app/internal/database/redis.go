package database

import (
	"context"
	redis "github.com/redis/go-redis/v9"
	"time"
)

type Redis struct {
	url      string
	password string
	db       int
	client   *redis.Client
}

func NewRedis(url, password string, db int) Redis {
	return Redis{url: url, password: password, db: db, client: redis.NewClient(&redis.Options{
		Addr:     url,
		Password: password,
		DB:       db,
	})}
}

func (r *Redis) Set(key, value string, ttl int) error {
	return r.client.Set(context.Background(), key, value, time.Duration(ttl)).Err()
}

func (r *Redis) Get(key string) string {
	return r.client.Get(context.Background(), key).Val()
}
