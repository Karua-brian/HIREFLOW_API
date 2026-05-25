package app

import (
	"database/sql"
	"job_board/internal/config"
	"log"
	"os"
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

	// Log the database name from env vars for verification:
	dbName := os.Getenv("DB_NAME")
	log.Printf("Env database name: %s\n", dbName)

	// Double-check connection by querying current database name:
	err = dbConn.QueryRow("SELECT current_database()").Scan(&dbName)
	log.Println("Connected DB:", dbName)

	return dbConn

}