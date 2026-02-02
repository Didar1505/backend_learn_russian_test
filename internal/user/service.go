package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid user id")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetMe(ctx context.Context, id uuid.UUID) (*User, error) {
	return s.GetByID(ctx, id)
}

func (s *Service) UpdateMe(ctx context.Context, id uuid.UUID, patch ProfilePatch) (*User, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid user id")
	}

	if patch.NativeLanguage != nil {
		if *patch.NativeLanguage == "" {
			return nil, errors.New("native_language cannot be empty")
		}
	}

	return s.repo.UpdateProfile(ctx, id, patch)
}

func (s *Service) TouchLastLogin(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("invalid user id")
	}
	return s.repo.UpdateLastLogin(ctx, id)
}
