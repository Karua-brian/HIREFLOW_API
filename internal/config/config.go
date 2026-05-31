package config

import (
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

// Config holds all configuration values for the application,
// such as database connection info, JWT secret, etc.
type Config struct {
	// Database configuration
	DBDSN 	  string
	PORT 	  string
	JWTSecret string
}

func LoadConfig(logger *zap.Logger) *Config {
	// Try loading .env locally
	// Only load .env file if not running in Docker, 
	// since Docker will pass env vars directly
	if err := godotenv.Load(); err != nil {
		logger.Info("No .env file found (using env vars)")
	} else {
		logger.Info(".env file loaded successfully")
	}

	// Load env variables into Config struct
	cfg := &Config{
		DBDSN: 	  os.Getenv("DB_DSN"),
		PORT: 	  os.Getenv("PORT"),
		JWTSecret: os.Getenv("JWT_SECRET"),
	}

	// Set defaults
	if cfg.PORT == "" {
		cfg.PORT = "8080"
	}

	// Validation: ensure required config values are set
	if cfg.DBDSN == "" {
		logger.Fatal("DB_DSN is required")
	}

	if cfg.JWTSecret == "" {
		logger.Fatal("JWT_SECRET is required")
	}

	return cfg
}