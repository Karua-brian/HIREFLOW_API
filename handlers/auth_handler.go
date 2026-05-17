package handlers

import (
	"encoding/json"
	"errors"
	"job_board/service"
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

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Basic validation (transport-level validation)
	if input.Email == "" || input.Password == "" || input.Role == "" {
		http.Error(w, "missing required fileds", http.StatusBadRequest)
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
			http.Error(w, "user already exists", http.StatusConflict)
			log.Printf("User with email %s already exists", input.Email)
			return
		}
		log.Printf("Error registering user: %v", err)

		http.Error(w, "failed to register user",  http.StatusInternalServerError)
		return
	}
	log.Printf("User with email %s successfully registered", input.Email)

	// Return success response
	w.WriteHeader(http.StatusCreated)
}

// Login handles POST /login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

	// Decode request body into struct
	// We only accept Email, Password from client
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	log.Printf("Login attempt for email: %s", input.Email)

	// Basic validation (transport-level validation)
	if input.Email == "" || input.Password == "" {
		http.Error(w, "missing required fields", http.StatusBadRequest)
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
			http.Error(w, "invalid email or password", http.StatusUnauthorized)
			return
		}

		http.Error(w, "failed to login", http.StatusInternalServerError)
		return
	}

	// Return the token in the response
	resp := map[string]string{
		"access_token" : token,
		"refresh_token": refresh,
		
	}

	// Return success response with token
	log.Printf("User with email %s successfully logged in", input.Email)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Refresh handles POST /refresh-token
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {

	// Decode req body into struct
	var input struct {
		RefreshToken string `json:"refresh_token"`
	}

	log.Printf("Refresh token attempt")
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if input.RefreshToken == "" {
		http.Error(w, "missing refresh token", http.StatusBadRequest)
		return
	}

	// Call service layer to refresh the token
	access, refresh, err := h.authService.Refresh(
		r.Context(),
		input.RefreshToken,
	)

	if err != nil {
		http.Error(w, "invalid refresh token", http.StatusUnauthorized)
		return
	}

	// Return the user ID associated with the refresh token
	resp := map[string]string{
		"access_token": access,
		"refresh_token": refresh,
	}

	log.Printf("Refresh token successful, new access token issued")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// 
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Decode req body into struct
	var input struct {
		RefreshToken string `json:"refresh_token"`
	}

	log.Printf("Logout attempt")
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if input.RefreshToken == "" {
		http.Error(w, "missing refresh token", http.StatusBadRequest)
		return
	}

	// Call service layer to logout (delete the refresh token)
	err := h.authService.Logout(
		r.Context(),
		input.RefreshToken,
	)

	if err != nil {
		http.Error(w, "failed to logout", http.StatusInternalServerError)
		return
	}

	log.Printf("Logout successful, refresh token invalidated")
	w.WriteHeader(http.StatusOK)
}