package service

import (
	"errors"
	"strings"

	"github.com/Didar1505/project_test.git/internal/course/dto"
	"github.com/Didar1505/project_test.git/internal/course/model"
	"github.com/Didar1505/project_test.git/internal/course/repo"
)

type CourseServiceImpl struct {
	repo repo.CourseRepository
}

func NewCourseService(repo repo.CourseRepository) *CourseServiceImpl {
	return &CourseServiceImpl{repo: repo}
}

func (s *CourseServiceImpl) ListPublished() ([]dto.CourseListItemDTO, error) {
	courses, err := s.repo.GetPublishedCourses()
	if err != nil {
		return nil, err
	}
	out := make([]dto.CourseListItemDTO, 0, len(courses))
	for _, c := range courses {
		out = append(out, dto.CourseListItemDTO{
			ID:           c.ID,
			Slug:         c.Slug,
			Title:        c.Title,
			Level:        c.Level,
			LanguageFrom: c.LanguageFrom,
			LanguageTo:   c.LanguageTo,
		})
	}
	return out, nil
}

func (s *CourseServiceImpl) GetPublishedBySlug(slug string) (*dto.CourseDetailDTO, error) {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return nil, errors.New("Slug is required")
	}

	course, err := s.repo.GetCourseWithModules(slug)
	if err != nil {
		return nil, errors.New("Course not found")
	}

	if !course.IsPublished {
		return nil, errors.New("course not found")
	}

	return mapCourseToDetailDTO(course), nil
}

func mapCourseToDetailDTO(c *model.Course) *dto.CourseDetailDTO {
	modules := make([]dto.ModuleDTO, 0, len(c.Modules))
	for _, m := range c.Modules {
		modules = append(modules, dto.ModuleDTO{
			ID:       m.ID,
			Title:    m.Title,
			Position: m.Position,
		})
	}

	return &dto.CourseDetailDTO{
		ID:           c.ID,
		Slug:         c.Slug,
		Title:        c.Title,
		Description:  c.Description,
		Level:        c.Level,
		LanguageFrom: c.LanguageFrom,
		LanguageTo:   c.LanguageTo,
		Modules:      modules,
	}
}