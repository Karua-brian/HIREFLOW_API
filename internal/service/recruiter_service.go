package service

import (
	"context"
	"job_board/internal/domain"
	"job_board/internal/repository"

	"github.com/google/uuid"
)

type RecruiterRequestService interface {
	// RequestRecruiterAccess allows a user to request recruiter access.
	RequestRecruiterAccess(ctx context.Context, req *domain.RecruiterRequest) error

	GetMyRecruiterRequest(ctx context.Context, requestID uuid.UUID) (*domain.RecruiterRequest, error)
}


type recruiterRequestService struct {

	recruiterRequestRepo repository.RecruiterRequestRepository
}

func NewRecruiterRequestService(recruiterRequestRepo repository.RecruiterRequestRepository) RecruiterRequestService {

	return &recruiterRequestService{
		recruiterRequestRepo: recruiterRequestRepo,
	}
}

// Recruiter
func (s *recruiterRequestService) RequestRecruiterAccess(ctx context.Context, req *domain.RecruiterRequest) error {

	req.Status = "pending" // default status for new requests

	existingRequest, err := s.recruiterRequestRepo.GetRecruiterRequestByUserID(ctx, req.RequestID)
	if err != nil && err != repository.ErrNotFound {
		return err
	}
	if existingRequest != nil && (existingRequest.Status == "pending" || existingRequest.Status == "approved") {
		return ErrRecruiterRequestAlreadyExists
	}

	return s.recruiterRequestRepo.CreateRecruiterRequest(ctx, req)
}

func (s *recruiterRequestService) GetMyRecruiterRequest(ctx context.Context, requestID uuid.UUID) (*domain.RecruiterRequest, error) {

	request, err := s.recruiterRequestRepo.GetRecruiterRequestByUserID(ctx, requestID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, ErrRecruiterRequestNotFound
		}
		return nil, err
	}

	return request, nil
}


