package models

import "time"

// UserType defines the type of a given user
type UserType string

const (
	// UserTypeAdmin the admin will handle the tickets
	UserTypeAdmin UserType = "admin"
	// UserTypeUser the one who will create the tickets
	UserTypeUser UserType = "user"
)

var validUserTypes = map[UserType]bool{
	UserTypeAdmin: true,
	UserTypeUser:  true,
}

// User represents a user of the application
type User struct {
	UserID   int64     `json:"userID"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Type     UserType  `json:"userType"`
	CreateAt time.Time `json:"createdAt"`
}

// HasValidType returns whether the user has a valid type or not
func (u User) HasValidType() bool {
	return validUserTypes[u.Type]
}
