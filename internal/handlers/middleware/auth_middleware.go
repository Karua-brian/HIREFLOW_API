package middleware

import (
	"context"
	"job_board/internal/domain"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// contextKey is declared in another file in this package to avoid redeclaration.

// JWTAuth is a middleware that checks for a valid JWT token in the Authorization header.
func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Get the token from the Authorization header
		authHeader := r.Header.Get("Authorization")

		// Check if the Authorization header is present and has the correct format
		if authHeader == "" {
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			return 
		}

		// Extract the token from the header (assuming "Bearer <token>")
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return 
		}

		tokenString := parts[1]

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid{
			http.Error(w, "Invalid token", http.StatusUnauthorized)
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

		ctx := context.WithValue(r.Context(), userContextKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}