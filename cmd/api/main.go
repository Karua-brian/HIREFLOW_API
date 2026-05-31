package main

import (
	"context"
	"job_board/internal/app"
	"job_board/internal/config"
	"job_board/internal/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {

	// Initialize logger first so we can log any issues during startup
	 // or "production" based on cfg.ENV
	log := logger.Init("development")

	// Load configuration from .env
	cfg := config.LoadConfig(log)

	// Initialize logger first so we can log any issues during startup
	 // or "production" based on cfg.ENV
	defer log.Sync() // Flush any buffered log entries

	// Initialize the entire application.
	app := app.NewApp(cfg, log)

	// Create HTTP server
	server := &http.Server{
		Addr: ":" + cfg.PORT,
		Handler: app.Router,
		ReadTimeout: 15 * time.Second, // Set reasonable timeouts to prevent hanging connections 
		WriteTimeout: 10 * time.Second, 
		IdleTimeout: 120 * time.Second,
	}

	// Start server in separate goroutine
	go func ()  {
		log.Info("Server running on: http://localhost:8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Info("Server failed:", zap.Error(err))
		}
	}()

	// Create channel to listen for shutdwon signals
	stop := make(chan os.Signal, 1)

	// Notify channel on interrupt/termination
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Block until signal received
	<-stop

	log.Info("Shutting down server...")

	// Create timeout context for shutdown
	ctx, cancel := context.WithTimeout(
		context.Background(),
		5 * time.Second,
	) 
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Info("graceful shutdown failed:", zap.Error(err))
	}
	log.Info("Server exited cleanly")
}