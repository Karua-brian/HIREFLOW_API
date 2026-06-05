package service

import (
	"context"
	"job_board/internal/contextkeys"
	"job_board/internal/domain"
	"job_board/internal/repository"
)

type RecruiterService interface {
	// RequestRecruiterAccess allows a user to request recruiter access.
	RequestRecruiterAccess(ctx context.Context, req *domain.RecruiterRequest) error

	// ListRecruiterRequests allows admins to view pending recruiter requests.
	ListRecruiterRequests(ctx context.Context, limit, offset int) ([]domain.RecruiterRequest, int64, error)

	UpdateRecruiterRequestStatus(ctx context.Context, id int64, status string) error

	GetMyRecruiterRequest(ctx context.Context, id int64) (*domain.RecruiterRequest, error)
}


type recruiterService struct {
	recruiterRequestRepo repository.RecruiterRequestRepository
}

func NewRecruiterService(recruiterRequestRepo repository.RecruiterRequestRepository) RecruiterService {
	return &recruiterService{
		recruiterRequestRepo: recruiterRequestRepo,
	}
}

// Recruiter
func (s *recruiterService) RequestRecruiterAccess(ctx context.Context, req *domain.RecruiterRequest) error {
	req.Status = "pending" // default status for new requests

	existingRequest, err := s.recruiterRequestRepo.GetRecruiterRequestByUserID(ctx, req.UserID)
	if err != nil && err != repository.ErrNotFound {
		return err
	}
	if existingRequest != nil && (existingRequest.Status == "pending" || existingRequest.Status == "approved") {
		return ErrRecruiterRequestAlreadyExists
	}

	return s.recruiterRequestRepo.CreateRecruiterRequest(ctx, req)
}

func (s *recruiterService) GetMyRecruiterRequest(ctx context.Context, userID int64) (*domain.RecruiterRequest, error) {
	request, err := s.recruiterRequestRepo.GetRecruiterRequestByUserID(ctx, userID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, ErrRecruiterRequestNotFound
		}
		return nil, err
	}

	return request, nil
}

func (s *recruiterService) ListRecruiterRequests(ctx context.Context, limit, offset int) ([]domain.RecruiterRequest, int64, error) {

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
	return s.recruiterRequestRepo.ListRecruiterRequests(ctx, limit, offset)
}

func (s *recruiterService) UpdateRecruiterRequestStatus(ctx context.Context, id int64, status string) error {

	// Validate status input
	if status != "approved" && status != "rejected" {
		return ErrInvalidStatus
	}

	// Check if request exists
	req, err := s.recruiterRequestRepo.GetRecruiterRequestByUserID(ctx, id)
	if err != nil {
		if err == repository.ErrNotFound {
			return repository.ErrNotFound
		}
		return err
	}

	if req == nil {
		return repository.ErrNotFound
	}

	return s.recruiterRequestRepo.UpdateRecruiterRequestStatus(ctx, id, status)
}
