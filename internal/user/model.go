package user

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID              uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey;column:id"`
	Email           *string    `gorm:"column:email"`
	AuthProvider    string     `gorm:"column:auth_provider"`
	FullName        *string    `gorm:"column:full_name"`
	NativeLanguage  string     `gorm:"column:native_language"`
	CreatedAt       time.Time  `gorm:"column:created_at"`
	LastLoginAt     *time.Time `gorm:"column:last_login_at"`
	ProviderSubject *string    `gorm:"column:provider_subject" json:"-"`
}

func (User) TableName() string {
	return "users"
}
