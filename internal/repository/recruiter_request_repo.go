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
	GetRecruiterRequestByUserID(ctx context.Context, requestID uuid.UUID) (*domain.RecruiterRequest, error)

}

type PostgresRecruiterRequestRepository struct {
	db *sql.DB
}

func NewPostgresRecruiterRequestRepository(db *sql.DB) *PostgresRecruiterRequestRepository {
	return &PostgresRecruiterRequestRepository{db: db}
}

func (r *PostgresRecruiterRequestRepository) CreateRecruiterRequest(ctx context.Context, req *domain.RecruiterRequest) error {
	query := `
	INSERT INTO recruiter_requests (request_id, company_name, company_website, message, status, rejection_reason)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING request_id, status, created_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		req.RequestID,
		req.CompanyName,
		req.CompanyWebsite,
		req.Message,
		req.Status,
		req.Reason,
	).Scan(&req.RequestID, &req.Status, &req.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresRecruiterRequestRepository) GetRecruiterRequestByUserID(ctx context.Context, requestID uuid.UUID) (*domain.RecruiterRequest, error) {
	query := `
	SELECT id, request_id, company_name, company_website, message, status, created_at
	FROM recruiter_requests
	WHERE request_id = $1
	`

	var req domain.RecruiterRequest
	err := r.db.QueryRowContext(ctx, query, requestID).Scan(
		&req.ID,
		&req.RequestID,
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

