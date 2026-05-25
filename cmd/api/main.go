package main

import (
	"context"
	"job_board/internal/app"
	"job_board/internal/config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	// Load configuration from .env
	cfg := config.LoadConfig()

	// Initialize the entire application.
	app := app.NewApp(cfg)

	// Create HTTP server
	server := &http.Server{
		Addr: ":" + cfg.PORT,
		Handler: app.Router,
	}

	// Run server in separate goroutine
	go func ()  {
		log.Printf("Server running on: http://localhost:%s\n", cfg.PORT)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Create channel to listen for shutdwon signals
	stop := make(chan os.Signal, 1)

	// Notify channel on interrupt/termination
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Block until signal received
	<-stop

	log.Println("Shutting down server...")

	// Create timeout context for shutdown
	ctx, cancel := context.WithTimeout(
		context.Background(),
		5 * time.Second,
	) 
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("graceful shutdown failed: %v", err)
	}
	log.Println("Server exited cleanly")
}