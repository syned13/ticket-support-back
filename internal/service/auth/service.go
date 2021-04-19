package service

import (
	"context"
	"errors"

	"github.com/syned13/ticket-support-back/internal/models"
	usersRepo "github.com/syned13/ticket-support-back/internal/repositories/users"
	"github.com/syned13/ticket-support-back/pkg/httputils"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrMissingName missing name
	ErrMissingName = httputils.NewBadRequestError("missing name")
	// ErrMissingPassword missing password
	ErrMissingPassword = httputils.NewBadRequestError("missing password")
	// ErrMissingEmail missing email
	ErrMissingEmail = httputils.NewBadRequestError("missing email")
	// ErrMissingType missing type
	ErrMissingType = errors.New("missing type")
	// ErrInvalidType invalid type
	ErrInvalidType = errors.New("invalid type")
	// ErrPasswordHashingFailed hashing password failed
	ErrPasswordHashingFailed = errors.New("hashing password failed")
)

var (
	generatePasswordHashFunction func(password []byte, cost int) ([]byte, error)
)

type service struct {
	repo usersRepo.Repository
}

func init() {
	generatePasswordHashFunction = bcrypt.GenerateFromPassword
}

func New(repo usersRepo.Repository) Service {
	return service{
		repo: repo,
	}
}

func (s service) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	err := validateCreateUserParams(user)
	if err != nil {
		return models.User{}, err
	}

	hashedPassword, err := generatePasswordHashFunction([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, ErrPasswordHashingFailed
	}

	user.Password = string(hashedPassword)

	createdUser, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return models.User{}, err
	}

	return createdUser, nil
}

func validateCreateUserParams(user models.User) error {
	if user.Email == "" {
		return ErrMissingEmail
	}

	if user.Name == "" {
		return ErrMissingName
	}

	if user.Password == "" {
		return ErrMissingPassword
	}

	if user.Type == "" {
		return ErrMissingType
	}

	if !user.HasValidType() {
		return ErrInvalidType
	}

	return nil
}

func (s service) Login(ctx context.Context, email, password string) (LoginResponse, error) {
	return LoginResponse{}, nil
}
