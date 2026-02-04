package otp

import (
	"context"
	"errors"
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrOTPNotFound = errors.New("otp code not found")

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Create(ctx context.Context, email string, codeHash string, expiresAt time.Time) error {
	row := OTPCode{
		Email:        email,
		CodeHash:     codeHash,
		Purpose:      "login",
		ExpiresAt:    expiresAt,
		AttemptsLeft: 5,
	}
	return r.db.WithContext(ctx).Create(&row).Error
}

func (r *GormRepository) GetLatestValid(ctx context.Context, email string, now time.Time) (*OTPCode, error) {
	var row OTPCode
	err := r.db.WithContext(ctx).
		Where("email = ? AND expires_at > ?", email, now).
		Order("created_at DESC").
		First(&row).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOTPNotFound
		}
		return nil, err
	}
	return &row, nil
}

func (r *GormRepository) DecrementAttempts(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).
		Model(&OTPCode{}).
		Where("id = ? AND attempts_left > 0", id).
		Update("attempts_left", gorm.Expr("attempts_left - 1"))

	return res.Error
}

func (r *GormRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&OTPCode{}, "id = ?", id).Error
}
