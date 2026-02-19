package service

import (
	"github.com/Didar1505/project_test.git/internal/course/dto"
	"github.com/google/uuid"
)

type ModuleService interface {
	GetModuleById(id uuid.UUID) (*dto.ModuleDetailDTO, error)
}
