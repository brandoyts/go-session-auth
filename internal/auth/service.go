package auth

import (
	"context"
	"errors"
	"time"

	"github.com/brandoyts/go-session-auth/internal/user"
	"github.com/google/uuid"
)

type Service struct {
	sessionRepository SessionRepository
	userRepository    user.DBRepository
}

const invalidCredentials = "invalid credentials"

func NewService(sessionRepository SessionRepository, userRepository user.DBRepository) *Service {
	return &Service{sessionRepository: sessionRepository, userRepository: userRepository}
}

func (s *Service) Login(ctx context.Context, in user.User) (*SessionMetadata, error) {
	// find user by email
	filter := user.User{
		Email: in.Email,
	}
	result, err := s.userRepository.FindOne(ctx, filter)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, errors.New(invalidCredentials)
	}

	// compare password against the hashed password
	if result.Password != in.Password {
		return nil, errors.New(invalidCredentials)
	}

	// create a new session
	sessionKey := uuid.NewString()
	sessionValue := result.ID
	sessionTTL := "30s"

	err = s.sessionRepository.Set(ctx, sessionKey, sessionValue, sessionTTL)
	if err != nil {
		return nil, err
	}

	ttl, _ := time.ParseDuration(sessionTTL)

	session := SessionMetadata{
		Key: sessionKey,
		TTL: time.Now().Add(ttl),
	}

	return &session, nil
}

func (s *Service) Logout(ctx context.Context, sessionId string) error {
	cache, err := s.sessionRepository.Get(ctx, sessionId)
	if err != nil {
		return err
	}

	user, err := s.userRepository.FindById(ctx, cache)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.New(invalidCredentials)
	}

	return s.sessionRepository.Delete(ctx, sessionId)
}
