package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
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
				(title, ticket_description, ticket_type, severity, ticket_priority, ticket_status, creator_id, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
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
		ticket.CreatorID).Scan(&ticketID, &createdAt, &updatedAt)

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
		ticket.CreatedAt = &createdAt.Time
	}

	if updatedAt.Valid {
		ticket.UpdatedAt = &updatedAt.Time
	}

	return ticket, nil
}

// GetTicket returns a ticket from the database based on the tickeID
func (r postgresRepository) GetTicket(ctx context.Context, ticketID int64) (models.Ticket, error) {
	query := `SELECT * FROM tickets WHERE id = $1`

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

	if errors.Is(err, pgx.ErrNoRows) {
		return models.Ticket{}, repository.ErrNotFound
	}

	if err != nil {
		return models.Ticket{}, err
	}

	return ticket, nil
}

// GetTickets returns all the tickets
func (r postgresRepository) GetTickets(ctx context.Context, lastID int64) ([]models.Ticket, int64, error) {
	query := `SELECT * FROM tickets WHERE id > $1 ORDER BY id LIMIT 1000`

	rows, err := r.pool.Query(ctx, query, lastID)
	if err != nil {
		fmt.Println(err)
		return nil, 0, err
	}

	tickets := []models.Ticket{}

	if err := pgxscan.NewScanner(rows).Scan(&tickets); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, 0, repository.ErrNotFound
		}

		return nil, 0, err
	}

	return tickets, tickets[len(tickets)-1].TicketID, nil
}

// GetTicketsByCreator returns all the tickets made by a single person
func (r postgresRepository) GetTicketsByCreator(ctx context.Context, creatorID int64, lastID int64) ([]models.Ticket, int64, error) {
	query := `SELECT * FROM tickets WHERE creator_id = $1 AND id > $2 ORDER BY id LIMIT 1000`

	rows, err := r.pool.Query(ctx, query, creatorID, lastID)
	if err != nil {
		fmt.Println(err)
		return nil, 0, err
	}

	tickets := []models.Ticket{}

	if err := pgxscan.NewScanner(rows).Scan(&tickets); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, 0, repository.ErrNotFound
		}

		return nil, 0, err
	}

	return tickets, tickets[len(tickets)-1].TicketID, nil
}

// UpdateTicket updates a ticket
func (r postgresRepository) UpdateTicket(ctx context.Context, ticket models.Ticket) (models.Ticket, error) {
	params := []interface{}{}
	setStatements := []string{}
	valuesCount := 0

	if ticket.Status != "" {
		valuesCount++
		params = append(params, ticket.Status)
		// setStatements = append(setStatements, fmt.Sprintf("ticket_status = $%d", valuesCount))
		setStatements = append(setStatements, fmt.Sprintf("ticket_status = '%s'", ticket.Status))
	}

	if ticket.OwnerID != nil {
		valuesCount++
		params = append(params, ticket.Status)
		// setStatements = append(setStatements, fmt.Sprintf("owner_id = $%d", valuesCount))
		setStatements = append(setStatements, fmt.Sprintf("owner_id = '%d'", *ticket.OwnerID))
	}

	if valuesCount == 0 {
		return models.Ticket{}, repository.ErrNothingToUpdate
	}

	setStatements = append(setStatements, "updated_at = NOW()")
	params = append(params, ticket.TicketID)

	query := fmt.Sprintf("UPDATE tickets SET %s WHERE id = %d RETURNING *", strings.Join(setStatements, ", "), ticket.TicketID)

	updatedTicket := models.Ticket{}

	_, err := r.pool.Exec(ctx, query)
	if err != nil {
		return models.Ticket{}, err
	}

	return updatedTicket, nil
}

func (r postgresRepository) SaveTicketChange(ctx context.Context, ticketChange models.TicketChange) error {
	query := `INSERT INTO tickets_changes 
			  (ticket_id, creator_id, to_status, changed_at) 
			  VALUES ($1, $2, $3, NOW())`

	_, err := r.pool.Exec(ctx, query, ticketChange.TicketID, ticketChange.CreatorID, ticketChange.To)
	if err != nil {
		return err
	}

	return nil
}

func (r postgresRepository) GetTicketChanges(ctx context.Context, creatorID int64) ([]models.TicketChange, error) {
	query := `SELECT * FROM tickets_changes where creator_id = $1`

	changes := []models.TicketChange{}

	rows, err := r.pool.Query(ctx, query, creatorID)
	if errors.Is(err, pgx.ErrNoRows) {
		return changes, nil
	}

	if err := pgxscan.NewScanner(rows).Scan(&changes); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}

		return nil, err
	}

	return changes, nil
}
