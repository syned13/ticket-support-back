package service

import "github.com/syned13/ticket-support-back/internal/models"

// LoginResponse login response
type LoginResponse struct {
	User  models.User `json:"user"`
	Token string      `json:"token"`
}
