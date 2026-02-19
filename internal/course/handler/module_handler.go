package handler

import (
	"net/http"

	"github.com/Didar1505/project_test.git/internal/course/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ModuleHandler struct {
	service service.ModuleService
}

func NewModuleHandler(svc service.ModuleService) *ModuleHandler {
	return &ModuleHandler{service: svc}
}

func (h *ModuleHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("modules/:id", h.GetModuleByID)
}

func (h *ModuleHandler) GetModuleByID(c *gin.Context) {
	id := c.Param("id")
	parsedID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid format"})
		return
	}

	data, err := h.service.GetModuleById(parsedID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}