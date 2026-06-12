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

type RecruiterRequestHandler interface {
	RequestRecruiterAccess(w http.ResponseWriter, r *http.Request)

	GetMyRecruiterRequest(w http.ResponseWriter, r *http.Request)
}

type recruiterRequestHandler struct {
	service service.RecruiterRequestService
	logger  *zap.Logger
}

func NewRecruiterRequestHandlers(s service.RecruiterRequestService, logger *zap.Logger) RecruiterRequestHandler {
	return &recruiterRequestHandler{
		service: s,
		logger:  logger,
	}
}

// RequestRecruiterAccesss handles the HTTP request for users to request recruiter access.
func (h *recruiterRequestHandler) RequestRecruiterAccess(w http.ResponseWriter, r *http.Request) {
	// Implementation for handling recruiter access requests
	var req dto.CreateRecruiterRequest

	if err := response.DecodeJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Basic validation (transport-level validation)
	if err := validator.ValidateRecruiterRequest(req.CompanyName, req.CompanyWebsite, req.Message); err != nil {
		response.Error(w, http.StatusBadRequest, "validation error", response.ValidationError{
			Field: "company_name/company_website/message",
			Error: err.Error(),
		})
		return
	}
	h.logger.Info("Received request to create recruiter access with body:", zap.Any("request_body", req))

	user, ok := contextkeys.UserFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "invalid user context")
	}

	userID := user.ID
	// Create domain object for service layer
	request := &domain.RecruiterRequest{
		RequestID:		userID,
		CompanyName:    req.CompanyName,
		CompanyWebsite: req.CompanyWebsite,
		Message:        req.Message,
	}

	// Call service layer to process the recruiter access request
	err := h.service.RequestRecruiterAccess(r.Context(), request)
	if err != nil {
		h.logger.Error("failed to create recruiter request",
			zap.Error(err),
			zap.Any("request", request),
		)

		if errors.Is(err, service.ErrRecruiterRequestAlreadyExists) {
			response.Error(w, http.StatusConflict, "already exists")
			return
		}

		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	var resp dto.RecruiterRequestResponse

	resp.RequestID = request.RequestID
	resp.Status = request.Status
	resp.Message = "Recruiter access request submitted successfully"

	response.JSON(w, http.StatusCreated, resp)
}

// GetMyRecruiterRequest allows users to check the status of their recruiter access request.
func (h *recruiterRequestHandler) GetMyRecruiterRequest(w http.ResponseWriter, r *http.Request) {

	user, ok := contextkeys.UserFromContext(r.Context())
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

	var resp dto.RecruiterRequestResponse
	resp.RequestID = request.ID
	resp.Status = request.Status
	resp.Message = "Recruiter request status retrieved successfully"

	response.JSON(w, http.StatusOK, resp)
}

