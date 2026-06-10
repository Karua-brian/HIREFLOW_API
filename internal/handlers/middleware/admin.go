package middleware

import (
	"job_board/internal/contextkeys"
	"job_board/pkg/response"
	"net/http"
)

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user, ok := contextkeys.UserFromContext(r.Context())

		if !ok || user == nil {
			response.Error(w, http.StatusUnauthorized, "unauthorized")
			return 
		}

		if user.Role != "admin" {
			response.Error(w, http.StatusForbidden, "admin access required")
			return 
		}

		next.ServeHTTP(w, r)	
	})
}