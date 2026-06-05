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

	"go.uber.org/zap"
)

type RecruiterHandler interface {
	RequestRecruiterAccess(w http.ResponseWriter, r *http.Request)

	ListRecruiterRequests(w http.ResponseWriter, r *http.Request)

	UpdateRecruiterRequestStatus(w http.ResponseWriter, r *http.Request)

	GetMyRecruiterRequest(w http.ResponseWriter, r *http.Request)
}

type recruiterHandler struct {
	service service.RecruiterService
	logger  *zap.Logger
}

func NewRecruiterHandlers(s service.RecruiterService, logger *zap.Logger) RecruiterHandler {
	return &recruiterHandler{
		service: s,
		logger:  logger,
	}
}

// RequestRecruiterAccesss handles the HTTP request for users to request recruiter access.
func (h *recruiterHandler) RequestRecruiterAccess(w http.ResponseWriter, r *http.Request) {
	// Implementation for handling recruiter access requests
	var req dto.CreateRecruiterRequest

	if err := response.DecodeJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Basic validation (transport-level validation)
	if err := validator.ValidateRecruiterRequest(req.CompanyName, req.CompanyWebsite, req.Message); err != nil {
		response.Error(w, http.StatusBadRequest, "validation error", response.ValidationError{
			Field: "company_name/company_website",
			Error: err.Error(),
		})
		return
	}
	h.logger.Info("Received request to create recruiter access with body:", zap.Any("request_body", req))

	user, ok := r.Context().Value(contextkeys.UserKey).(*domain.User) 
	if !ok {
		response.Error(w, http.StatusUnauthorized, "invalid user context")
	}
	userID := user.ID
	// Create domain object for service layer
	request := &domain.RecruiterRequest{
		UserID:         userID,
		CompanyName:    req.CompanyName,
		CompanyWebsite: req.CompanyWebsite,
		Message:        req.Message,
	}

	// Call service layer to process the recruiter access request
	if err := h.service.RequestRecruiterAccess(r.Context(), request); err != nil {
		if errors.Is(err, service.ErrRecruiterRequestAlreadyExists) {
			response.Error(w, http.StatusConflict, "a pending or approved recruiter request already exists for this user")
			return
		}

		response.Error(w, http.StatusInternalServerError, "failed to create request")
		return
	}

	var resp dto.RecruiterResponse
	resp.ID = request.ID
	resp.Status = request.Status
	resp.Message = "Recruiter access request submitted successfully"

	response.JSON(w, http.StatusCreated, resp)
}

// GetMyRecruiterRequest allows users to check the status of their recruiter access request.
func (h *recruiterHandler) GetMyRecruiterRequest(w http.ResponseWriter, r *http.Request) {
	
	user, ok := r.Context().Value(contextkeys.UserKey).(*domain.User) 
	if !ok {
		response.Error(w, http.StatusUnauthorized, "invalid user context")
	}
	userID := user.ID

	request, err := h.service.GetMyRecruiterRequest(r.Context(), userID)
	if err != nil {
		if errors.Is(err, service.ErrRecruiterRequestNotFound) {
			response.Error(w, http.StatusNotFound, "no recruiter request found for this user")
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to retrieve recruiter request status")
		return
	}

	var resp dto.RecruiterResponse
	resp.ID = request.ID
	resp.Status = request.Status
	resp.Message = "Recruiter request status retrieved successfully"

	response.JSON(w, http.StatusOK, resp)
}

// ListRecruiterRequests handles the HTTP request for admins to list recruiter access requests with pagination.
func (h *recruiterHandler) ListRecruiterRequests(w http.ResponseWriter, r *http.Request) {
	// Implementation for listing recruiter access requests
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit, offset, err := validator.ParsePaginationParams(limitStr, offsetStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid pagination parameters", response.ValidationError{
			Field: "limit/offset",
			Error: err.Error(),
		})
		return
	}

	h.logger.Info("Received request to list recruiter requests with pagination", zap.Int("limit", limit), zap.Int("offset", offset))

	requests, total, err := h.service.ListRecruiterRequests(r.Context(), limit, offset)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to retrieve recruiter requests")
		return
	}

	var resp dto.ListRecruiterRequestsResponse
	resp.Requests = make([]dto.RecruiterRequestSummary, len(requests))
	for i, req := range requests {
		resp.Requests[i] = dto.RecruiterRequestSummary{
			ID:             req.ID,
			UserID:         req.UserID,
			CompanyName:    req.CompanyName,
			CompanyWebsite: req.CompanyWebsite,
			Message:        req.Message,
			Status:         req.Status,
		}
	}
	resp.Total = total
	resp.Limit = limit
	resp.Offset = offset

	h.logger.Info("Listed recruiter requests with pagination, total requests:", zap.Int64("total_requests", total))

	response.JSON(w, http.StatusOK, resp)
}

func (h *recruiterHandler) UpdateRecruiterRequestStatus(w http.ResponseWriter, r *http.Request) {
	// Implementation for updating recruiter request status (approve/reject)
	var req dto.UpdateRecruiterRequestStatusRequest

	// Decode JSON body into DTO
	if err := response.DecodeJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Basic validation (transport-level validation)
	if req.ID <= 0 {
		response.Error(w, http.StatusBadRequest, "invalid recruiter request ID", response.ValidationError{
			Field: "id",
			Error: "ID must be a positive integer",
		})
		return
	}

	h.logger.Info("Received request to update recruiter request status with body:", zap.Any("request_body", req))

	// Call service layer to update the recruiter request status
	err := h.service.UpdateRecruiterRequestStatus(r.Context(), req.ID, req.Status)
	if err != nil {
		if errors.Is(err, service.ErrRecruiterRequestNotFound) {
			response.Error(w, http.StatusNotFound, "recruiter request not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to update recruiter request status")
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{
		"message": "Recruiter request status updated successfully",
	})
}