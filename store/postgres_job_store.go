package store

import (
	"context"
	"database/sql"
	"job_board/domain"
)

type PostgresJobStore struct {
	db *sql.DB
}

func NewPostgresJobStore(db *sql.DB) *PostgresJobStore {
	return &PostgresJobStore{db: db}
}

// Implement Create
func (s *PostgresJobStore) Create(ctx context.Context, job *domain.Job) error {

	query := `
	INSERT INTO jobs (title, description, company, created_by)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at
	`

	// Use QueryRowContext so we respect request cancellation
	return s.db.QueryRowContext(
		ctx,
		query,
		job.Title,
		job.Description,
		job.Company,
		job.CreatedBy,
	).Scan(&job.ID, &job.CreatedAt)
}

// Implement list
func (s *PostgresJobStore) List(ctx context.Context, limit, offset int) ([]domain.Job, int64, error) {

	query := `
	SELECT id, title, description, company, created_at, created_by
	FROM jobs
	ORDER BY created_at DESC
	LIMIT $1 OFFSET $2
	`

	rows, err := s.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var jobs []domain.Job

	for rows.Next() {
		var job domain.Job
		
		if err := rows.Scan(
			&job.ID,
			&job.Title,
			&job.Description,
			&job.Company,
			&job.CreatedAt,
			&job.CreatedBy,
		); err != nil {
			return nil, 0, err
		}

		jobs = append(jobs, job)
	}

	var total int64
	err = s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM jobs").Scan(&total)
	if err != nil {
		return nil, 0, err
	}


	return jobs, total, rows.Err()
}