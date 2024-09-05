package infrastructure

import (
	"time"
	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

type TokenStorage interface {
	SaveToken(token string, expiration time.Duration) error
	GetToken(token string) (string, error)
}

type RedisTokenStorage struct {
	client *redis.Client
}

func NewRedisTokenStorage(client *redis.Client) TokenStorage {
	return &RedisTokenStorage{client: client}
}

func (r *RedisTokenStorage) SaveToken(token string, expiration time.Duration) error {
	return r.client.Set(context.Background(), token, token, expiration).Err()
}

func (r *RedisTokenStorage) GetToken(token string) (string, error) {
	val, err := r.client.Get(context.Background(), token).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}
