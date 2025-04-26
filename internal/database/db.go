package database

import (
	"context"
	"time"
)

type Client interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string, source interface{}) error
	HSet(ctx context.Context, key, field string, value interface{}, expiration time.Duration) error
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	Exists(ctx context.Context, key string) (bool, error)
	DelByKey(ctx context.Context, key string) error
}
