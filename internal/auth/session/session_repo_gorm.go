package session

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrSessionNotFound = errors.New("session_not_found_or_expired")

type GormSessionRepository struct {
	db *gorm.DB
}

func NewGormSessionRepository(db *gorm.DB) *GormSessionRepository {
	return &GormSessionRepository{db: db}
}

func (r *GormSessionRepository) Create(ctx context.Context, s *Session) error {
	return r.db.WithContext(ctx).Create(s).Error
}

func (r *GormSessionRepository) GetByRefreshHash(ctx context.Context, refreshHash string, now time.Time) (*Session, error) {
	var s Session
	err := r.db.WithContext(ctx).
		Where("refresh_token_hash = ? AND revoked_at IS NULL AND expires_at > ?", refreshHash, now).
		First(&s).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}
	return &s, nil
}

func (r *GormSessionRepository) Revoke(ctx context.Context, sessionID uuid.UUID, now time.Time) error {
	res := r.db.WithContext(ctx).
		Model(&Session{}).
		Where("id = ? AND revoked_at IS NULL", sessionID).
		Update("revoked_at", now)

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrSessionNotFound
	}
	return nil
}

// Rotate — refresh token rotation: меняем hash + expires_at (и можем хранить revoked_at = null)
func (r *GormSessionRepository) Rotate(ctx context.Context, sessionID uuid.UUID, newRefreshHash string, newExpiresAt time.Time, now time.Time) error {
	updates := map[string]any{
		"refresh_token_hash": newRefreshHash,
		"expires_at":         newExpiresAt,
		"revoked_at":         nil,
	}

	res := r.db.WithContext(ctx).
		Model(&Session{}).
		Where("id = ? AND revoked_at IS NULL", sessionID).
		Updates(updates)

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrSessionNotFound
	}
	_ = now
	return nil
}
