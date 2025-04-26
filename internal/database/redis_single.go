package database

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"live-chat-kafka/config"
	"live-chat-kafka/internal/models"
	"time"
)

type redisClient struct {
	cfg    config.Redis
	client *redis.Client
}

func NewRedisSingleClient(ctx context.Context, cfg config.Redis) (Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		DialTimeout:  time.Second * 3,
		ReadTimeout:  time.Second * 3,
		WriteTimeout: time.Second * 3,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("fail ping err : %w", err)
	}

	return &redisClient{
		cfg:    cfg,
		client: client,
	}, nil
}

func (r *redisClient) Set(ctx context.Context, key string, data interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, jsonData, expiration).Err()
}

func (r *redisClient) HSet(ctx context.Context, key, field string, data interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return r.client.HSet(ctx, key, field, jsonData, expiration).Err()
}

func (r *redisClient) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	result, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, models.GetCustomErr(models.ErrNotFoundServerInfo)
	}

	return result, nil
}
