package handlers

import (
	"encoding/json"
	"errors"
	"job_board/internal/service"
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
	
	// Decode request body into struct
	// We only accept Email, Password, Role from client
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"` // "recruiter" or "admin"
	}

	if err := response.DecodeJSON(r, &input); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Basic validation (transport-level validation)
	if input.Email == "" || input.Password == "" || input.Role == "" {
		response.Error(w, http.StatusBadRequest, "missing required fileds", response.ValidationError{
			Field: "email/password/role",
			Error: "email, password and role are required",
		})
		return
	}

	// Call service layer (business rules happen there)
	err := h.authService.Register(
		r.Context(),
		input.Email,
		input.Password,
		input.Role,
	)


	log.Printf("Attempting to register user with email: %s and role: %s", input.Email, input.Role)
	if err != nil {
		if errors.Is(err, service.ErrUserExists) {
			response.Error(w, http.StatusConflict, "user already exists")
			log.Printf("User with email %s already exists", input.Email)
			return
		}
		log.Printf("Error registering user: %v", err)

		response.Error(w, http.StatusInternalServerError, "failed to register user")
		return
	}
	log.Printf("User with email %s successfully registered", input.Email)

	// Return success response
	response.JSON(w, http.StatusCreated, map[string]string{
		"message": "user registered successfully", 
		})	
}

// Login handles POST /login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

	// Decode request body into struct
	// We only accept Email, Password from client
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Use the response helper to decode JSON body and handle errors
	if err := response.DecodeJSON(r, &input); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	log.Printf("Login attempt for email: %s", input.Email)

	// Basic validation (transport-level validation)
	if input.Email == "" || input.Password == "" {
		response.Error(w, http.StatusBadRequest, "missing required fields", response.ValidationError{
			Field: "email/password",
			Error: "email and password are required",
		})
		return
	}

	// Call service layer
	token, refresh, err := h.authService.Login(
		r.Context(),
		input.Email,
		input.Password,
	)

	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			response.Error(w, http.StatusUnauthorized, "invalid email or password")
			return
		}

		response.Error(w, http.StatusInternalServerError, "failed to login")
		return
	}

	// Return the token in the response
	resp := map[string]string{
		"access_token" : token,
		"refresh_token": refresh,
		
	}

	// Return success response with token
	log.Printf("User with email %s successfully logged in", input.Email)

	// Use the response helper to send a JSON response
	response.JSON(w, http.StatusOK, resp)
}

// Refresh handles POST /refresh-token
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {

	// Decode req body into struct
	var input struct {
		RefreshToken string `json:"refresh_token"`
	}

	log.Printf("Refresh token attempt")
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest,"invalid request body")
		return
	}

	// Basic validation (transport-level validation)
	if input.RefreshToken == "" {
		response.Error(w, http.StatusBadRequest, "missing refresh token", response.ValidationError{
			Field: "refresh_token",
			Error: "refresh token is required",
		})
		return
	}

	// Call service layer to refresh the token
	access, refresh, err := h.authService.Refresh(
		r.Context(),
		input.RefreshToken,
	)

	if err != nil {
		response.Error(w, http.StatusUnauthorized,"invalid refresh token")
		return
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
	var input struct {
		RefreshToken string `json:"refresh_token"`
	}

	log.Printf("Logout attempt")
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if input.RefreshToken == "" {
		response.Error(w, http.StatusBadRequest, "missing refresh token", response.ValidationError{
			Field: "refresh_token",
			Error: "refresh token is required",
		})
		return
	}

	// Call service layer to logout (delete the refresh token)
	err := h.authService.Logout(
		r.Context(),
		input.RefreshToken,
	)

	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to logout")
		return
	}

	log.Printf("Logout successful, refresh token invalidated")
	// Return success response
	response.JSON(w, http.StatusOK, map[string]string{
		"message": "logged out successfully",
	})
}