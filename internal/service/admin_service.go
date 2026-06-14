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
	notificationRepo repository.NotificationRepo
}

func NewAdminService(adminRepo repository.AdminRepository, notificationRepo repository.NotificationRepo) AdminService {
	return &adminService{
		adminRepo: adminRepo,
		notificationRepo: notificationRepo,
	}
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
	
	req, err := s.adminRepo.GetRecruiterRequestByID(ctx, requestID)
	if err != nil {
		return err
	}

	if req.Status != "pending" {
		return ErrAlreadyAppliedRequest
	}

	if req == nil {
		return ErrRecruiterRequestNotFound
	}

	err = s.adminRepo.ApproveRecruiterRequest(ctx, requestID)
	if err != nil {
		return err
	}

	notification := &domain.Notification{
		UserID: req.UserID,
		Type: "approval",
		Title: "Recruiter Request Approved",
		Message: "Congratulations! Your recruiter request has been approved.", 
		Link: "/jobs/create",
	}

	err = s.notificationRepo.CreateNotification(ctx, notification)
	if err != nil {
		return err
	}
	
	return nil
}

func (s *adminService) RejectRecruiterRequest(ctx context.Context, reason string, requestID uuid.UUID) error {

	req, err := s.adminRepo.GetRecruiterRequestByID(ctx, requestID)
	if err != nil {
		return err
	}

	if req == nil {
		return ErrRecruiterRequestNotFound
	}

	if req.Status != "pending" {
		return ErrAlreadyAppliedRequest
	}

	err = s.adminRepo.RejectRecruiterRequest(ctx, reason, requestID)
	if err != nil {
		return err
	}

	notification := &domain.Notification{
		UserID: req.UserID,
		Type: "rejection",
		Title: "Recruiter Request Rejected",
		Message: reason,
		Link: "/recruiter/requests",
	}

	err = s.notificationRepo.CreateNotification(ctx, notification)
	if err != nil {
		return err
	}

	return nil
}
