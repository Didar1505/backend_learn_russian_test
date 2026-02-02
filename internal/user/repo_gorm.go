package user

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Create(ctx context.Context, u *User) error {
	// Если ты генеришь UUID в Go:
	// if u.ID == uuid.Nil { u.ID = uuid.New() }

	return r.db.WithContext(ctx).Create(u).Error
}

func (r *GormRepository) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	var u User
	err := r.db.WithContext(ctx).First(&u, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *GormRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	var u User
	err := r.db.WithContext(ctx).First(&u, "email = ?", email).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *GormRepository) GetByProviderSubject(ctx context.Context, provider string, subject string) (*User, error) {
	var u User

	// Если поля provider_subject пока нет — этот метод лучше не использовать
	// до добавления миграции и поля в model.
	err := r.db.WithContext(ctx).First(&u,
		"auth_provider = ? AND provider_subject = ?",
		provider, subject,
	).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *GormRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	now := time.Now().UTC()
	res := r.db.WithContext(ctx).
		Model(&User{}).
		Where("id = ?", id).
		Update("last_login_at", now)

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *GormRepository) UpdateProfile(ctx context.Context, id uuid.UUID, patch ProfilePatch) (*User, error) {
	updates := map[string]any{}

	if patch.FullName != nil {
		updates["full_name"] = patch.FullName
	}
	if patch.NativeLanguage != nil {
		updates["native_language"] = *patch.NativeLanguage
	}

	if len(updates) == 0 {
		// Нечего обновлять — вернём текущего пользователя
		return r.GetByID(ctx, id)
	}

	res := r.db.WithContext(ctx).
		Model(&User{}).
		Where("id = ?", id).
		Updates(updates)

	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, ErrNotFound
	}

	return r.GetByID(ctx, id)
}
