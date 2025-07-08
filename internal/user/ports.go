package user

import "context"

type DBRepository interface {
	// All(ctx context.Context) ([]User, error)
	// Find(ctx context.Context, user User) ([]User, error)
	FindOne(ctx context.Context, user User) (*User, error)
	FindById(ctx context.Context, id string) (*User, error)
	Create(ctx context.Context, user User) (string, error)
	// Update(ctx context.Context, id string, user User) error
	// Delete(ctx context.Context, id string) error
}
