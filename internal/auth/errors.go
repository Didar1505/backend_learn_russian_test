package auth

import "errors"

var (
	ErrOTPNotFound     = errors.New("otp_not_found_or_expired")
	ErrSessionNotFound = errors.New("session_not_found_or_expired")
	ErrInvalidRefresh  = errors.New("invalid_refresh_token")
)
