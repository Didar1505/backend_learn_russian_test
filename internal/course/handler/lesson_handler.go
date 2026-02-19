package handler

import (
	"net/http"

	"github.com/Didar1505/project_test.git/internal/course/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LessonHandler struct {
	service service.LessonService
}

func NewLessonHandler(svc service.LessonService) *LessonHandler {
	return &LessonHandler{service: svc}
}

func (h *LessonHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("lessons/:id", h.GetLessonWithSections)
} 

func (h *LessonHandler) GetLessonWithSections(c *gin.Context) {
	id := c.Param("id")
	parsedID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid format"})
		return
	}
	lesson, err := h.service.GetLessonWithSections(parsedID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, lesson)
}