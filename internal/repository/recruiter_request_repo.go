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
	GetMyRecruiterRequestByUserID(ctx context.Context, userID uuid.UUID) (*domain.RecruiterRequest, error)

}

type PostgresRecruiterRequestRepository struct {
	db *sql.DB
}

func NewPostgresRecruiterRequestRepository(db *sql.DB) *PostgresRecruiterRequestRepository {
	return &PostgresRecruiterRequestRepository{db: db}
}

func (r *PostgresRecruiterRequestRepository) CreateRecruiterRequest(ctx context.Context, req *domain.RecruiterRequest) error {
	query := `
	INSERT INTO recruiter_requests (user_id, company_name, company_website, message, status, rejection_reason)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, status, created_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		req.UserID,
		req.CompanyName,
		req.CompanyWebsite,
		req.Message,
		req.Status,
		req.Reason,
	).Scan(&req.ID, &req.Status, &req.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresRecruiterRequestRepository) GetMyRecruiterRequestByUserID(ctx context.Context, userID uuid.UUID) (*domain.RecruiterRequest, error) {
	query := `
	SELECT id, user_id, company_name, company_website, message, status, rejection_reason, created_at, updated_at
	FROM recruiter_requests
	WHERE user_id = $1
	`

	var req domain.RecruiterRequest
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&req.ID,
		&req.UserID,
		&req.CompanyName,
		&req.CompanyWebsite,
		&req.Message,
		&req.Status,
		&req.Reason,
		&req.CreatedAt,
		&req.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound // Not found
		}
		return nil, err
	}

	return &req, nil
}

