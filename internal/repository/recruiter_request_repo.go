package repository

import (
	"context"
	"database/sql"
	"job_board/internal/domain"

	"github.com/google/uuid"
)

type RecruiterRequestRepository interface {
	// CreateRecruiterRequest inserts a new recruiter request into the database
	CreateRecruiterRequest(ctx context.Context, req *domain.RecruiterRequest) error

	// GetRecruiterRequestByID retrieves a recruiter request by its ID
	GetRecruiterRequestByUserID(ctx context.Context, recruiterID uuid.UUID) (*domain.RecruiterRequest, error)

}

type PostgresRecruiterRequestRepository struct {
	db *sql.DB
}

func NewPostgresRecruiterRequestRepository(db *sql.DB) *PostgresRecruiterRequestRepository {
	return &PostgresRecruiterRequestRepository{db: db}
}

func (r *PostgresRecruiterRequestRepository) CreateRecruiterRequest(ctx context.Context, req *domain.RecruiterRequest) error {
	query := `
	INSERT INTO recruiter_requests (recruiter_id, company_name, company_website, message, status)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING recruiter_id, status, created_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		req.RecruiterID,
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

func (r *PostgresRecruiterRequestRepository) GetRecruiterRequestByUserID(ctx context.Context, recruiterID uuid.UUID) (*domain.RecruiterRequest, error) {
	query := `
	SELECT id, recruiter_id, company_name, company_website, message, status, created_at
	FROM recruiter_requests
	WHERE recruiter_id = $1
	`

	var req domain.RecruiterRequest
	err := r.db.QueryRowContext(ctx, query, recruiterID).Scan(
		&req.ID,
		&req.RecruiterID,
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

