package store

import (
	"context"
	"job_board/domain"
)

type ApplicationTxStore interface {
	Create(ctx context.Context, app *domain.Application) error 
	Exists(ctx context.Context, jobID, userID int64) (bool, error)
}