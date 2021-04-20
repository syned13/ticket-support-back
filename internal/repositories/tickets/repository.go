package repository

import (
	"context"

	"github.com/syned13/ticket-support-back/internal/models"
)

type Repository interface {
	SaveTicket(ctx context.Context, ticket models.Ticket) (models.Ticket, error)
	GetTicket(ctx context.Context, ticketID int64) (models.Ticket, error)
	GetTickets(ctx context.Context) ([]models.Ticket, int64, error)
	GetTicketsByCreator(ctx context.Context, id int64) ([]models.Ticket, int64, error)
	UpdateTicket(ctx context.Context, ticket models.Ticket) ([]models.Ticket, error)
}
