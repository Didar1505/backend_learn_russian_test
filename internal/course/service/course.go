package service

import "github.com/Didar1505/project_test.git/internal/course/dto"

type CourseService interface {
	ListPublished() ([]dto.CourseListItemDTO, error)
	GetPublishedBySlug(slug string) (*dto.CourseDetailDTO, error)
}
