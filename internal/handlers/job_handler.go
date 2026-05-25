package handlers

import (
	"encoding/json"
	"errors"
	"job_board/internal/domain"
	"job_board/internal/service"
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
	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Company     string `json:"company"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	log.Printf("Received request to create job with body: %v", r.Body)

	// Basic validation (transport-level validation)
	if input.Title == "" || input.Company == "" {
		http.Error(w, "missing required fields", http.StatusBadRequest)
		return
	}
	

	// Create domain job object
	job := &domain.Job{
		Title:       input.Title,
		Description: input.Description,
		Company:     input.Company,
	}

	// Call service layer (business rules happen there)
	err := h.service.CreateJob(r.Context(), job)
	if err != nil {
		h.mapError(w, err)
		return
	}

	// Success response
	log.Printf("Job created successfully with ID %d", job.ID)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(job)
}

// ListJobs handles GET /jobs
func (h *JobHandler) ListJobs(w http.ResponseWriter, r *http.Request) {
	// Parse Query parameters
	limit := 10
	offset := 0

	log.Printf("Received request to list jobs with query params: %v", r.URL.Query())
	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil {
			limit = v
		}
	}

	if o := r.URL.Query().Get("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil {
			offset = v
		}
	}

	jobs, total, err := h.service.ListJobs(r.Context(), limit, offset)
	if err != nil {
		h.mapError(w, err)
		return
	}

	// Structured response
	resp := map[string]interface{}{
		"data": jobs,
		"limit": limit,
		"offset": offset,
		"total": total,
	}

	log.Printf("Listed jobs with limit %d and offset %d, total jobs: %d", limit, offset, total)

	// Success response

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}

func (h *JobHandler) ApplyToJob(w http.ResponseWriter, r *http.Request) {
	jobIDStr := chi.URLParam(r, "id")
	jobID, err := strconv.ParseInt(jobIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid job id", http.StatusBadRequest)
		return
	}

	log.Printf("Received application request for Job ID %d", jobID)
	err = h.service.ApplyToJob(r.Context(), jobID)
	if err != nil {
		// Log the status before returning
		if errors.Is(err, service.ErrAlreadyApplied) {
			http.Error(w, err.Error(), http.StatusConflict)
			log.Printf("User has already applied to Job ID %d", jobID)
			return
		} else {
			log.Printf("Failed to apply to Job ID %d: %v", jobID, err)
		}
		h.mapError(w, err)
		return
	}

	// Succes response
	log.Printf("Job ID %d successfully applied to by the current user", jobID)
	w.WriteHeader(http.StatusCreated)
}

// Health handles GET /health for health checks
func (h *JobHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// mapError maps service errors to HTTP responses.
func (h *JobHandler) mapError(w http.ResponseWriter, err error) {
	switch err {
	case service.ErrInvalidRole:
		http.Error(w, err.Error(), http.StatusForbidden)
	case service.ErrAlreadyApplied:
		http.Error(w, err.Error(), http.StatusConflict)
	case service.ErrUnauthorized:
		http.Error(w, err.Error(), http.StatusUnauthorized)
	case service.ErrForbidden:
		http.Error(w, err.Error(), http.StatusForbidden)
	default:
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
