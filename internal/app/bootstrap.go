package app

import (
	"database/sql"
	"job_board/internal/config"
	"job_board/internal/db"
	"os"
	"time"

	"go.uber.org/zap"
)

// InitDB initializes the database connection using the provided
// configuration and returns a sql.DB instance.
func InitDB(cfg *config.Config, logger *zap.Logger) *sql.DB {
	// Initialize the database connection:
	dbConn, err := sql.Open("postgres", cfg.DBDSN)
	if err != nil {
		logger.Fatal("DB connection failed:", zap.Error(err))
	}

	// Verify database connection:
	if err := dbConn.Ping(); err != nil {
		logger.Fatal("DB ping failed:", zap.Error(err))
	}

	// Retry ping a few times in case db is still starting up (especially in Docker):
	for i := 0; i < 10; i++ {
		err = dbConn.Ping()
		if err == nil {
			break
		}
		// Log the error and wait before retrying:
		logger.Info("DB ping attempt failed:", zap.Int("attempt", i+1), zap.Error(err))
		// Wait a bit before retrying
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		logger.Fatal("DB ping failed after retries:", zap.Error(err))
	}

	/* Run database migrations to ensure schema is up to date:
	db.RunMigrations(dbConn, logger)
	*/
	
	// Log the database name from env vars for verification:
	dbName := os.Getenv("DB_NAME")
	logger.Info("Env database name:", zap.String("db_name", dbName))

	// Double-check connection by querying current database name:
	err = dbConn.QueryRow("SELECT current_database()").Scan(&dbName)
	logger.Info("Connected DB:", zap.String("db_name", dbName))

	return dbConn

}