package auth

import "context"

type Mailer interface {
	SendOTP(ctx context.Context, email string, code string) error
}
