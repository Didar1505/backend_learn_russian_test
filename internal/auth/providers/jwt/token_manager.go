package jwt

import (
	"time"

	"github.com/google/uuid"
)

type TokenManager interface {
	SignAccess(userID uuid.UUID, ttl time.Duration) (string, error)
	VerifyAccess(token string) (uuid.UUID, error)
}
