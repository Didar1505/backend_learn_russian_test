package model

import (
	"time"

	"github.com/google/uuid"
)

type Module struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey;column:id"`
	CourseID uuid.UUID `gorm:"type:uuid;column:course_id;index"`
	Title    string    `gorm:"size:255;not null"`
	// Position controls ordering inside a course (1..N)
	Position int `gorm:"not null;default:0;index"`

	CreatedAt time.Time
	UpdatedAt time.Time

	Lessons []Lesson `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
