package service

import (
	"context"

	"github.com/syned13/ticket-support-back/internal/models"
	usersRepo "github.com/syned13/ticket-support-back/internal/repositories/users"
)

type service struct {
	repo usersRepo.Repository
}

func New(repo usersRepo.Repository) Service {
	return service{
		repo: repo,
	}
}

func (s service) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	return models.User{}, nil
}

func (s service) Login(ctx context.Context, email, password string) (LoginResponse, error) {
	return LoginResponse{}, nil
}
