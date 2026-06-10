package service

import (
	"context"
	"job_board/internal/domain"
	"job_board/internal/repository"

	"github.com/google/uuid"
)

type RecruiterService interface {
	// RequestRecruiterAccess allows a user to request recruiter access.
	RequestRecruiterAccess(ctx context.Context, req *domain.RecruiterRequest) error

	GetMyRecruiterRequest(ctx context.Context, recruiterID uuid.UUID) (*domain.RecruiterRequest, error)
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

	existingRequest, err := s.recruiterRequestRepo.GetRecruiterRequestByUserID(ctx, req.RecruiterID)
	if err != nil && err != repository.ErrNotFound {
		return err
	}
	if existingRequest != nil && (existingRequest.Status == "pending" || existingRequest.Status == "approved") {
		return ErrRecruiterRequestAlreadyExists
	}

	return s.recruiterRequestRepo.CreateRecruiterRequest(ctx, req)
}

func (s *recruiterService) GetMyRecruiterRequest(ctx context.Context, recruiterID uuid.UUID) (*domain.RecruiterRequest, error) {
	request, err := s.recruiterRequestRepo.GetRecruiterRequestByUserID(ctx, recruiterID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, ErrRecruiterRequestNotFound
		}
		return nil, err
	}

	return request, nil
}


