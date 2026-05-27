package handlers

import (
	"errors"
	"job_board/internal/domain"
	"job_board/internal/dto"
	"job_board/internal/service"
	"job_board/internal/validator"
	"job_board/pkg/response"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// JobHandler holds dependencies for job related HTTP handlers.
type JobHandler struct {
	service service.JobService
}

// Constructor - dependency injection
func NewJobHandlers(s service.JobService) *JobHandler {
	return &JobHandler{service: s}
}


// CreateJob handles POST /jobs
func (h *JobHandler) CreateJob(w http.ResponseWriter, r *http.Request) {

	// Decode request body into domain.Job
	// We only accept Title, Description, Company from client
	var req dto.CreateJobRequest
	
	if err := response.DecodeJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Basic validation (transport-level validation)
	if err := validator.ValidateCreateJob(req.Title, req.Company); err != nil {
		response.Error(w, http.StatusBadRequest, "validation error", response.ValidationError{
			Field: "title/company",
			Error: err.Error(),
		})
		return
	}
	log.Printf("Received request to create job with body: %v", r.Body)
	

	// Create domain job object
	job := &domain.Job{
		Title:       req.Title,
		Description: req.Description,
		Company:     req.Company,
	}

	// Call service layer (business rules happen there)
	err := h.service.CreateJob(r.Context(), job)

	if err != nil {
		if errors.Is(err, service.ErrUnauthorized) {
			h.mapError(w, err) // Map to 401 unauthorized
			log.Printf("unauthorized attempt to create job: %v", err)
			return
		}
	}

	// Success response
	log.Printf("Job created successfully with ID %d", job.ID)
	response.JSON(w, http.StatusOK, job)	
}

// ListJobs handles GET /jobs
func (h *JobHandler) ListJobs(w http.ResponseWriter, r *http.Request) {

	// Parse query params for pagination
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	// Set default values if not provided
	limit := 10
	offset := 0

	// Basic validation (transport-level validation)
	if err := validator.ValidateListJobs(limit, offset); err != nil {
		response.Error(w, http.StatusBadRequest, "validation error", response.ValidationError{
			Field: "limit/offset",
			Error: err.Error(),
		})
		return
	}

	// Parse limit and offset if provided
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil || parsedLimit < 0 {
			response.Error(w, http.StatusBadRequest, "invalid limit parameter")
			return
		}
		limit = parsedLimit
	}

	if offsetStr != "" {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err != nil || parsedOffset < 0 {
			response.Error(w, http.StatusBadRequest, "invalid offset parameter")
			return
		}
		offset = parsedOffset
	}

	log.Printf("Listing jobs with limit %d and offset %d", limit, offset)

	// Call service layer to get jobs and total count
	jobs, total, err := h.service.ListJobs(r.Context(), limit, offset)
	if err != nil {
		log.Printf("Failed to list jobs: %v", err)
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
			Title:       job.Title,
			Description: job.Description,
			Company:     job.Company,
		}
	}
	// Include pagination metadata in the response
	resp.Limit = limit
	resp.Offset = offset
	resp.Total = total

	log.Printf("Listed jobs with limit %d and offset %d, total jobs: %d", limit, offset, total)

	// Success response
	response.JSON(w, http.StatusOK, resp)

}

func (h *JobHandler) ApplyToJob(w http.ResponseWriter, r *http.Request) {

	// Extract job ID from URL path
	jobIDStr := chi.URLParam(r, "id")
	jobID, err := strconv.ParseInt(jobIDStr, 10, 64)

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
	log.Printf("Received request to apply to Job ID %d", jobID)

	// Call service layer to apply to the job
	err = h.service.ApplyToJob(r.Context(), jobID)
	if err != nil {
		// Log the status before returning
		if errors.Is(err, service.ErrAlreadyApplied) {
			h.mapError(w, err) // Map to 409 Conflict
			log.Printf("User has already applied to Job ID %d", jobID)
			return
		} else if errors.Is(err, service.ErrInvalidRole) {
			h.mapError(w, err) // Map to 403 Forbidden
			log.Printf("User with invalid role attempted to apply to Job ID %d", jobID)
			return
		} else if errors.Is(err, service.ErrUnauthorized) {
			h.mapError(w, err) // Map to 401 Unauthorized
			log.Printf("Unauthorized user attempted to apply to Job ID %d", jobID)
			return
		} else if errors.Is(err, service.ErrForbidden) {
			h.mapError(w, err) // Map to 403 Forbidden
			log.Printf("Forbidden action: user attempted to apply to Job ID %d", jobID)
			return
		 }

		log.Printf("Failed to apply to Job ID %d: %v", jobID, err)
		response.Error(w, http.StatusInternalServerError, "failed to apply to job")
		return
	}

	// Succes response
	log.Printf("Job ID %d successfully applied to by the current user", jobID)
	response.JSON(w, http.StatusCreated, map[string]string{
		"message": "application successful",
	})
}

// Health handles GET /health for health checks
func (h *JobHandler) Health(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

// mapError maps service errors to HTTP responses.
func (h *JobHandler) mapError(w http.ResponseWriter, err error) {
	switch err {
	case service.ErrInvalidRole:
		response.Error(w, http.StatusForbidden, "only job seekers can apply to jobs")
	case service.ErrAlreadyApplied:
		response.Error(w, http.StatusConflict, "user has already applied to this job")
	case service.ErrUnauthorized:
		response.Error(w, http.StatusUnauthorized, "unauthorized")
	default:
		response.Error(w, http.StatusInternalServerError, "internal server error")
	}
}
