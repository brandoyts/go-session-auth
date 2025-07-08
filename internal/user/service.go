package user

import (
	"context"
)

type Service struct {
	repo DBRepository
}

func NewService(repo DBRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, user User) (string, error) {
	result, err := s.repo.Create(ctx, user)
	if err != nil {
		return "", err
	}

	return result, nil
}

func (s *Service) FindOne(ctx context.Context, user User) (*User, error) {
	result, err := s.repo.FindOne(ctx, user)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Service) FindUserById(ctx context.Context, id string) (*User, error) {
	result, err := s.repo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Service) FindUserByEmail(ctx context.Context, user User) (*User, error) {
	filter := User{
		Email: user.Email,
	}
	result, err := s.repo.FindOne(ctx, filter)
	if err != nil {
		return nil, err
	}

	return result, nil
}
