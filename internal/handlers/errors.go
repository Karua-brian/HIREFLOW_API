package handlers

import (
	"job_board/internal/service"
	"job_board/pkg/response"
	"net/http"
)

// mapError maps service errors to HTTP responses.
func (h *jobHandler) mapError(w http.ResponseWriter, err error) {
	switch err {
	case service.ErrInvalidRole:
		response.Error(w, http.StatusForbidden, "only job seekers can apply to jobs")
	case service.ErrAlreadyApplied:
		response.Error(w, http.StatusConflict, "user has already applied to this job")
	case service.ErrUnauthorized:
		response.Error(w, http.StatusUnauthorized, "unauthorized")
	case service.ErrForbidden:
		response.Error(w, http.StatusForbidden, "forbidden")		
	default:
		response.Error(w, http.StatusInternalServerError, "internal server error")
	}
}