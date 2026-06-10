package repository

import (
	"context"
	"database/sql"
	"job_board/internal/domain"

	"github.com/google/uuid"
)

type AdminRepository interface {

	// ListRecruiterRequests retrieves all recruiter requests with pagination
		ListRecruiterRequests(ctx context.Context, limit, offset int) ([]domain.RecruiterRequest, int64, error)

	// GetRecruiterRequestByID retrieves a recruiter request by its ID
	GetRecruiterRequestByUserID(ctx context.Context, recruiterID uuid.UUID) (*domain.RecruiterRequest, error)

	// UpdateRecruiterRequestStatus updates the status of a recruiter request
	UpdateRecruiterRequestStatus(ctx context.Context, recruiterID uuid.UUID, status string) error
}

type PostgresAdminRepository struct {
	db *sql.DB
}
func NewPostgresAdminRepository(db *sql.DB) *PostgresAdminRepository {
	return &PostgresAdminRepository{db: db}
}

func (r *PostgresAdminRepository) ListRecruiterRequests(ctx context.Context, limit, offset int) ([]domain.RecruiterRequest, int64, error) {
	query := `
	SELECT 
		id, 
		recruiter_id,
		company_name,
		company_website, 
		status, 
		created_at, 
		updated_at
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
			&req.RecruiterID,
			&req.CompanyName,
			&req.CompanyWebsite,
			&req.Status,
			&req.CreatedAt,
			&req.UpdatedAt,
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

func (r *PostgresAdminRepository) GetRecruiterRequestByUserID(ctx context.Context, recruiterID uuid.UUID) (*domain.RecruiterRequest, error) {
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

func (r *PostgresAdminRepository) UpdateRecruiterRequestStatus(ctx context.Context, recruiterID uuid.UUID, status string) error {
	query := `
	UPDATE recruiter_requests
	SET status = $1
	WHERE id = $2
	`

	_, err := r.db.ExecContext(ctx, query, status, recruiterID)
	return err
}