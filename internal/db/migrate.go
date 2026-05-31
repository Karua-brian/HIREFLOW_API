package db

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"go.uber.org/zap"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations runs the database migrations using the provided database URL
func RunMigrations(db *sql.DB, logger *zap.Logger) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		logger.Fatal("migration driver error:", zap.Error(err))
	} 
	
	// Create a new migrate instance with the file source and database URL
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)

	if err != nil {
		logger.Fatal("migration init error:", zap.Error(err))
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Fatal("migration run error:", zap.Error(err))
	} 

	logger.Info("Migrations applied successfully")
	
}