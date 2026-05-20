	package app

	import (
		"database/sql"
		"job_board/handlers"
		"job_board/middleware"
		"job_board/service"
		"job_board/store"
		"job_board/internal/db"
		"log"
		"log/slog"
		"net/http"
		"os"

		"github.com/go-chi/chi/v5"
		"github.com/joho/godotenv"
		_ "github.com/lib/pq"
	)

	// App encapsulates the entire application, including the router and all dependencies.
	type App struct {
		Router http.Handler
	}

	// NewApp initializes the application, sets up all dependencies, and returns an instance of App.
	func NewApp() *App {

		// Load .env file 
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		// Load env variables
		dsn := os.Getenv("DB_DSN")

		/* Fallback to constructing the DSN if DB_DSN is not set
		if dsn == "" {
			dbHost := os.Getenv("DB_HOST")
			dbPort := os.Getenv("DB_PORT")
			dbName := os.Getenv("DB_NAME")
			dbUser := os.Getenv("DB_USER")
			dbPassword := os.Getenv("DB_PASSWORD")

			dsn = "postgres://" + dbUser + ":" + dbPassword + "@" + dbHost + ":" + dbPort + "/" + dbName + "?sslmode=disable"
		}
			*/

		// Initialize the database connection:
		dbConn, err := sql.Open("postgres", dsn)
		if err != nil {
			log.Fatalf("Failed to connect to the database: %v", err)
		}

		// Verify database connection:
		err = dbConn.Ping()
		if err != nil {
			log.Fatalf("Failed to ping the database: %v", err)
		}

		// Run database migrations: after successful connection
		db.RunMigrations(dbConn) 

		dbName := os.Getenv("DB_NAME")
		log.Printf("Env database name: %s\n", dbName)
		err = db.QueryRow("SELECT current_database()").Scan(&dbName)
		log.Println("Connected DB:", dbName)

		// =================
		// Stores
		// =================

		jobStore := store.NewPostgresJobStore(db)
		applicationStore := store.NewPostgresApplicationStore(db)
		userStore := store.NewPostgresUserStore(db)
		refreshTokenStore := store.NewPostgresRefreshTokenStore(db)

		// =================
		// Worker Pool
		// =================

		workerPool := service.NewWorker(100, 4) // 100 jobs, 4 workers
		workerPool.Start()

		// =================
		// Services
		// =================

		jobService := service.NewJobService(jobStore, applicationStore, workerPool)
		authService := service.NewAuthService(userStore, refreshTokenStore)

		// =================
		// Handlers
		// =================

		// Initialize handlers (set up routes below)
		jobHandler := handlers.NewJobHandlers(jobService)
		authHandler := handlers.NewAuthHandlers(authService)

		// =================
		// Router
		// =================

		r := chi.NewRouter()

		// Global middleware

		logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
		r.Use(middleware.Logging(logger))

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
			r.Use(middleware.JWTAuth)

			// Job routes
			
			r.Post("/jobs", jobHandler.CreateJob)

			// Application routes
			r.Post("/jobs/{id}/apply", jobHandler.ApplyToJob)
		})

		return &App{
			Router: r,
		}

	}
