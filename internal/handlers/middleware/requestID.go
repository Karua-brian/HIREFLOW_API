package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type key int

const requestIDKey key = 0

// RequestID is a middleware that generates a unique request ID for each incoming HTTP request
// and adds it to the request context. This allows us to trace logs and correlate them with specific requests.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Generate a new UUID as the request ID
		requestID := uuid.New().String()

		// Add the request ID to the request context
		ctx := context.WithValue(r.Context(), requestIDKey, requestID)

		// Call the next handler with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetRequestID retrieves the request ID from the context. This can be used in handlers or other middleware to access the request ID.
func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(requestIDKey).(string); ok {
		return requestID
	}
	return ""
}