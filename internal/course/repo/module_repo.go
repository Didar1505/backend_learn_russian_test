package repo

import (
	"github.com/Didar1505/project_test.git/internal/course/model"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type ModuleRepository struct {
	DB *gorm.DB
	log *zerolog.Logger
}

func NewModuleRepository(db *gorm.DB, log *zerolog.Logger) *ModuleRepository {
	return &ModuleRepository{
		DB: db,
		log: log,
	}
}

func (r *ModuleRepository) GetModuleWithLessons(id uuid.UUID) (*model.Module, error) {
	var module model.Module
	if err := r.DB.Where("id = ?", id).Preload("Lessons").First(&module).Error; err != nil {
		r.log.Error().Err(err).Msg("Failed to fetch the module with lessons")
		return nil, err
	}
	return &module, nil
}