package store

import (
	"context"
	"job_board/domain"
)

type ApplicationStore interface {

	// Create inserts a new application for a job.
	Create(ctx context.Context, app *domain.Application) error

	// Exists checks if user already applied
	Exists(ctx context.Context, jobID, userID int64) (bool, error)

	// Inserts a trasactional new applicaation for a job
	CreateTx(ctx context.Context, fn func(ApplicationTxStore) error) error
}

