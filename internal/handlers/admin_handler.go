package handlers

import (
	"errors"
	"job_board/internal/dto"
	"job_board/internal/service"
	"job_board/internal/validator"
	"job_board/pkg/response"
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AdminHandler interface {
	
	ListRecruiterRequests(w http.ResponseWriter, r *http.Request)

	UpdateRecruiterRequestStatus(w http.ResponseWriter, r *http.Request)
}

type adminHandler struct {
	adminService 	service.AdminService
	logger 			*zap.Logger
}

func NewAdminHandlers(adminService service.AdminService, logger *zap.Logger) AdminHandler {
	return &adminHandler{
		adminService: 	adminService,	
		logger: 		logger,
	}
}

// ListRecruiterRequests handles the HTTP request for admins to list recruiter access requests with pagination.
func (h *adminHandler) ListRecruiterRequests(w http.ResponseWriter, r *http.Request) {
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

	requests, total, err := h.adminService.ListRecruiterRequests(r.Context(), limit, offset)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to retrieve recruiter requests")
		return
	}

	var resp dto.ListRecruiterRequestsResponse
	resp.Requests = make([]dto.RecruiterRequestSummary, len(requests))
	for i, req := range requests {
		resp.Requests[i] = dto.RecruiterRequestSummary{
			ID:             req.ID,
			RecruiterID:    req.RecruiterID,
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

func (h *adminHandler) UpdateRecruiterRequestStatus(w http.ResponseWriter, r *http.Request) {
	// Implementation for updating recruiter request status (approve/reject)
	var req dto.UpdateRecruiterRequestStatusRequest

	// Decode JSON body into DTO
	if err := response.DecodeJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Basic validation (transport-level validation)
	if req.ID == uuid.Nil {
		response.Error(w, http.StatusBadRequest, "invalid recruiter request ID", response.ValidationError{
			Field: "id",
			Error: "ID must be a valid UUID",
		})
		return
	}

	h.logger.Info("Received request to update recruiter request status with body:", zap.Any("request_body", req))

	// Call service layer to update the recruiter request status
	err := h.adminService.UpdateRecruiterRequestStatus(r.Context(), req.ID, req.Status)
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
