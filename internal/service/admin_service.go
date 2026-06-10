package service

import (
	"context"
	"job_board/internal/contextkeys"
	"job_board/internal/domain"
	"job_board/internal/repository"

	"github.com/google/uuid"
)

type AdminService interface {
	
	// ListRecruiterRequests allows admins to view pending recruiter requests.
	ListRecruiterRequests(ctx context.Context, limit, offset int) ([]domain.RecruiterRequest, int64, error)

	UpdateRecruiterRequestStatus(ctx context.Context, recruiterID uuid.UUID, status string) error
}

type adminService struct {
	adminRepo repository.AdminRepository
}

func NewAdminService(adminRepo repository.AdminRepository) AdminService {
	return &adminService{adminRepo: adminRepo}
}

func (s *adminService) ListRecruiterRequests(ctx context.Context, limit, offset int) ([]domain.RecruiterRequest, int64, error) {

	// Only admins should be able to list recruiter requests.
	// This check should ideally be done in the handler layer, but we can also enforce it here for extra safety.
	// Extract user from context
	user, ok := contextkeys.UserFromContext(ctx)
	if !ok {
		return nil, 0, ErrUnauthorized
	}

	if user.Role != "admin" {
		return nil, 0, ErrForbidden
	}
	return s.adminRepo.ListRecruiterRequests(ctx, limit, offset)
}

func (s *adminService) UpdateRecruiterRequestStatus(ctx context.Context, recruiterID uuid.UUID, status string) error {

	// Validate status input
	if status != "approved" && status != "rejected" {
		return ErrInvalidStatus
	}

	// Check if request exists
	req, err := s.adminRepo.GetRecruiterRequestByUserID(ctx, recruiterID)
	if err != nil {
		if err == repository.ErrNotFound {
			return repository.ErrNotFound
		}
		return err
	}

	if req == nil {
		return repository.ErrNotFound
	}

	return s.adminRepo.UpdateRecruiterRequestStatus(ctx, recruiterID, status)
}