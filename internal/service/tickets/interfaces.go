package service

import (
	"context"

	"github.com/syned13/ticket-support-back/internal/models"
)

type Service interface {
	CreateTicket(ctx context.Context, ticket models.Ticket) (models.Ticket, error)
	GetTickets(ctx context.Context, userID int64, userType models.UserType, lastID int64) (GetTicketsResponse, error)
	GetTicket(ctx context.Context, ticketID int64) (models.Ticket, error)
}
