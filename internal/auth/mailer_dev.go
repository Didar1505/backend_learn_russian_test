package auth

import (
	"context"
	"log"
)

type DevMailer struct{}

func NewDevMailer() *DevMailer { return &DevMailer{} }

func (m *DevMailer) SendOTP(ctx context.Context, email, code string) error {
	log.Printf("[DEV OTP] email=%s code=%s\n", email, code)
	return nil
}
