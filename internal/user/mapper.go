package user

import "time"

func UserToResponse(u User) UserResponse {
	resp := UserResponse{
		ID:             u.ID.String(),
		Email:          u.Email,
		AuthProvider:   u.AuthProvider,
		FullName:       u.FullName,
		NativeLanguage: u.NativeLanguage,
		CreatedAt:      u.CreatedAt.UTC().Format(time.RFC3339),
	}

	if u.LastLoginAt != nil {
		s := u.LastLoginAt.UTC().Format(time.RFC3339)
		resp.LastLoginAt = &s
	}

	return resp
}

func UsersToResponse(users []User) []UserResponse {
	out := make([]UserResponse, 0, len(users))
	for _, u := range users {
		out = append(out, UserToResponse(u))
	}
	return out
}


func UserToProfile(req UpdateProfileRequest) ProfilePatch {
	return ProfilePatch{
		FullName:       req.FullName,
		NativeLanguage: req.NativeLanguage,
	}
}