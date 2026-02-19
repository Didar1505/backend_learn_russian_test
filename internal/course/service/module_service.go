package service

import (
	"errors"

	"github.com/Didar1505/project_test.git/internal/course/dto"
	"github.com/Didar1505/project_test.git/internal/course/repo"
	"github.com/google/uuid"
)

type ModuleServiceImpl struct {
	repo repo.ModuleRepository
}

func NewModuleService(repo repo.ModuleRepository) *ModuleServiceImpl {
	return &ModuleServiceImpl{repo: repo}
}

func (s *ModuleServiceImpl) GetModuleById(id uuid.UUID) (*dto.ModuleDetailDTO, error) {
	module, err := s.repo.GetModuleWithLessons(id)
	if err != nil {
		return nil, errors.New("module not found")
	}

	lessons := make([]dto.LessonListItemDTO, 0)
	for _, l := range module.Lessons {
		if !l.IsPublished {
			continue
		}

		lessons = append(lessons, dto.LessonListItemDTO{
			ID:                l.ID,
			Title:             l.Title,
			Summary:           l.Summary,
			Position:          l.Position,
			EstimatedMinutes:  l.EstimatedMinutes,
		})
	}

	return &dto.ModuleDetailDTO{
		ID:       module.ID,
		Title:    module.Title,
		Position: module.Position,
		Lessons:  lessons,
	}, nil
}
