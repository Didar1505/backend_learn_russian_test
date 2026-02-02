package user

import (
	"context"

	"github.com/google/uuid"
)

// Repository — доменный контракт.
// Ни HTTP, ни GORM сюда не лезут (кроме типов uuid/ctx).
type Repository interface {
	// Create сохраняет нового пользователя.
	// Ожидается, что ID/CreatedAt могут быть заполнены на уровне БД или сервиса.
	Create(ctx context.Context, u *User) error

	// GetByID получает пользователя по ID.
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)

	// GetByEmail получает пользователя по email.
	// Возвращает (nil, ErrNotFound), если не найден.
	GetByEmail(ctx context.Context, email string) (*User, error)

	// GetByProviderSubject для соцлогинов: (google|vk, subject/id у провайдера).
	// Если пока не используешь provider_subject — можешь временно оставить заглушкой/не вызывать.
	GetByProviderSubject(ctx context.Context, provider string, subject string) (*User, error)

	// UpdateLastLogin фиксирует факт входа.
	UpdateLastLogin(ctx context.Context, id uuid.UUID) error

	// UpdateProfile — безопасное обновление редактируемых полей.
	UpdateProfile(ctx context.Context, id uuid.UUID, patch ProfilePatch) (*User, error)
}

// ProfilePatch — что пользователь реально может менять (не auth_provider, не created_at).
type ProfilePatch struct {
	FullName       *string
	NativeLanguage *string
}
