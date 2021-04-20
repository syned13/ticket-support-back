package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/syned13/ticket-support-back/internal/models"
	repository "github.com/syned13/ticket-support-back/internal/repositories/tickets"
)

var (
	// ErrMissingPool missing pool
	ErrMissingPool = errors.New("missing pool")
)

type postgresRepository struct {
	pool *pgxpool.Pool
}

// New returns a new postgress repository
func New(pool *pgxpool.Pool) (repository.Repository, error) {
	if pool == nil {
		return nil, ErrMissingPool
	}

	return postgresRepository{
		pool: pool,
	}, nil
}

// SaveTicket saves a ticket in the database
func (r postgresRepository) SaveTicket(ctx context.Context, ticket models.Ticket) (models.Ticket, error) {
	return models.Ticket{}, nil
}

// GetTicket returns a ticket from the database based on the tickeID
func (r postgresRepository) GetTicket(ctx context.Context, ticketID int64) (models.Ticket, error) {
	return models.Ticket{}, nil
}

// GetTickets returns all the tickets
func (r postgresRepository) GetTickets(ctx context.Context) ([]models.Ticket, int64, error) {
	return nil, 0, nil
}

// GetTicketsByCreator returns all the tickets made by a single person
func (r postgresRepository) GetTicketsByCreator(ctx context.Context, id int64) ([]models.Ticket, int64, error) {
	return nil, 0, nil
}

// UpdateTicket updates a ticket
func (r postgresRepository) UpdateTicket(ctx context.Context, ticket models.Ticket) ([]models.Ticket, error) {
	return nil, nil
}
