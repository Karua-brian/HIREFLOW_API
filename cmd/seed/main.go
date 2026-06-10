package main

import (
	"context"
	"database/sql"
	"job_board/internal/config"
	"job_board/internal/db"

	"go.uber.org/zap"
)

func main() {

	logger,_ := zap.NewProduction()

	cfg := config.LoadConfig(logger)

	conn, err := sql.Open("postgres", cfg.DBDSN)
	if err != nil {
		logger.Fatal("db connection failed", zap.Error(err))
	}

	if err := conn.Ping(); err != nil {
		logger.Fatal("db ping failed", zap.Error(err))
	}

	if err := db.SeedAdmin(context.Background(), conn, logger); err != nil {
		logger.Fatal("admin seed failed", zap.Error(err))
	}

	logger.Info("seed completed successfully")
}