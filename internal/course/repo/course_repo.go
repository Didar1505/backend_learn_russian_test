package repo

import (
	"github.com/Didar1505/project_test.git/internal/course/model"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type CourseRepository struct {
	DB *gorm.DB
	log *zerolog.Logger
}

func NewCourseRepository(db *gorm.DB, log *zerolog.Logger) *CourseRepository {
	return &CourseRepository{
		DB: db,
		log: log,
	}
}

func (r *CourseRepository) GetPublishedCourses() ([]model.Course, error) {
	var courses []model.Course
	if err := r.DB.Where("is_published = ?", true).Find(&courses).Error; err != nil {
		r.log.Error().Err(err).Msg("Failed to fetch published courses")
		return nil, err
	}	
	return courses, nil
}

func (r *CourseRepository) GetCourseWithModules(slug string) (*model.Course, error) {
	var course model.Course
	if err := r.DB.Where("slug = ?", slug).Preload("Modules").First(&course).Error; err != nil {
		r.log.Error().Err(err).Msg("Failed to fetch course with modules")
		return nil, err
	}
	return &course, nil
}
