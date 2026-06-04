package repository

import (
	"context"
	"database/sql"
	"job_board/internal/domain"
)

type RecruiterRequestRepository interface {
	// CreateRecruiterRequest inserts a new recruiter request into the database
	CreateRecruiterRequest(ctx context.Context, req *domain.RecruiterRequest) error

	// GetRecruiterRequestByID retrieves a recruiter request by its ID
	GetRecruiterRequestByID(ctx context.Context, id int64) (*domain.RecruiterRequest, error)

	// ListRecruiterRequests retrieves all recruiter requests with pagination
	ListRecruiterRequests(ctx context.Context, limit, offset int) ([]domain.RecruiterRequest, int64, error)

	// UpdateRecruiterRequestStatus updates the status of a recruiter request
	UpdateRecruiterRequestStatus(ctx context.Context, id int64, status string) error
}

type PostgresRecruiterRequestRepository struct {
	db *sql.DB
}

func NewPostgresRecruiterRequestRepository(db *sql.DB) *PostgresRecruiterRequestRepository {
	return &PostgresRecruiterRequestRepository{db: db}
}

func (r *PostgresRecruiterRequestRepository) CreateRecruiterRequest(ctx context.Context, req *domain.RecruiterRequest) error {
	query := `
	INSERT INTO recruiter_requests (user_id, company_name, company_website, message, status)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, created_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		req.UserID,
		req.CompanyName,
		req.CompanyWebsite,
		req.Message,
		req.Status,
	).Scan(&req.ID, &req.Status, &req.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresRecruiterRequestRepository) GetRecruiterRequestByID(ctx context.Context, id int64) (*domain.RecruiterRequest, error) {
	query := `
	SELECT id, user_id, company_name, company_website, message, status, created_at
	FROM recruiter_requests
	WHERE id = $1
	`

	var req domain.RecruiterRequest
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&req.ID,
		&req.UserID,
		&req.CompanyName,
		&req.CompanyWebsite,
		&req.Message,
		&req.Status,
		&req.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, err
	}

	return &req, nil
}

func (r *PostgresRecruiterRequestRepository) ListRecruiterRequests(ctx context.Context, limit, offset int) ([]domain.RecruiterRequest, int64, error) {
	query := `
	SELECT id, user_id, company_name, company_website, message, status, created_at
	FROM recruiter_requests
	ORDER BY created_at DESC
	LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var requests []domain.RecruiterRequest
	for rows.Next() {
		var req domain.RecruiterRequest
		if err := rows.Scan(
			&req.ID,
			&req.UserID,
			&req.CompanyName,
			&req.CompanyWebsite,
			&req.Message,
			&req.Status,
			&req.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		requests = append(requests, req)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	// Get total count for pagination
	var total int64
	err = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM recruiter_requests`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return requests, total, nil
}

func (r *PostgresRecruiterRequestRepository) UpdateRecruiterRequestStatus(ctx context.Context, id int64, status string) error {
	query := `
	UPDATE recruiter_requests
	SET status = $1
	WHERE id = $2
	`

	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}