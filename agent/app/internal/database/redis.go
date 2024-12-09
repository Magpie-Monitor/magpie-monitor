package database

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisDatabase interface {
	Set(key, value string, ttl int) error
	Get(key string) string
}

type Redis struct {
	url      string
	password string
	db       int
	client   *redis.Client
}

func NewRedis(url, password string, db int) RedisDatabase {
	client := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: password,
		DB:       db,
	})

	_, err := client.Ping(context.TODO()).Result()
	if err != nil {
		log.Println(err)
		panic(err)
	}

	return &Redis{url: url, password: password, db: db, client: client}
}

func (r *Redis) Set(key, value string, ttl int) error {
	return r.client.Set(context.Background(), key, value, time.Duration(ttl)).Err()
}

func (r *Redis) Get(key string) string {
	return r.client.Get(context.Background(), key).Val()
}

var db RedisDatabase = &Redis{}
