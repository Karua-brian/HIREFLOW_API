package handlers

import (
	"errors"
	"job_board/internal/dto"
	"job_board/internal/service"
	"job_board/internal/validator"
	"job_board/pkg/response"
	"log"
	"net/http"
)

// AuthHandler holds dependencies for authentication related HTTP handlers.
type AuthHandler struct {
	authService service.AuthService
}

// Constructor - dependecy injection
func NewAuthHandlers(s service.AuthService) *AuthHandler {
	return &AuthHandler{authService: s}
}

// Register handles POST /register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	
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
		req.Role,
	); err != nil {
		response.Error(w, http.StatusBadRequest, "validation error", response.ValidationError{
			Field: "email/password/role",
			Error: err.Error(),
		})
		return	
	}
	
	// Call service layer to register the user and handle specific errors
	err := h.authService.Register(
		r.Context(),
		req.Email,
		req.Password,
		req.Role,
	)
	log.Printf("Registration attempt for email: %s", req.Email)

	// Handle specific service errors and return appropriate HTTP responses
	if err != nil {
		if errors.Is(err, service.ErrUserEmailExists) {
			h.mapError(w, err) // Map to 409 Conflict
			log.Printf("Registration failed for email %s: email already exists", req.Email)
			return
		}
	}

	// Return success response
	response.JSON(w, http.StatusCreated, map[string]string{
		"message": "user registered successfully", 
		})	
	log.Printf("User with email %s successfully registered", req.Email)	
}

// Login handles POST /login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

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
	log.Printf("Login attempt for email: %s", req.Email)

	// Call service layer to login the user and handle specific errors
	token, refresh, err := h.authService.Login(
		r.Context(),
		req.Email,
		req.Password,
	)

	// Handle specific service errors and return appropriate HTTP responses
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			h.mapError(w, err) // Map to 401 Unauthorized
			log.Printf("Login failed for email %s: invalid credentials", req.Email)
			return
		}
	}

	// Return the token in the response
	resp := map[string]string{
		"access_token" : token,
		"refresh_token": refresh,
		
	}

	// Return success response with token
	log.Printf("User with email %s successfully logged in", req.Email)

	// Use the response helper to send a JSON response
	response.JSON(w, http.StatusOK, resp)
}

// Refresh handles POST /refresh-token
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {

	// Decode req body into struct
	var req dto.RefreshTokenRequest

	if err := response.DecodeJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest,"invalid request body")
		return
	}

	// Basic validation (transport-level validation)
	if err := validator.ValidateRefreshToken(req.RefreshToken); err != nil {
		response.Error(w, http.StatusBadRequest,"validation error", response.ValidationError{
			Field: "refresh_token",
			Error: err.Error(),
		})
		return	
	}
	log.Printf("Refresh token attempt")

	// Call service layer to refresh the token
	access, refresh, err := h.authService.Refresh(
		r.Context(),
		req.RefreshToken,
	)

	if err != nil {
		if errors.Is(err, service.ErrInvalidRefreshToken) {
			h.mapError(w, err) // Map to 401 Unauthorized
			log.Printf("Refresh token failed: invalid refresh token")
			return
		}
	}

	// Return the user ID associated with the refresh token
	resp := map[string]string{
		"access_token": access,
		"refresh_token": refresh,
	}
	log.Printf("Refresh token successful, new access token issued")

	// Return success response with new access token
	response.JSON(w, http.StatusOK, resp)	
}

// 
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Decode req body into struct
	var req dto.RefreshTokenRequest

	if err := response.DecodeJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Basic validation
	if err := validator.ValidateRefreshToken(req.RefreshToken); err != nil {
		response.Error(w, http.StatusBadRequest, "validation error", response.ValidationError{
			Field: "refresh_token",
			Error: err.Error(),
		})
		return
	}
	log.Printf("Logout attempt")

	// Call service layer to logout (delete the refresh token)
	err := h.authService.Logout(
		r.Context(),
		req.RefreshToken,
	)

	if err != nil {
		if errors.Is(err, service.ErrInvalidRefreshToken) {
			h.mapError(w, err) // Map to 401 Unauthorized
			log.Printf("Logout failed: invalid refresh token")
			return
		}
	}	

	log.Printf("Logout successful, refresh token invalidated")
	// Return success response
	response.JSON(w, http.StatusOK, map[string]string{
		"message": "logged out successfully",
	})
}

func (h *AuthHandler) mapError(w http.ResponseWriter, err error) {
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