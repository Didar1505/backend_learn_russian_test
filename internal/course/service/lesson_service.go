package service

import (
	"errors"

	"github.com/Didar1505/project_test.git/internal/course/dto"
	"github.com/Didar1505/project_test.git/internal/course/repo"
	"github.com/google/uuid"
)

type LessonServiceImpl struct {
	repo repo.LessonRepository
}

func NewLessonService(repo repo.LessonRepository) *LessonServiceImpl {
	return &LessonServiceImpl{repo: repo}
}

func (s *LessonServiceImpl) GetLessonWithSections(id uuid.UUID) (*dto.LessonDetailDTO, error) {
	lesson, err := s.repo.GetLessonWithSections(id)
	if err != nil {
		return nil, errors.New("Lesson not found")
	}
	if !lesson.IsPublished {
		return nil, errors.New("lesson not found")
	}

	sections := make([]dto.LessonSectionDTO, 0, len(lesson.Sections))
	for _, sec := range lesson.Sections {
		sections = append(sections, dto.LessonSectionDTO{
			ID:       sec.ID,
			Title:    sec.Title,
			Position: sec.Position,
		})
	}

	return &dto.LessonDetailDTO{
		ID:               lesson.ID,
		CourseID:         lesson.CourseID,
		ModuleID:         lesson.ModuleID,
		Title:            lesson.Title,
		Summary:          lesson.Summary,
		Position:         lesson.Position,
		EstimatedMinutes: lesson.EstimatedMinutes,
		SectionCount:     len(sections),
		Sections:         sections,
	}, nil
}
