package auth

import "context"

type SessionRepository interface {
	Set(ctx context.Context, key string, value interface{}, ttl string) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}
