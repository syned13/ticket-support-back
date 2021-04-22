package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/syned13/ticket-support-back/internal/models"
	usersRepo "github.com/syned13/ticket-support-back/internal/repositories/users"
	"github.com/syned13/ticket-support-back/pkg/httputils"
	"golang.org/x/crypto/bcrypt"
)

const (
	tokenSub      = "ticket-support-back"
	tokenDuration = time.Hour * 24
)

var (
	// ErrMissingName missing name
	ErrMissingName = httputils.NewBadRequestError("missing name")
	// ErrMissingPassword missing password
	ErrMissingPassword = httputils.NewBadRequestError("missing password")
	// ErrMissingEmail missing email
	ErrMissingEmail = httputils.NewBadRequestError("missing email")
	// ErrDuplicateFields duplicate fields
	ErrDuplicateFields = httputils.NewBadRequestError("duplicate fields")
	// ErrMissingType missing type
	ErrMissingType = errors.New("missing type")
	// ErrInvalidType invalid type
	ErrInvalidType = errors.New("invalid type")
	// ErrPasswordHashingFailed hashing password failed
	ErrPasswordHashingFailed = errors.New("hashing password failed")
	// ErrInvalidCredentials invalid credentials
	ErrInvalidCredentials = httputils.NewBadRequestError("invalid credentials")
	// ErrGeneratingIDFailed generating id failed
	ErrGeneratingIDFailed = errors.New("generating id failed")
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
	if errors.Is(err, usersRepo.ErrDuplicateField) {
		return models.User{}, ErrDuplicateFields
	}

	if err != nil {
		return models.User{}, err
	}

	createdUser.Password = ""

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
	err := validateLoginParams(email, password)
	if err != nil {
		return LoginResponse{}, err
	}

	user, err := s.repo.GetUserByEmail(ctx, email)
	if errors.Is(err, usersRepo.ErrNotFound) {
		return LoginResponse{}, ErrInvalidCredentials
	}

	if err != nil {
		return LoginResponse{}, err
	}

	if !isPasswordCorrect(password, user.Password) {
		return LoginResponse{}, ErrInvalidCredentials
	}

	token, err := generateToken(user)
	if err != nil {
		return LoginResponse{}, err
	}

	user.Password = ""

	return LoginResponse{User: user, Token: token}, nil
}

func validateLoginParams(email, password string) error {
	if email == "" {
		return ErrMissingEmail
	}

	if password == "" {
		return ErrMissingPassword
	}

	return nil
}

func generateToken(user models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      fmt.Sprint(user.UserID),
		"iss":      tokenSub,
		"userType": user.Type,
		"iat":      time.Now(),
		"exp":      time.Now().Add(tokenDuration), // TODO: make the adding a const
	})

	signedToken, err := token.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
	if err != nil {
		return "", fmt.Errorf("error signing token: " + err.Error())
	}

	return signedToken, nil
}

func isPasswordCorrect(enteredPassword string, storedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(enteredPassword))
	return err == nil
}
