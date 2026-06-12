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
	GetRecruiterRequestByUserID(ctx context.Context, requestID uuid.UUID) (*domain.RecruiterRequest, error)

	//
	ApproveRecruiterRequest(ctx context.Context, requestID uuid.UUID) error

	//
	RejectRecruiterRequest(ctx context.Context, reason string, requestID uuid.UUID) error
	
	// UpdateUserRole updates the role of a user
	UpdateUserRole(ctx context.Context, userID uuid.UUID, role string) error

	// UpdateRecruiterRequestStatus updates the status of a recruiter request
	// UpdateRecruiterRequestStatus(ctx context.Context, requestID uuid.UUID, status string) error

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
		request_id,
		company_name,
		message,
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
			&req.RequestID,
			&req.CompanyName,
			&req.Message,
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

func (r *PostgresAdminRepository) GetRecruiterRequestByUserID(ctx context.Context, requestID uuid.UUID) (*domain.RecruiterRequest, error) {
	query := `
	SELECT id, request_id, company_name, message, status, created_at
	FROM recruiter_requests
	WHERE request_id = $1
	`

	var req domain.RecruiterRequest
	err := r.db.QueryRowContext(ctx, query, requestID).Scan(
		&req.ID,
		&req.RequestID,
		&req.CompanyName,
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

func (r *PostgresAdminRepository) ApproveRecruiterRequest(ctx context.Context, requestID uuid.UUID) error {
	
	query := `
	UPDATE recruiter_requests
	SET status = 'approved',
		updated_at = NOW()
	WHERE request_id = $2
	`

	result, err := r.db.ExecContext(ctx, query, requestID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err 
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil 
} 

func (r *PostgresAdminRepository) RejectRecruiterRequest(ctx context.Context, reason string, requestID uuid.UUID) error {
	
	query := `
	UPDATE recruiter_requests 
	SET status = 'rejected'
		rejection_reason = $1,
		updated_at = NOW()
	WHERE request_id = $2
	`

	result, err := r.db.ExecContext(ctx, query, reason, requestID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err 
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil 
}

func (r *PostgresAdminRepository) UpdateUserRole(ctx context.Context, userID uuid.UUID, role string) error {

	query :=
	`
		UPDATE users
		SET role = 'recruiter',
			updated_at = NOW()
		WHERE id = $2	
	`

	result, err := r.db.ExecContext(ctx, query, role, userID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
} 

// func (r *PostgresAdminRepository) UpdateRecruiterRequestStatus(ctx context.Context, requestID uuid.UUID, status string) error {
// 	query := `
// 	UPDATE recruiter_requests
// 	SET status = $1
// 	WHERE id = $2
// 	`

// 	result, err := r.db.ExecContext(ctx, query, requestID, status)
// 	if err != nil {
// 		return err
// 	}

// 	rows, err := result.RowsAffected()
// 	if err != nil {
// 		return err 
// 	}

// 	if rows == 0 {
// 		return ErrNotFound
// 	}

// 	return nil
// }