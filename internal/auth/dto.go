package auth

import "github.com/Didar1505/project_test.git/internal/user"

type OTPRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type OTPVerifyRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,len=6"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type AuthResponse struct {
	AccessToken  string           `json:"access_token"`
	RefreshToken string           `json:"refresh_token"`
	User         user.UserResponse `json:"user"`
}
