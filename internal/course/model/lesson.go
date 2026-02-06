package model

import (
	"time"

	"github.com/google/uuid"
)

type Lesson struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey;column:id"`
	CourseID uuid.UUID `gorm:"type:uuid;column:course_id;index"`
	ModuleID *uuid.UUID `gorm:"type:uuid;column:module_id;index"` // optional

	Title   string `gorm:"size:255;not null"`
	Summary string `gorm:"type:text"`

	Position        int  `gorm:"not null;default:0;index"`
	EstimatedMinutes int `gorm:"not null;default:0"`
	IsPublished     bool `gorm:"not null;default:false"`

	CreatedAt time.Time
	UpdatedAt time.Time

	Sections []LessonSection `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
