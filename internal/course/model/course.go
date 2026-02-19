package model

import (
	"time"

	"github.com/google/uuid"
)

type Course struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey;column:id"`
	Slug         string    `gorm:"size:120;not null;uniqueIndex"`
	Title        string    `gorm:"size:255;not null"`
	Description  string    `gorm:"type:text"`
	Level        string    `gorm:"size:10;not null"` // "A1", "A2", "B1"...
	LanguageFrom string    `gorm:"size:10;not null"` // "en", "de", "kk"...
	LanguageTo   string    `gorm:"size:10;not null"` // "ru"
	IsPublished  bool      `gorm:"column:is_published;not null;default:false"`

	CreatedAt time.Time
	UpdatedAt time.Time

	Modules []Module `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Lessons []Lesson `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
