package handlers

import (
	"job_board/internal/dto"
	"job_board/internal/service"
	"job_board/internal/validator"
	"job_board/pkg/response"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AdminHandler interface {
	ListRecruiterRequests(w http.ResponseWriter, r *http.Request)

	ApproveRecruiterRequest(w http.ResponseWriter, r *http.Request)

	RejectRecruiterRequest(w http.ResponseWriter, r *http.Request)
}

type adminHandler struct {
	adminService service.AdminService
	logger       *zap.Logger
}

func NewAdminHandlers(adminService service.AdminService, logger *zap.Logger) AdminHandler {
	return &adminHandler{
		adminService: adminService,
		logger:       logger,
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
			ID:          req.ID,
			RequestID: 	 req.RequestID,
			CompanyName: req.CompanyName,
			Message:     req.Message,
			Status:      req.Status,
		}
	}
	resp.Total = total
	resp.Limit = limit
	resp.Offset = offset

	h.logger.Info("Listed recruiter requests with pagination, total requests:", zap.Int64("total_requests", total))

	response.JSON(w, http.StatusOK, resp)
}

func (h *adminHandler) ApproveRecruiterRequest(w http.ResponseWriter, r *http.Request) {

	requestIDStr := chi.URLParam(r, "id")

	requestID, err := uuid.Parse(requestIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid recruiter request id")
		return
	}

	err = h.adminService.ApproveRecruiterRequest(r.Context(), requestID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to approve recruiter request")
		h.logger.Info("failed to approve recruiter request", zap.Error(err), zap.String("request_id", requestID.String()))
		return
	}

	h.logger.Info("recruiter request approved", zap.String("request_id", requestID.String()))

	response.JSON(w, http.StatusOK, map[string]string{
		"message": "recruiter request approved successfully",
	})
}

func (h *adminHandler) RejectRecruiterRequest(w http.ResponseWriter, r *http.Request) {

	requestIDStr := chi.URLParam(r, "id")

	requestID, err := uuid.Parse(requestIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid recruiter request id")
		return
	}

	var req dto.RejectRecruiterRequest

	// Use the response helper to decode JSON body and handle errors
	if err := response.DecodeJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if strings.TrimSpace(req.Reason) == "" {
		response.Error(w, http.StatusBadRequest, "rejection reason is required")
		return
	}

	err = h.adminService.RejectRecruiterRequest(r.Context(), req.Reason, requestID)
}
