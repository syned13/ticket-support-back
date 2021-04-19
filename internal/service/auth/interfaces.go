package service

import (
	"context"

	"github.com/syned13/ticket-support-back/internal/models"
)

// Service defines the auth related methods
type Service interface {
	CreateUser(ctx context.Context, user models.User) (models.User, error)
	Login(ctx context.Context, email, password string) (LoginResponse, error)
}
