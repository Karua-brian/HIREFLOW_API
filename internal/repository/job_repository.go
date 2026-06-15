package repository

import (
	"context"
	"database/sql"
	"fmt"
	"job_board/internal/domain"

	"github.com/google/uuid"
)

// JobStore defines how the service interacts with persistence
// The service does NOT care whether this is Postgres, MySQL, or memory
type JobRepository interface {

	// Create inserts a new job into the database
	// It should set job.ID and job.CreatedAt.
	Create(ctx context.Context, job *domain.Job) error

	// List retrieves all jobs
	// Later we will add pagination and filtering -> List returns jobs with limit
	//  and offset and also the total number of jobs
	List(ctx context.Context, limit, offset int) ([]domain.Job, int64, error)
}

type PostgresJobRepository struct {
	db *sql.DB
}

func NewPostgresJobRepo(db *sql.DB) *PostgresJobRepository {
	return &PostgresJobRepository{db: db}
}

// Implement Create
func (s *PostgresJobRepository) Create(ctx context.Context, job *domain.Job) error {

	query := `
	INSERT INTO jobs (recruiter_user_id, title, description, company_name, location, salary_range)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, created_at
	`

	// Use QueryRowContext so we respect request cancellation
	err := s.db.QueryRowContext(
		ctx,
		query,
		job.RecruiterUserID,
		job.Title,
		job.Description,
		job.Company,
		job.Location,
		job.Salary,
	).Scan(&job.ID, &job.CreatedAt)
	if job.ID == uuid.Nil {
		return fmt.Errorf("job ID not returned from database")
	}
	
	if err != nil {
		return fmt.Errorf("created job failed: %w", err)
	}

	return nil
}

// Implement list
func (s *PostgresJobRepository) List(ctx context.Context, limit, offset int) ([]domain.Job, int64, error) {

	query := `
	SELECT id, title, description, company_name, location, salary_range, created_at 
	FROM jobs
	ORDER BY created_at DESC
	LIMIT $1 OFFSET $2 
	`

	rows, err := s.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	var jobs []domain.Job

	for rows.Next() { // Iterate over the rows and scan into Job structs
		var job domain.Job
		
		if err := rows.Scan(
			&job.ID,
			&job.Title,
			&job.Description,
			&job.Company,
			&job.Location,
			&job.Salary,
			&job.CreatedAt,
		); err != nil {
			return nil, 0, err
		}

		jobs = append(jobs, job)
	}

	//var total int64
	// err = s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM jobs").Scan(&total)
	// if err != nil {
	// 	return nil, 0, err
	// }


	return jobs, 0, rows.Err()
}