package service

import (
	"context"
	"errors"
	"job_board/domain"
	"job_board/middleware"
	"job_board/store"
	"sync"
	"testing"
)

// mockJobStore implemnts store.JobStore for testing purposes
// It allows us to test JobService without a database
type mockJobStore struct {
	createFunc func(ctx context.Context, job *domain.Job) error
	listFunc   func(ctx context.Context, limit, offset int) ([]domain.Job, int64, error)
}

type mockWorker struct {
	called bool
	event  domain.ApplicationEvent
}

type mockApplicationStore struct {
	createFunc func(ctx context.Context, app *domain.Application) error
	existsFunc func(ctx context.Context, jobID, userID int64) (bool, error)
}

// CreateTx implements [store.ApplicationStore].

func (m *mockApplicationStore) Create(ctx context.Context, app *domain.Application) error {
	return m.createFunc(ctx, app)
}

func (m *mockApplicationStore) Exists(ctx context.Context, jobID, userID int64) (bool, error) {
	return m.existsFunc(ctx, jobID, userID)
}

func (m *mockJobStore) Create(ctx context.Context, job *domain.Job) error {
	return m.createFunc(ctx, job)
}

func (m *mockJobStore) List(ctx context.Context, limit, offset int) ([]domain.Job, int64, error) {
	return m.listFunc(ctx, limit, offset)
}

func (m *mockWorker) Enqueue(event domain.ApplicationEvent) {
	m.called = true
	m.event = event
}

func (m *mockApplicationStore) CreateTx(ctx context.Context, fn func(store.ApplicationTxStore) error) error {
	return fn(m)
}

func TestCreateJob_Unauthorized(t *testing.T) {

	// Create a mock store - will not be called in this test
	mockStore := &mockJobStore{}

	// Create the service
	svc := NewJobService(mockStore, &mockApplicationStore{}, &mockWorker{},)

	// Context without user -> simulated missing authentication
	ctx := context.Background()

	job := &domain.Job{
		Title:       "Backend Engineer",
		Description: "Build APIs",
		Company:     "Acme Corp",
	}

	// Call CreateJob
	err := svc.CreateJob(ctx, job)

	// Expected ErrUnauthorized
	if !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
}

func TestCreateJob_Forbidden(t *testing.T) {

	// Create a mock store -> will not be called in this test
	mockStore := &mockJobStore{}

	// Create the service
	svc := NewJobService(mockStore, &mockApplicationStore{}, &mockWorker{},)

	// Context with a user who is not recruiter/admin
	ctx := middleware.WithUser(context.Background(), &domain.User{
		ID:   2,
		Role: "applicant", // not allowed
	})

	job := &domain.Job{
		Title:       "Backend Engineer",
		Description: "Build APIs",
		Company:     "Acme Corp",
	}

	// Call CreateJob
	err := svc.CreateJob(ctx, job)

	// Expected ErrUnauthorized
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestCreateJob_Success(t *testing.T) {

	// Track if store.Create was called
	called := false

	mockStore := &mockJobStore{
		createFunc: func(ctx context.Context, job *domain.Job) error {
			called = true

			// Ensure createdBy is set
			if job.CreatedBy != 1 {
				t.Fatalf("expected createdBy=1, got %d", job.CreatedBy)
			}
			return nil
		},
	}

	svc := NewJobService(mockStore, &mockApplicationStore{}, &mockWorker{},)

	// Context with recruiter/admin user
	ctx := middleware.WithUser(context.Background(), &domain.User{
		ID:   1,
		Role: "recruiter",
	})

	job := &domain.Job{
		Title:       "Backend Engineer",
		Description: "Build APIs",
		Company:     "Acme Corp",
	}

	err := svc.CreateJob(ctx, job)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !called {
		t.Fatalf("expected store.Create to be called")
	}
}
func TestApplyToJob_Concurrent(t *testing.T) {
	var mu sync.Mutex
	created := false

	mockAppStore := &mockApplicationStore{
		existsFunc: func(ctx context.Context, jobID, userID int64) (bool, error) {
			mu.Lock()
			defer mu.Unlock()
			return created, nil
		},
		createFunc: func(ctx context.Context, app *domain.Application) error {
			mu.Lock()
			defer mu.Unlock()

			if created {
				return errors.New("duplicate key value violates unique constraint")
			}

			created = true
			return nil
		},
	}

	worker := NewWorker(100, 4)
	worker.Start()

	svc := NewJobService(&mockJobStore{}, mockAppStore, worker,)

	ctx := middleware.WithUser(context.Background(), &domain.User{
		ID:   1,
		Role: "applicant",
	})

	jobID := int64(1)

	var wg sync.WaitGroup
	attempts := 20
	success := 0

	errCh := make(chan error, attempts)

	for i := 0; i < attempts; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			err := svc.ApplyToJob(ctx, jobID)
			if err == nil {
				mu.Lock()
				success++
				mu.Unlock()
			} else if !errors.Is(err, ErrAlreadyApplied) {
				// send unexpected error to channel to be handled by the main goroutine
				errCh <- err
			}
		}()
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	if success != 1 {
		t.Fatalf("expected 1 success, got %d", success)
	}
}
