package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type OTPCode struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey;column:id"`
	Email        string    `gorm:"column:email;index"`
	CodeHash     string    `gorm:"column:code_hash"`
	Purpose      string    `gorm:"column:purpose"`
	ExpiresAt    time.Time `gorm:"column:expires_at;index"`
	AttemptsLeft int       `gorm:"column:attempts_left"`
	CreatedAt    time.Time `gorm:"column:created_at"`
}

func (OTPCode) TableName() string { return "otp_codes" }

type OTPRepository interface {
	Create(ctx context.Context, email string, codeHash string, expiresAt time.Time) error
	GetLatestValid(ctx context.Context, email string, now time.Time) (*OTPCode, error)
	DecrementAttempts(ctx context.Context, id uuid.UUID) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
}
