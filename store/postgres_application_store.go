package store

import (
	"context"
	"database/sql"
	"job_board/domain"
	"github.com/jackc/pgconn"
	"errors"
)

type PostgresApplicationStore struct {
	db *sql.DB
}

func NewPostgresApplicationStore(db *sql.DB) * PostgresApplicationStore {
	return &PostgresApplicationStore{db: db}
}

// Create inserts a new application with idempotency safety
func (s *PostgresApplicationStore) Create(ctx context.Context, app *domain.Application) error {
	
	query := `
	INSERT INTO applications (job_id, user_id)
	VALUES ($1, $2)
	RETURNING id, created_at
	`

	err := s.db.QueryRowContext(
		ctx,
		query,
		app.JobID,
		app.UserID,
	) .Scan(&app.ID, &app.CreatedAt)

	if err != nil {
		// Detect unique violation
		var pgError *pgconn.PgError

		if errors.As(err, &pgError) && pgError.Code == "23505" {
			return ErrDuplicate // Application already exists, treat as success
		}
	}
	return err
}

func (s *PostgresApplicationStore) Exists(ctx context.Context, jobID, userID int64) (bool, error) {

	query := `
	SELECT 1
	FROM applications
	WHERE job_id = $1 AND user_id = $2
	`

	var exists int 
	err := s.db.QueryRowContext(ctx, query, jobID, userID).Scan(&exists)

	if errors.Is(err, sql.ErrNoRows) {
		return false, ErrNotFound // No application found
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *PostgresApplicationStore) CreateTx(ctx context.Context, fn func(ApplicationTxStore) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Wrap store with tx
	txStore := &txApplicationStore{tx: tx}

	if err := fn(txStore); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
