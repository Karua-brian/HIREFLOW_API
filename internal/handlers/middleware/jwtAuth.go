package middleware

import (
	"context"
	"job_board/internal/contextkeys"
	"job_board/internal/domain"
	"job_board/internal/validator"
	"job_board/pkg/response"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// JWTAuth is a middleware that checks for a valid JWT token in the Authorization header.
func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Get the token from the Authorization header
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

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid{
			response.Error(w, http.StatusUnauthorized, "Invalid token")
			return 
		}

		// Extract user info from token claims and add it to the request context
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return 
		}

		userIDFloat, ok1 := claims["user_id"].(float64)
		role, ok2 := claims["role"].(string)

		if !ok1 || !ok2 {
			http.Error(w, "Invalid token payload", http.StatusUnauthorized)
			return 
		}

		// Inject user into context
		user := &domain.User{
			ID: int64(userIDFloat),
			Role: role,
		}

		ctx := context.WithValue(r.Context(), contextkeys.UserKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}