package service

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisClient インターフェースを定義
type RedisClient interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
}
