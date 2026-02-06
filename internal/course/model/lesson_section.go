package model

import (
	"time"

	"github.com/google/uuid"
)

type LessonSection struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey;column:id"`
	LessonID uuid.UUID `gorm:"type:uuid;column:lesson_id;not null;index"`

	Title    string `gorm:"size:255;not null"` // e.g. "Vocabulary"
	Position int    `gorm:"not null;default:0;index"`

	CreatedAt time.Time
	UpdatedAt time.Time

	Items []LessonItem `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
