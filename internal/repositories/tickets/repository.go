package repository

import (
	"context"
	"errors"

	"github.com/syned13/ticket-support-back/internal/models"
)

var (
	// ErrNotFound not found
	ErrNotFound = errors.New("not found")
	// ErrNothingToUpdate nothing to update
	ErrNothingToUpdate = errors.New("nothing to update")
)

type Repository interface {
	SaveTicket(ctx context.Context, ticket models.Ticket) (models.Ticket, error)
	GetTicket(ctx context.Context, ticketID int64) (models.Ticket, error)
	GetTickets(ctx context.Context, lastID int64) ([]models.Ticket, int64, error)
	GetTicketsByCreator(ctx context.Context, creatorID int64, lastID int64) ([]models.Ticket, int64, error)
	UpdateTicket(ctx context.Context, ticket models.Ticket) (models.Ticket, error)
	SaveTicketChange(ctx context.Context, ticketChange models.TicketChange) error
	GetTicketChanges(ctx context.Context, creatorID int64) ([]models.TicketChange, error)
}
