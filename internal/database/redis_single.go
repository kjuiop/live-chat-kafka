package database

import (
	"context"
	"encoding/json"
	"errors"
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

func (r *redisClient) Get(ctx context.Context, key string, dest interface{}) error {

	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return fmt.Errorf("chat room data is not exist : %s", key)
		}
		return err
	}

	if err := json.Unmarshal([]byte(val), dest); err != nil {
		return fmt.Errorf("unmarshal chat room info err: %w", err)
	}

	return nil
}

func (r *redisClient) HSet(ctx context.Context, key, field string, data interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return r.client.HSet(ctx, key, field, jsonData, expiration).Err()
}

func (r *redisClient) HGet(ctx context.Context, key, mapKey string) (string, error) {

	result, err := r.client.HGet(ctx, key, mapKey).Result()
	if err != nil {
		return "", fmt.Errorf("fail hget data, err : %w", err)
	}

	return result, nil
}

func (r *redisClient) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	result, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, models.GetCustomErr(models.ErrNotFoundServerInfo)
	}

	return result, nil
}

func (r *redisClient) Exists(ctx context.Context, key string) (bool, error) {
	isExist, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return isExist == 1, nil
}

func (r *redisClient) DelByKey(ctx context.Context, key string) error {

	if err := r.client.Del(ctx, key).Err(); err != nil {
		return err
	}

	return nil
}
