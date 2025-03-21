package database

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"live-chat-kafka/config"
	"live-chat-kafka/internal/models"
	"time"
)

const (
	LiveChatServerInfo = "live-chat-server-info"
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

func (r *redisClient) GetAvailableServerList() (map[string]string, error) {
	result, err := r.client.HGetAll(context.TODO(), LiveChatServerInfo).Result()
	if err != nil {
		return nil, models.GetCustomErr(models.ErrNotFoundServerInfo)
	}

	return result, nil
}
