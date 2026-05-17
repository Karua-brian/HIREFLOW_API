package store

import (
	"context"
	"job_board/domain"
)

// JobStore defines how the service interacts with persistence
// The service does NOT care whether this is Postgres, MySQL, or memory
type JobStore interface {

	// Create inserts a new job into the database
	// It should set job.ID and job.CreatedAt.
	Create(ctx context.Context, job *domain.Job) error

	// List retrieves all jobs
	// Later we will add pagination and filtering -> List returns jobs with limit
	//  and offset and also the total number of jobs
	List(ctx context.Context, limit, offset int) ([]domain.Job, int64, error)
}