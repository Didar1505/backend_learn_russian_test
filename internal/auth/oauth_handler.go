package auth

import (
	"net/http"

	"github.com/Didar1505/project_test.git/internal/auth/providers/oauth"
	"github.com/gin-gonic/gin"
	goauth "google.golang.org/api/oauth2/v2"
)

func (h *Handler) RegisterOAuthRoutes(rg *gin.RouterGroup) {
	rg.GET("/login", oauth.LoginRedirectHandler)
	rg.GET("/callback", oauth.Auth(), h.GoogleCallback)
}

func (h *Handler) GoogleCallback(c *gin.Context) {
	val, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing_oauth_user"})
		return
	}

	var info *goauth.Userinfo
	switch v := val.(type) {
	case *goauth.Userinfo:
		info = v
	case goauth.Userinfo:
		info = &v
	default:
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_oauth_user"})
		return
	}

	ua := c.GetHeader("User-Agent")
	ip := c.ClientIP()

	resp, err := h.service.LoginWithGoogle(
		c.Request.Context(),
		info.Email,
		info.Name,
		info.Id,
		ua,
		ip,
	)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
