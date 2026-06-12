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

	ApproveRecruiterRequest(ctx context.Context, requestID uuid.UUID) error

	RejectRecruiterRequest(ctx context.Context, reason string, requestID uuid.UUID) error

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

func (s *adminService) ApproveRecruiterRequest(ctx context.Context, requestID uuid.UUID) error {

	req, err := s.adminRepo.GetRecruiterRequestByUserID(ctx, requestID)
	if err != nil {
		return err
	}

	if req.Status != "pending" {
		return ErrAlreadyAppliedRequest
	}

	err = s.adminRepo.ApproveRecruiterRequest(ctx, requestID)
	if err != nil {
		return err
	}

	return s.adminRepo.UpdateUserRole(ctx, req.RequestID, "recruiter")
}

func (s *adminService) RejectRecruiterRequest(ctx context.Context, reason string, requestID uuid.UUID) error {

	req, err := s.adminRepo.GetRecruiterRequestByUserID(ctx, requestID)
	if err != nil {
		return err
	}

	if req.Status != "pending" {
		return ErrAlreadyAppliedRequest
	}

	return  s.adminRepo.RejectRecruiterRequest(ctx, reason, requestID)
}