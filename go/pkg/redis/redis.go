package redis

import (
	"context"
	"encoding/json"
	"log"
	"time"

	redis "github.com/redis/go-redis/v9"
)

type Redis struct {
	url      string
	password string
	db       int
	client   *redis.Client
}

func NewRedis(url, password string, db int) Redis {
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

	return Redis{url: url, password: password, db: db, client: client}
}

func (r *Redis) Set(key, value string, ttl int) error {
	return r.client.Set(context.Background(), key, value, time.Duration(ttl)).Err()
}

func (r *Redis) Get(key string) string {
	// r.client.Keys()
	return r.client.Get(context.Background(), key).Val()
}

func (r *Redis) HKeys(pattern string) ([]string, error) {
	return r.client.HKeys(context.Background(), pattern).Result()
}

//	func (r *Redis) HDel(key string) error {
//		return r.client.HDel()
//	}
func (r *Redis) Del(key string) error {
	cmd := r.client.Del(context.Background(), key)
	return cmd.Err()
}

func (r *Redis) HGet(id string, field string) (string, error) {
	return r.client.HGet(context.Background(), id, field).Result()
}

func (r *Redis) HGetAll(id string, v any) error {
	cmd := r.client.HGetAll(context.Background(), id)
	if cmd.Err() != nil {
		return cmd.Err()
	}

	return cmd.Scan(v)
}

func (r *Redis) HSet(id string, v any) error {

	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	var hSet map[string]interface{}
	err = json.Unmarshal(data, &hSet)
	if err != nil {
		return err
	}

	cmd := r.client.HSet(context.Background(), id, hSet)

	_, err = cmd.Result()
	return err
}
