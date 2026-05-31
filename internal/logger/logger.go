package logger

import (
	"go.uber.org/zap"
)

// Init initializes the global logger instance.
// This function should be called at the start of the application.
func Init(env string) *zap.Logger {
	var logger *zap.Logger
	var err error

	// Configure logger based on environment
	switch env {
	case "production":
		logger, err = zap.NewProduction()
	default:
		logger, err = zap.NewDevelopment()
	}

	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}

	return logger
}
