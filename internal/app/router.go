package app

import (
	"job_board/internal/handlers"
	"job_board/internal/handlers/middleware"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func NewRouter(jobHandler handlers.JobHandler, authHandler handlers.AuthHandler, recruiterHandler handlers.RecruiterHandler) http.Handler {
	r := chi.NewRouter()

	// Set up logging middleware with slog
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Add logging middleware to the router
	r.Use(middleware.Logging(logger))
	r.Use(middleware.CORS)

	// Public routes
	r.Post("/register", authHandler.Register)
	r.Post("/login", authHandler.Login)
	r.Post("/refresh", authHandler.Refresh)
	r.Post("/logout", authHandler.Logout)
	
	r.Get("/jobs", jobHandler.ListJobs)
	r.Get("/health", jobHandler.Health)

	// Protected routes (require authentication)
	r.Group(func(r chi.Router) {
		// Apply authentication middleware to all routes in this group
		r.Use(middleware.RequestID)
		r.Use(middleware.JWTAuth)

		// Job routes - only authenticated users can create jobs
		r.Post("/jobs", jobHandler.CreateJob)

		// Application routes - only authenticated users can apply to jobs
		r.Post("/jobs/{id}/apply", jobHandler.ApplyToJob)

		// Recruiter access request route - only authenticated users can request recruiter access
		r.Post("/recruiter/requests", recruiterHandler.RequestRecruiterAccesss)
		r.Get("/recruiter/status", recruiterHandler.GetMyRecruiterRequest) // Endpoint for users to check their recruiter request status

		// Admin routes for managing recruiter requests
		r.Get("/admin/recruiter-requests", recruiterHandler.ListRecruiterRequests)
		r.Put("/admin/recruiter-requests/{id}", recruiterHandler.UpdateRecruiterRequestStatus)

	})

	return r
}