package handlers

import (
	"errors"
	"job_board/internal/contextkeys"
	"job_board/internal/domain"
	"job_board/internal/dto"
	"job_board/internal/service"
	"job_board/internal/validator"
	"job_board/pkg/response"
	"net/http"

	"github.com/google/uuid"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type JobHandler interface {
	CreateJob(w http.ResponseWriter, r *http.Request)
	ListJobs(w http.ResponseWriter, r *http.Request)
	ApplyToJob(w http.ResponseWriter, r *http.Request)
	Health(w http.ResponseWriter, r *http.Request)
}

// JobHandler holds dependencies for job related HTTP handlers.
type jobHandler struct {
	service service.JobService
	logger  *zap.Logger
}

// Constructor - dependency injection
func NewJobHandlers(s service.JobService, logger *zap.Logger) JobHandler {
	return &jobHandler{
		service: s,
		logger:  logger,
	}
}

// CreateJob handles POST /jobs
func (h *jobHandler) CreateJob(w http.ResponseWriter, r *http.Request) {

	// Decode request body into domain.Job
	// We only accept Title, Description, Company from client

	user, ok := contextkeys.UserFromContext(r.Context())
		if !ok {
			response.Error(w, http.StatusUnauthorized, "invalid user context")
			return 
		}
		userID := user.ID

	var req dto.CreateJobRequest

	if err := response.DecodeJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Basic validation (transport-level validation)
	if err := validator.ValidateCreateJob(req.Title, req.Company, req.Location, req.Salary); err != nil {
		response.Error(w, http.StatusBadRequest, "validation error", response.ValidationError{
			Field: "title/company",
			Error: err.Error(),
		})
		return
	}
	h.logger.Info("Received request to create job with body:", zap.Any("request_body", req))

	// Create domain job object
	job := &domain.Job{
		RecruiterUserID:  userID,
		Title:       req.Title,
		Description: req.Description,
		Company:     req.Company,
		Location:    req.Location,
		Salary:      req.Salary,
		// CreatedBy will be set in the service layer based on the authenticated user
	}

	// Call service layer (business rules happen there)
	err := h.service.CreateJob(r.Context(), job, userID)
	if err != nil {
		if errors.Is(err, service.ErrUnauthorized) {
			h.mapError(w, err) // Map to 401 unauthorized
			h.logger.Info("unauthorized attempt to create job:", zap.Error(err))
			return
		}
	}

	if err != nil {
    	h.logger.Error("create job failed", zap.Error(err))
    	response.Error(w, http.StatusInternalServerError, err.Error())
    	return 
	}

	// Success response
	h.logger.Info("Job created successfully with ID", zap.Any("job_id", job.ID))
	response.JSON(w, http.StatusOK, job)
}

// ListJobs handles GET /jobs
func (h *jobHandler) ListJobs(w http.ResponseWriter, r *http.Request) {

	// Parse query params for pagination
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	// Basic validation (transport-level validation)
	limit, offset, err := validator.ParsePaginationParams(limitStr, offsetStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid pagination parameters", response.ValidationError{
			Field: "limit/offset",
			Error: err.Error(),
		})
		return
	}

	h.logger.Info("Listing jobs with limit and offset:", zap.Int("limit", limit), zap.Int("offset", offset))

	// Call service layer to get jobs and total count
	jobs, total, err := h.service.ListJobs(r.Context(), limit, offset)
	if err != nil {
		h.logger.Info("Failed to list jobs:", zap.Error(err))
		response.Error(w, http.StatusInternalServerError, "failed to list jobs")
		return
	}

	// Structured response
	var resp dto.ListJobsResponse

	// Map domain jobs to DTO job summaries
	resp.Jobs = make([]dto.JobSummary, len(jobs))
	for i, job := range jobs {
		resp.Jobs[i] = dto.JobSummary{
			ID:          job.ID,
			RecruiterUserID:  job.RecruiterUserID,
			Title:       job.Title,
			Description: job.Description,
			Company:     job.Company,
			Location:    job.Location,
			Salary:      job.Salary,
		}
	}
	// Include pagination metadata in the response
	resp.Limit = limit
	resp.Offset = offset
	resp.Total = total

	h.logger.Info("Listed jobs with limit and offset, total jobs:", zap.Int64("total_jobs", total))

	// Success response
	response.JSON(w, http.StatusOK, resp)

}

func (h *jobHandler) ApplyToJob(w http.ResponseWriter, r *http.Request) {

	// Extract job ID from URL path
	jobIDStr := chi.URLParam(r, "id")
	jobID, err := uuid.Parse(jobIDStr)

	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid job id")
		return
	}

	// Basic validation (transport-level validation)
	if err := validator.ValidateApplyJob(jobID); err != nil {
		response.Error(w, http.StatusBadRequest, "validation error", response.ValidationError{
			Field: "job_id/role",
			Error: err.Error(),
		})
		return
	}
	h.logger.Info("Received request to apply to Job ID", zap.Any("job_id", jobID))

	// Call service layer to apply to the job
	err = h.service.ApplyToJob(r.Context(), jobID)
	if err != nil {
		// Log the status before returning
		if errors.Is(err, service.ErrAlreadyApplied) {
			h.mapError(w, err) // Map to 409 Conflict
			h.logger.Info("User has already applied to Job ID:", zap.Any("job_id", jobID))
			return
		} else if errors.Is(err, service.ErrInvalidRole) {
			h.mapError(w, err) // Map to 403 Forbidden
			h.logger.Info("User with invalid role attempted to apply to Job ID", zap.Any("job_id", jobID))
			return
		} else if errors.Is(err, service.ErrUnauthorized) {
			h.mapError(w, err) // Map to 401 Unauthorized
			h.logger.Info("Unauthorized user attempted to apply to Job ID", zap.Any("job_id", jobID))
			return
		} else if errors.Is(err, service.ErrForbidden) {
			h.mapError(w, err) // Map to 403 Forbidden
			h.logger.Info("Forbidden action: user attempted to apply to Job ID", zap.Any("job_id", jobID))
			return
		}

		h.logger.Info("Failed to apply to Job ID:", zap.Any("job_id", jobID), zap.Error(err))
		response.Error(w, http.StatusInternalServerError, "failed to apply to job")
		return
	}

	// Succes response
	h.logger.Info("Job successfully applied to by the current user", zap.Any("job_id", jobID))
	response.JSON(w, http.StatusCreated, map[string]string{
		"message": "application successful",
	})
}

// Health handles GET /health for health checks
func (h *jobHandler) Health(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}
