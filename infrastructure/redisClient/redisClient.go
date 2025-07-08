package redisClient

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func NewClient() (*redis.Client, error) {
	options := redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		Protocol: 2,
	}

	client := redis.NewClient(&options)

	err := client.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}

	return client, err
}
