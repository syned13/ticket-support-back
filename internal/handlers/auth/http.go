package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/syned13/ticket-support-back/internal/models"
	authService "github.com/syned13/ticket-support-back/internal/service/auth"
	"github.com/syned13/ticket-support-back/pkg/httputils"
)

var (
	// ErrMissingContentType missing content type
	ErrMissingContentType = httputils.NewBadRequestError("missing content type")
	// ErrInvalidBody invalid request body
	ErrInvalidBody = httputils.NewBadRequestError("invalid request body")
)

type HTTPHandler interface {
	HandleLogin(ctx context.Context) http.HandlerFunc
	HandleSignup(ctx context.Context) http.HandlerFunc
}

type httpHandler struct {
	service authService.Service
}

func SetupRoutes(ctx context.Context, service authService.Service, router *mux.Router) {
	handler := httpHandler{service: service}

	router.HandleFunc("/login", handler.HandleLogin(ctx)).Methods(http.MethodPost)
	router.HandleFunc("/signup", handler.HandleSignup(ctx)).Methods(http.MethodPost)
}

func (h httpHandler) HandleLogin(ctx context.Context) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

	}
}

func (h httpHandler) HandleSignup(ctx context.Context) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		err := validateContentType(*r)
		if err != nil {
			httputils.RespondWithError(rw, err)
			return
		}

		user := models.User{}

		err = json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			httputils.RespondWithError(rw, ErrInvalidBody)
			return
		}

		user.Type = models.UserTypeUser

		user, err = h.service.CreateUser(ctx, user)
		if err != nil {
			fmt.Println("creating_user_failed: " + err.Error())
			httputils.RespondWithError(rw, err)
			return
		}

		httputils.RespondJSON(rw, http.StatusCreated, user)
	}
}

func validateContentType(r http.Request) error {
	if r.Header.Get("Content-Type") == "" {
		return ErrMissingContentType
	}

	return nil
}
