package db

import (
	"context"
	"database/sql"
	"job_board/internal/domain"
	"job_board/internal/repository"
	"os"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func SeedAdmin(ctx context.Context, db *sql.DB, logger *zap.Logger) error {

	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	if adminEmail == "" || adminPassword == "" {
		logger.Info("Admin seed skipped (missing env vars)")
		return nil	
	}

	repo := repository.NewPostgresUserRepo(db)

	existing, err := repo.GetUserByEmail(ctx, adminEmail)
	if err != nil {
		return err
	}

	if existing != nil {
		logger.Info("Admin already exixts, skipping seed")
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	logger.Info("TRYING TO SEED ADMIN", zap.String("email", adminEmail))

	user := &domain.User{
		Email:    	adminEmail,
		Password: 	string(hash),
		Role: 		"admin",
	}

	err = repo.CreateUser(ctx, user)
	if err != nil {
		logger.Info("seed insert failed", zap.Error(err))
		return err
	}

	logger.Info("Admin user created successfully")

	return nil
}
