package profile

import "time"

// UserData represents the user data in the profile response.
type UserData struct {
	ID          string    `json:"id"`
	DisplayName string    `json:"display_name"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	CreatedAt   time.Time `json:"created_at"`
}

// UserResponse represents the user info response payload.
type UserResponse struct {
	Data    *UserData `json:"data"`
	Message string    `json:"message"`
}

// UpdateUserResponse represents the user update response payload.
type UpdateUserResponse struct {
	Message string `json:"message"`
}

