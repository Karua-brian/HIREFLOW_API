package app

import (
	"job_board/internal/handlers"
	"job_board/internal/handlers/middleware"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func NewRouter(
		jobHandler handlers.JobHandler, 
		authHandler handlers.AuthHandler, 
		recruiterRequestHandler handlers.RecruiterRequestHandler,
		adminHandler handlers.AdminHandler,
		notificationHandler handlers.NotificationHandler,
	) http.Handler {

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
		r.Post("/recruiter/requests", recruiterRequestHandler.RequestRecruiterAccess)
		r.Get("/recruiter/requests/me", recruiterRequestHandler.GetMyRecruiterRequest) // Endpoint for users to check their recruiter request status
	})
	
	// Protected AdminOnly routes
	r.Route("/admin", func(r chi.Router) {
		r.Use(middleware.RequestID)
		r.Use(middleware.JWTAuth)
		r.Use(middleware.AdminOnly)

		r.Get("/recruiter-requests", adminHandler.ListRecruiterRequests)
		r.Post("/recruiter-requests/{id}/approve", adminHandler.ApproveRecruiterRequest)
		r.Post("/recruiter-requests/{id}/reject", adminHandler.RejectRecruiterRequest)
	})

	r.Route("/notifications", func(r chi.Router) {
		r.Use(middleware.RequestID)
		r.Use(middleware.JWTAuth)

		//r.Get("/notifications", notificationHandler.CreateNotification)
		r.Get("/me", notificationHandler.GetMyNotifications)
		r.Patch("/read-all", notificationHandler.MarkAllAsRead)
	})

	return r
}