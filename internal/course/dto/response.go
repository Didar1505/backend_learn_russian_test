package dto

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// COURSE ONLY
type CourseListItemDTO struct {
	ID           uuid.UUID `json:"id"`
	Slug         string    `json:"slug"`
	Title        string    `json:"title"`
	Level        string    `json:"level"`
	LanguageFrom string    `json:"language_from"`
	LanguageTo   string    `json:"language_to"`
}
type CourseDetailDTO struct {
	ID           uuid.UUID   `json:"id"`
	Slug         string      `json:"slug"`
	Title        string      `json:"title"`
	Description  string      `json:"description"`
	Level        string      `json:"level"`
	LanguageFrom string      `json:"language_from"`
	LanguageTo   string      `json:"language_to"`
	Modules      []ModuleDTO `json:"modules"`
}
type ModuleDTO struct {
	ID       uuid.UUID           `json:"id"`
	Title    string              `json:"title"`
	Position int                 `json:"position"`
	Lessons  []LessonListItemDTO `json:"lessons"`
}

// MODULE ONLY
type LessonListItemDTO struct {
	ID               uuid.UUID `json:"id"`
	Title            string    `json:"title"`
	Summary          string    `json:"summary"`
	Position         int       `json:"position"`
	EstimatedMinutes int       `json:"estimated_minutes"`
}
type ModuleDetailDTO struct {
	ID       uuid.UUID           `json:"id"`
	Title    string              `json:"title"`
	Position int                 `json:"position"`
	Lessons  []LessonListItemDTO `json:"lessons"`
}

type LessonItemDTO struct {
	ID       uuid.UUID      `json:"id"`
	Kind     string         `json:"kind"`
	Position int            `json:"position"`
	Payload  datatypes.JSON `json:"payload"`
	Meta     datatypes.JSON `json:"meta"`
}
type LessonSectionDTO struct {
	ID       uuid.UUID       `json:"id"`
	Title    string          `json:"title"`
	Position int             `json:"position"`
	Items    []LessonItemDTO `json:"items,omitempty"`
}
type LessonDetailDTO struct {
	ID               uuid.UUID          `json:"id"`
	CourseID         uuid.UUID          `json:"course_id"`
	ModuleID         *uuid.UUID         `json:"module_id,omitempty"`
	Title            string             `json:"title"`
	Summary          string             `json:"summary"`
	Position         int                `json:"position"`
	EstimatedMinutes int                `json:"estimated_minutes"`
	SectionCount     int                `json:"section_count"`
	Sections         []LessonSectionDTO `json:"sections"`
}
