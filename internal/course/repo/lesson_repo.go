package repo

import (
	"github.com/Didar1505/project_test.git/internal/course/model"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type LessonRepository struct {
	DB  *gorm.DB
	log *zerolog.Logger
}

func NewLessonRepository(db *gorm.DB, log *zerolog.Logger) *LessonRepository {
	return &LessonRepository{DB: db, log: log}
}

func (r *LessonRepository) GetLessonWithSections(id uuid.UUID) (*model.Lesson, error) {
	var lesson model.Lesson
	if err := r.DB.Where("id = ?", id).Preload("Sections").First(&lesson).Error; err != nil {
		r.log.Error().Err(err).Msg("Failed to fetch the lesson with sections from DB")
		return nil, err
	}
	return &lesson, nil
}

func (r *LessonRepository) GetSectionItems(section_id uuid.UUID) (*model.LessonSection, error) {
	var section model.LessonSection
	if err := r.DB.Where("id = ?", section_id).Preload("Items").First(&section).Error; err != nil {
		r.log.Error().Err(err).Msg("Failed to fetch Lesson section and items from DB")
		return nil, err
	}
	return &section, nil
}
