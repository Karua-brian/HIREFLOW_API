package app

import (
	"database/sql"
	"job_board/internal/config"
	"job_board/internal/db"
	"log"
	"os"
	"time"
)

// InitDB initializes the database connection using the provided
// configuration and returns a sql.DB instance.
func InitDB(cfg *config.Config) *sql.DB {
	// Initialize the database connection:
	dbConn, err := sql.Open("postgres", cfg.DBDSN)
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}

	// Verify database connection:
	if err := dbConn.Ping(); err != nil {
		log.Fatalf("DB ping failed: %v", err)
	}

	// Retry ping a few times in case db is still starting up (especially in Docker):
	for i := 0; i < 10; i++ {
		err = dbConn.Ping()
		if err == nil {
			break
		}
		// Log the error and wait before retrying:
		log.Printf("DB ping attempt %d failed: %v\n", i+1, err)
		// Wait a bit before retrying
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatalf("DB ping failed after retries: %v", err)
	}

	// Run database migrations to ensure schema is up to date:
	db.RunMigrations(dbConn)

	// Log the database name from env vars for verification:
	dbName := os.Getenv("DB_NAME")
	log.Printf("Env database name: %s\n", dbName)

	// Double-check connection by querying current database name:
	err = dbConn.QueryRow("SELECT current_database()").Scan(&dbName)
	log.Println("Connected DB:", dbName)

	return dbConn

}