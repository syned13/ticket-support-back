package repository

import (
	"context"

	"github.com/syned13/ticket-support-back/internal/models"
)

// Repository defines the data-persistance related methods for a user
type Repository interface {
	CreateUser(ctx context.Context, user models.User) (models.User, error)
	GetUser(ctx context.Context, userID int) (models.User, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
}
