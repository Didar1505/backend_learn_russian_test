package user

type UserResponse struct {
	ID             string  `json:"id"`
	Email          *string `json:"email,omitempty"`
	AuthProvider   string  `json:"auth_provider"`
	FullName       *string `json:"full_name,omitempty"`
	NativeLanguage string  `json:"native_language"`
	CreatedAt      string  `json:"created_at"`
	LastLoginAt    *string `json:"last_login_at,omitempty"`
}

type UpdateProfileRequest struct {
	FullName       *string `json:"full_name"`
	NativeLanguage *string `json:"native_language"`
}
