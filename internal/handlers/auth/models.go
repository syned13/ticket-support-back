package handlers

// LoginRequest has the fields for a login request body
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
