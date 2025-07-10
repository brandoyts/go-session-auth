package user

import "context"

type DBRepository interface {
	FindOne(ctx context.Context, user User) (*User, error)
	FindById(ctx context.Context, id string) (*User, error)
	Create(ctx context.Context, user User) (string, error)
}
