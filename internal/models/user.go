package models

// UserType defines the type of a given user
type UserType string

const (
	// UserTypeAdmin the admin will handle the tickets
	UserTypeAdmin UserType = "admin"
	// UserTypeUser the one who will create the tickets
	UserTypeUser UserType = "user"
)

// User represents a user of the application
type User struct {
	UserID   string   `json:"userID"`
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Password string   `json:"password"`
	Type     UserType `json:"userType"`
}
