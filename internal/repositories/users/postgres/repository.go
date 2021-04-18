package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/syned13/ticket-support-back/internal/models"
	repository "github.com/syned13/ticket-support-back/internal/repositories/users"
)

var (
	// ErrMissingPool missing pool
	ErrMissingPool = errors.New("missing pool")
)

type PgxIface interface {
	Begin(context.Context) (pgx.Tx, error)
	Close(context.Context) error
}

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

// CreateUser saves a user in the database
func (r postgresRepository) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	query := `INSERT INTO users
			(name, email, password, user_type, created_at)
			VALUES ($1, $2, $3, $4, NOW() )
			RETURNING id, created_at `

	var userID sql.NullInt64
	var createdAt sql.NullTime

	err := r.pool.QueryRow(ctx, query, user.Name, user.Password, user.Type).Scan(&userID, &createdAt)
	if err != nil {
		fmt.Println(err)
		fmt.Println(err.Error())
		// TODO: handle postgres specific errors
		return models.User{}, err
	}

	if userID.Valid {
		user.UserID = userID.Int64
	}

	if createdAt.Valid {
		user.CreateAt = createdAt.Time
	}

	return user, nil
}

// GetUser gets a user from the database based on the userID
func (r postgresRepository) GetUser(ctx context.Context, userID int) (models.User, error) {
	query := `SELECT * FROM users WHERE id = $1`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return models.User{}, err
	}

	user := models.User{}

	err = rows.Scan(&user.UserID, &user.Name, &user.Password, &user.Type, &user.CreateAt)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// GetUserByEmail returns a user from the dabase based on the email
func (r postgresRepository) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	query := `SELECT * FROM users WHERE email = $1`

	rows, err := r.pool.Query(ctx, query, email)
	if err != nil {
		return models.User{}, err
	}

	user := models.User{}

	err = rows.Scan(&user.UserID, &user.Name, &user.Password, &user.Type, &user.CreateAt)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
