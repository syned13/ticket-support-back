package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/randallmlough/pgxscan"
	"github.com/syned13/ticket-support-back/internal/models"
	repository "github.com/syned13/ticket-support-back/internal/repositories/tickets"
)

var (
	// ErrMissingPool missing pool
	ErrMissingPool = errors.New("missing pool")
)

var (
	// TODO: look for other errors and map them
	errorCodes = map[string]error{}
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
	query := `INSERT INTO tickets 
				(title, description, ticket_type, severity, ticket_priority, ticket_status, creator_id, owner_id, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
				RETURNING id, created_at, updated_at`

	var ticketID sql.NullInt64
	var createdAt, updatedAt sql.NullTime

	err := r.pool.QueryRow(ctx, query,
		ticket.Title,
		ticket.Description,
		ticket.Type,
		ticket.Severity,
		ticket.Priority,
		ticket.Status,
		ticket.CreatorID,
		ticket.OwnerID).Scan(&ticketID, &createdAt, &updatedAt)

	if pgErr, ok := err.(*pgconn.PgError); ok {
		if _, ok := errorCodes[pgErr.Code]; ok {
			return models.Ticket{}, errorCodes[pgErr.Code]
		}
	}

	if err != nil {
		return models.Ticket{}, err
	}

	// TODO: this to another function
	if ticketID.Valid {
		ticket.TicketID = ticketID.Int64
	}

	if createdAt.Valid {
		ticket.CreatedAt = createdAt.Time
	}

	if updatedAt.Valid {
		ticket.UpdatedAt = updatedAt.Time
	}

	return models.Ticket{}, nil
}

// GetTicket returns a ticket from the database based on the tickeID
func (r postgresRepository) GetTicket(ctx context.Context, ticketID int64) (models.Ticket, error) {
	query := `SELECT * FROM tickers WHERE id = $1`

	ticket := models.Ticket{}

	err := r.pool.QueryRow(ctx, query, ticketID).Scan(
		&ticket.TicketID,
		&ticket.Title,
		&ticket.Description,
		&ticket.Type,
		&ticket.Severity,
		&ticket.Priority,
		&ticket.Status,
		&ticket.CreatorID,
		&ticket.OwnerID,
		&ticket.CreatedAt,
		&ticket.UpdatedAt,
		&ticket.ResolvedAt,
	)

	if err != nil {
		fmt.Println(err)
		return models.Ticket{}, err
	}

	return ticket, nil
}

// GetTickets returns all the tickets
func (r postgresRepository) GetTickets(ctx context.Context, lastID int64) ([]models.Ticket, int64, error) {
	query := `SELECT * FROM tickets WHERE lastID > $1 ORDER BY id LIMIT 10`

	rows, err := r.pool.Query(ctx, query, lastID)
	if err != nil {
		fmt.Println(err)
		return nil, 0, err
	}

	tickets := []models.Ticket{}

	if err := pgxscan.NewScanner(rows).Scan(&tickets); err != nil {
		return nil, 0, err
	}

	return tickets, tickets[len(tickets)-1].TicketID, nil
}

// GetTicketsByCreator returns all the tickets made by a single person
func (r postgresRepository) GetTicketsByCreator(ctx context.Context, creatorID int64, lastID int64) ([]models.Ticket, int64, error) {
	query := `SELECT * FROM tickets WHERE creator_id = $1 AND lastID > $2 ORDER BY id LIMIT 10`

	rows, err := r.pool.Query(ctx, query, creatorID, lastID)
	if err != nil {
		fmt.Println(err)
		return nil, 0, err
	}

	tickets := []models.Ticket{}

	if err := pgxscan.NewScanner(rows).Scan(&tickets); err != nil {
		return nil, 0, err
	}

	return tickets, tickets[len(tickets)-1].TicketID, nil
}

// UpdateTicket updates a ticket
func (r postgresRepository) UpdateTicket(ctx context.Context, ticket models.Ticket) ([]models.Ticket, error) {
	return nil, nil
}
