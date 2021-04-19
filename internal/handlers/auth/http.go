package handlers

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	authService "github.com/syned13/ticket-support-back/internal/service/auth"
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

	}
}
