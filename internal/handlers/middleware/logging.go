package middleware

import (
	"context"
	"job_board/internal/contextkeys"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type responseRecorder struct {
	http.ResponseWriter
	status int
}

func (r *responseRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func Logging(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request)  {

			start := time.Now()
			reqID := uuid.New().String()

			// Attach request ID to context
			ctx := context.WithValue(r.Context(), contextkeys.RequestIDKey, reqID)
			r = r.WithContext(ctx)

			// Add request ID to response header
			w.Header().Set("X-Request-ID", reqID)

			rec := &responseRecorder{
				ResponseWriter: w,
				status: 		http.StatusOK,	
			}

			// Call next handler
			next.ServeHTTP(rec, r)

			duration := time.Since(start)

			// Structured logging
			logger.Info(
				"http_request",
				"request_id", reqID,
				"method", r.Method,
				"path", r.URL.Path,
				"status", rec.status,
				"duration_ms", duration.Milliseconds(),
				"remote_addr", r.RemoteAddr,
			)
			})
	}
}