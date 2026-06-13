package handlers

import (
	"errors"
	"job_board/internal/dto"
	"job_board/internal/service"
	"job_board/internal/validator"
	"job_board/pkg/response"
	"net/http"
	"strings"

	"go.uber.org/zap"
)
type AuthHandler interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Refresh(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
}

// AuthHandler holds dependencies for authentication related HTTP handlers.
type authHandler struct {
	authService service.AuthService
	logger *zap.Logger
}

// Constructor - dependecy injection
func NewAuthHandlers(s service.AuthService, logger *zap.Logger) AuthHandler {
	return &authHandler{
		authService: s,
		logger: logger,
	}
}

// Register handles POST /register
func (h *authHandler) Register(w http.ResponseWriter, r *http.Request) {
	
	// import the Register request struct from the dto package
	var req dto.RegisterRequest

	// Use the response helper to decode JSON body and handle errors
	if err := response.DecodeJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	
	// Basic validation (transport-level validation)
	if err := validator.ValidateRegister(
		req.Email,
		req.Password,
	); err != nil {
		response.Error(w, http.StatusBadRequest, "validation error", response.ValidationError{
			Field: "email/password",
			Error: err.Error(),
		})
		return	
	}
	
	// Call service layer to register the user and handle specific errors
	err := h.authService.Register(
		r.Context(),
		req.Email,
		req.Password,
	)
	h.logger.Info("Registration attempt for email:", zap.String("email", req.Email))

	// Handle specific service errors and return appropriate HTTP responses
	if err != nil {
		if errors.Is(err, service.ErrUserEmailExists) {
			h.mapError(w, err) // Map to 409 Conflict
			h.logger.Info("Registration failed for email:", zap.String("email", req.Email), zap.Error(err))
			return
		}
	}

	// Return success response
	response.JSON(w, http.StatusCreated, map[string]string{
		"message": "user registered successfully", 
		})	
	h.logger.Info("User with email successfully registered", zap.String("email", req.Email))	
}

// Login handles POST /login
func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {

	// import the Login request struct from the dto package
	var req dto.LoginRequest 

	// Use the response helper to decode JSON body and handle errors
	if err := response.DecodeJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Basic validation (transport-level validation)
	if err := validator.ValidateLogin(
		req.Email,
		req.Password,
	); err != nil {
		response.Error(w, http.StatusBadRequest, "validation error", response.ValidationError{
			Field: "email/password",
			Error: err.Error(),
		})
		return	
	}
	h.logger.Info("Login attempt for email:", zap.String("email", req.Email))

	// Call service layer to login the user and handle specific errors
	token, refresh, user, err := h.authService.Login(
		r.Context(),
		req.Email,
		req.Password,
	)

	// Handle specific service errors and return appropriate HTTP responses
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			h.mapError(w, err) // Map to 401 Unauthorized
			h.logger.Info("Login failed for email: invalid credentials", zap.String("email", req.Email))
			return
		}
	}

	// Return the token in the response
	resp := dto.LoginResponse{
		AccessToken: token,
		RefreshToken: refresh,
		User: dto.UserDTO{
			ID: user.ID.String(),
			Email: user.Email,
			Role: user.Role,
		},
	}

	// Return success response with token
	h.logger.Info("User with email successfully logged in", zap.String("email", req.Email))

	// Use the response helper to send a JSON response
	response.JSON(w, http.StatusOK, resp)
}

// Refresh handles POST /refresh-token
func (h *authHandler) Refresh(w http.ResponseWriter, r *http.Request) {

	authHeader := r.Header.Get("Authorization")

	// Validate the header format
	if err := validator.ValidateJWTHeader(authHeader); err != nil {
		response.Error(w, http.StatusUnauthorized, "Missing or invalid Authorization header")
		return
	}

	// Expected format: "Bearer <token>"
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		response.Error(w, http.StatusUnauthorized, "Invalid Authorization header format")
		return 
	}

	// Extract the token string (remove "Bearer " prefix)
	refreshToken := parts[1]
	
	h.logger.Info("Refresh token attempt")
	// Call service layer to refresh the token and handle specific errors
	access, refresh, err := h.authService.Refresh(r.Context(), refreshToken)
	if err != nil {
		if errors.Is(err, service.ErrInvalidRefreshToken) {
			h.mapError(w, err) // Map to 401 Unauthorized
			h.logger.Info("Refresh token failed: invalid refresh token")
			return
		}
	}

	// Return the user ID associated with the refresh token
	resp := map[string]string{
		"access_token": access,
		"refresh_token": refresh,
	}
	h.logger.Info("Refresh token successful, new access token issued")

	// Return success response with new access token
	response.JSON(w, http.StatusOK, resp)	
}

// 
func (h *authHandler) Logout(w http.ResponseWriter, r *http.Request) {
	
	authHeader := r.Header.Get("Authorization")

	// Validate the header format
	if err := validator.ValidateJWTHeader(authHeader); err != nil {
		response.Error(w, http.StatusUnauthorized, "Missing or invalid Authorization header")
		return
	}

	// Expected format: "Bearer <token>"
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		response.Error(w, http.StatusUnauthorized, "Invalid Authorization header format")
		return 
	}

	// Extract the token string (remove "Bearer " prefix)
	tokenString := parts[1]

	h.logger.Info("Logout attempt with token")

	err := h.authService.Logout(r.Context(), tokenString)
	if err != nil {
		if errors.Is(err, service.ErrInvalidRefreshToken) {
			h.mapError(w, err) // Map to 401 Unauthorized
			h.logger.Info("Logout failed: invalid refresh token")
			return
		}
	}

	h.logger.Info("Logout successful, refresh token invalidated")
	// Return success response
	response.JSON(w, http.StatusOK, map[string]string{
		"message": "logged out successfully",
	})
}

func (h *authHandler) mapError(w http.ResponseWriter, err error) {
	switch err {
	case service.ErrUserEmailExists:
		response.Error(w, http.StatusConflict, "email exists")
	case service.ErrInvalidCredentials:
		response.Error(w, http.StatusUnauthorized, "invalid email or password")
	case service.ErrInvalidRefreshToken:
		response.Error(w, http.StatusUnauthorized, "invalid refresh token")	
	default:
		response.Error(w, http.StatusInternalServerError, "internal server error")
	}
}