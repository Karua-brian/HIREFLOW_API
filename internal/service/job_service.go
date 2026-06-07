package service

import (
	"context"
	"job_board/internal/contextkeys"
	"job_board/internal/domain"
	"job_board/internal/repository"

	"github.com/google/uuid"
)

// Service depends on an interface from store
// We can mock this in tests
// JobService defines business logic operations
type JobService interface {

	// CreateJob applies business rules before storing.
	CreateJob(ctx context.Context, job *domain.Job) error

	// ListJobs returns jobs (public endpoint)
	ListJobs(ctx context.Context, limit, offset int) ([]domain.Job, int64, error)

	// ApplyToJob returns jobs applied 
	ApplyToJob(ctx context.Context, jobId uuid.UUID) error
}


// Now implementation
type jobService struct {
	jobRepository repository.JobRepository // Dependency injected
	appRepository repository.ApplicationRepository
	worker JobWorker
}

// Constructor -> injects store dependency
func NewJobService(
	jobRepository repository.JobRepository, 
	appRepository repository.ApplicationRepository, 
	worker   JobWorker,
	) JobService {
	return &jobService{
		jobRepository: jobRepository,
		appRepository: appRepository,
		worker:   worker,
	}
}

// CreateJob implements business rules and delegates to the underlying store if supported.
func (s *jobService) CreateJob(ctx context.Context, job *domain.Job) error {

	// Extract user from context
	// Authenticatioon middleware should have 
	user, ok := contextkeys.UserFromContext(ctx)
	if !ok {
		// No authenticated user -> reject
		return ErrUnauthorized
	}

	// Only recruiters or admins can create jobs.
	if user.Role != "recruiter" && user.Role != "admin" {
		return ErrForbidden
	}

	// Call store to persist
	return s.jobRepository.Create(ctx, job)
}

// ListJobs returns jobs by delegating to the underlying store if supported.
func (s *jobService) ListJobs(ctx context.Context, limit, offset int) ([]domain.Job, int64, error) {
	// No authentication required
	// Public endpoint
	
	// Enfore defaults
	if limit <= 0 {
		limit = 10
	}

	// Prevent abuse
	if limit > 100 {
		limit = 100
	}

	if offset < 0 {
		offset = 0
	}

	jobs, total, err := s.jobRepository.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return jobs, total, nil
}

func (s *jobService) ApplyToJob(ctx context.Context, jobID uuid.UUID) error {
	return s.appRepository.CreateTx(ctx, func(txRepository repository.ApplicationTxRepository) error {

	user, ok := contextkeys.UserFromContext(ctx)
	if !ok {
		return ErrUnauthorized
	}

	// Only applicants can apply
	if user.Role != "applicant" {
		return ErrInvalidRole
	}

	// Check if already applied 
	exists, err := txRepository.Exists(ctx, jobID, user.ID)
	if err != nil {
		return err
	}

	if exists {
		return ErrAlreadyApplied
	}
 
	// Create same application inside the same transaction
	app := &domain.Application{
		JobID: jobID,
		UserID: user.ID,
	}

	if err := txRepository.Create(ctx, app); err != nil {
		return err
	}

	// Push event to worker AFTER successful DB operation
	s.worker.Enqueue(domain.ApplicationEvent{
		JobID: jobID,
		UserID: user.ID,
	})
	
	return nil
	
	})
}	
