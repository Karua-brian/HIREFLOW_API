package repository

import (
	"context"
	"job_board/internal/domain"
	"database/sql"

	"github.com/lib/pq"
)

type ApplicationTxRepository interface {
	Create(ctx context.Context, app *domain.Application) error 
	Exists(ctx context.Context, jobID, userID int64) (bool, error)
}

// txApplicationStore implements ApplicationStore inside a transaction
type txApplicationRepository struct {
	tx *sql.Tx
}

// Create inserts a new application in the context of a transaction
func (t *txApplicationRepository) Create(ctx context.Context, app *domain.Application) error {

	query := `
	INSERT INTO applications (job_id, user_id)
	VALUES ($1, $2)
	RETURNING id, created_at
	`
	err := t.tx.QueryRowContext(
		ctx,
		query,
		app.JobID,
		app.UserID,
	).Scan(&app.ID, &app.CreatedAt)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			return ErrAlreadyApplied
		}
		return err
	}
	return nil
}

// Exists checks if an application already exists in the context of a transaction
func (t *txApplicationRepository) Exists(ctx context.Context, jobID, userID int64) (bool, error) {
	var exists int
	query := `
	SELECT 1
	FROM applications
	WHERE job_id = $1 AND user_id = $2
	`
	err := t.tx.QueryRowContext(
		ctx,
		query,
		jobID,
		userID,
	).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

