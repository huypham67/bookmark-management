package response

// ErrorResponse represents an error API response.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// Error creates an error response.
func Error(code string, message string) ErrorResponse {
	return ErrorResponse{
		Error:   code,
		Message: message,
	}
}
