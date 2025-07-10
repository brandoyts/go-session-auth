package redisClient

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func NewClient(options *redis.Options) (*redis.Client, error) {
	client := redis.NewClient(options)

	err := client.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}

	return client, err
}
