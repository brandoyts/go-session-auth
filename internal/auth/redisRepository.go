package auth

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) *RedisRepository {
	return &RedisRepository{client: client}
}

func (rr *RedisRepository) Set(ctx context.Context, key string, value interface{}, ttl string) error {
	parsedDuration, err := time.ParseDuration(ttl)
	if err != nil {
		return err
	}

	return rr.client.Set(ctx, key, value, parsedDuration).Err()
}

func (rr *RedisRepository) Get(ctx context.Context, key string) (string, error) {
	result := rr.client.Get(ctx, key)
	if result.Err() != nil {
		return "", result.Err()
	}

	return result.Val(), nil
}

func (rr *RedisRepository) Delete(ctx context.Context, key string) error {
	return rr.client.Del(ctx, key).Err()
}
