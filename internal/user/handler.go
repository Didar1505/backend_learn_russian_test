package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes — удобная регистрация роутов модуля
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/user", h.GetMe)
	rg.PATCH("/user", h.UpdateMe)
}

// GetMe: GET /me
func (h *Handler) GetMe(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	u, err := h.service.GetMe(c.Request.Context(), userID)
	if err != nil {
		if err == ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user_not_found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}

	c.JSON(http.StatusOK, UserToResponse(*u))
}

// UpdateMe: PATCH /me
func (h *Handler) UpdateMe(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_json"})
		return
	}

	updated, err := h.service.UpdateMe(c.Request.Context(), userID, UserToProfile(req))
	if err != nil {
		if err == ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user_not_found"})
			return
		}
		// Валидационные ошибки сервиса
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, UserToResponse(*updated))
}

// getUserIDFromContext — ожидаем, что auth middleware положил туда userID.
// Поддерживаем несколько форматов, чтобы было меньше боли при интеграции.
func getUserIDFromContext(c *gin.Context) (uuid.UUID, bool) {
	v, exists := c.Get("userID")
	if !exists || v == nil {
		return uuid.Nil, false
	}

	switch t := v.(type) {
	case uuid.UUID:
		if t == uuid.Nil {
			return uuid.Nil, false
		}
		return t, true
	case string:
		id, err := uuid.Parse(t)
		if err != nil || id == uuid.Nil {
			return uuid.Nil, false
		}
		return id, true
	default:
		return uuid.Nil, false
	}
}
