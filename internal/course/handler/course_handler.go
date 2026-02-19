package handler

import (
	"net/http"

	"github.com/Didar1505/project_test.git/internal/course/service"
	"github.com/gin-gonic/gin"
)

type CourseHandler struct {
	service service.CourseService
}

func NewCourseHandler(svc service.CourseService) *CourseHandler {
	return &CourseHandler{service: svc}
}

func (h *CourseHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/courses", h.ListCourses)
	r.GET("/courses/:slug", h.GetCourseBySlug)
}

func (h *CourseHandler) ListCourses(c *gin.Context) {
	data, err := h.service.ListPublished()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} 
	c.JSON(http.StatusOK, data)
}

func (h *CourseHandler) GetCourseBySlug(c *gin.Context) {
	slug := c.Param("slug")
	data, err := h.service.GetPublishedBySlug(slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, data)
	}
	c.JSON(http.StatusOK, data)
}