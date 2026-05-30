package auth

// RegisterUserResponse represents the user registration response payload.
type RegisterUserResponse struct {
	Data    UserData `json:"data"`
	Message string   `json:"message"`
}

// LoginResponse represents the user login response payload.
type LoginResponse struct {
	Data    string `json:"data"`
	Message string `json:"message"`
}

// UserData represents the user data in the response.
type UserData struct {
	ID          string    `json:"id"`
	DisplayName string    `json:"display_name"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	CreatedAt   interface{} `json:"created_at"`
}

