package service

import (
	"github.com/Didar1505/project_test.git/internal/course/dto"
	"github.com/google/uuid"
)

type LessonService interface {
	GetLessonWithSections(uuid.UUID) (*dto.LessonDetailDTO, error)
}