package database

import (
	"context"
	"time"
)

type Client interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	HSet(ctx context.Context, key, field string, value interface{}, expiration time.Duration) error
	HGetAll(ctx context.Context, key string) (map[string]string, error)
}
