package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID               uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey;column:id"`
	UserID           uuid.UUID  `gorm:"type:uuid;column:user_id;index"`
	RefreshTokenHash string     `gorm:"column:refresh_token_hash;uniqueIndex"`
	UserAgent        *string    `gorm:"column:user_agent"`
	IP               *string    `gorm:"column:ip"`
	ExpiresAt        time.Time  `gorm:"column:expires_at;index"`
	CreatedAt        time.Time  `gorm:"column:created_at"`
	RevokedAt        *time.Time `gorm:"column:revoked_at"`
}

func (Session) TableName() string { return "sessions" }

type SessionRepository interface {
	Create(ctx context.Context, s *Session) error
	GetByRefreshHash(ctx context.Context, refreshHash string, now time.Time) (*Session, error)
	Revoke(ctx context.Context, sessionID uuid.UUID, now time.Time) error
	Rotate(ctx context.Context, sessionID uuid.UUID, newRefreshHash string, newExpiresAt time.Time, now time.Time) error
}
