package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type LessonItem struct {
	ID              uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey;column:id"`
	LessonSectionID uuid.UUID `gorm:"type:uuid;column:lesson_section_id;not null;index;uniqueIndex:uidx_section_item_pos,priority:1"`

	// Examples: "info", "mcq", "match", "listen_type", "speak", "order_words"
	Kind string `gorm:"size:50;not null;index"`

	Position int `gorm:"not null;default:0;uniqueIndex:uidx_section_item_pos,priority:2"`

	// Fully flexible payload per Kind.
	// Postgres: jsonb, MySQL: json
	Payload datatypes.JSON `gorm:"type:jsonb;not null"`

	// Optional metadata for clients (difficulty, XP, etc.)
	Meta datatypes.JSON `gorm:"type:jsonb;not null;default:'{}'"`

	IsPublished bool `gorm:"not null;default:true"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
