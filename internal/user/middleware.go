package user

import "github.com/gin-gonic/gin"

func FakeAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Подставь реальный user uuid из БД
		c.Set("userID", "550e8400-e29b-41d4-a716-446655440000")
		c.Next()
	}
}
