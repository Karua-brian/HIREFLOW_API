package app

import (
	"job_board/internal/config"
	"job_board/internal/handlers"
	"job_board/internal/repository"
	"job_board/internal/service"
	"net/http"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

// App encapsulates the entire application, including the router and all dependencies.
	type App struct {
		Router http.Handler
	}

	// NewApp initializes the application, sets up all dependencies, and returns an instance of App.
	func NewApp(cfg *config.Config, logger *zap.Logger) *App {

		// Initialize db connection first, since many components depend on it
		db := InitDB(cfg, logger)

		// =================
		// Stores
		// =================

		jobRepo := repository.NewPostgresJobStore(db)
		applicationRepo := repository.NewPostgresApplicationStore(db)
		userRepo := repository.NewPostgresUserStore(db)
		refreshTokenRepo := repository.NewPostgresRefreshTokenStore(db)

		// =================
		// Worker Pool
		// =================

		workerPool := service.NewWorker(100, 4) // 100 jobs, 4 workers
		workerPool.Start()

		// =================
		// Services
		// =================

		jobService := service.NewJobService(jobRepo, applicationRepo, workerPool)
		authService := service.NewAuthService(userRepo, refreshTokenRepo, logger)

		// =================
		// Handlers
		// =================

		jobHandler := handlers.NewJobHandlers(jobService, logger)
		authHandler := handlers.NewAuthHandlers(authService, logger)

		// =================
		// Router
		// =================
		
		router := NewRouter(jobHandler, authHandler)

		return &App{
			Router: router,
		}

	}
