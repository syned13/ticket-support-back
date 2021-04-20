package repository

import (
	"context"
	"errors"

	"github.com/syned13/ticket-support-back/internal/models"
)

var (
	// ErrDuplicateField duplicate field
	ErrDuplicateField = errors.New("duplicate field")
	// ErrNotFound not found
	ErrNotFound = errors.New("not found")
)

// Repository defines the data-persistance related methods for a user
type Repository interface {
	CreateUser(ctx context.Context, user models.User) (models.User, error)
	GetUser(ctx context.Context, userID int) (models.User, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
}
