package store

import (
	"context"
	"database/sql"
	"fmt"
	"job_board/security"
	"log"
	"strings"
	"time"
)

type PostgresRefreshTokenStore struct {
	db *sql.DB
}

func NewPostgresRefreshTokenStore(db *sql.DB) *PostgresRefreshTokenStore {
	return &PostgresRefreshTokenStore{db: db}
}

func (s *PostgresRefreshTokenStore) SaveToken(ctx context.Context, userID int64, token string, expires time.Time) error {

	token = strings.TrimSpace(token)

	tokenHash := security.HashToken(token)

	// Insert the refresh token into the database
	query := `
	INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
	VALUES ($1, $2, $3)
	`
	_, err := s.db.ExecContext(
		ctx,
		query,
		userID,
		tokenHash,
		expires,
	)

	log.Printf("Token Raw: %s", token)
	log.Println("Save Token:", token)
	log.Println("Save Hash:", tokenHash)

	return err
}

func (s *PostgresRefreshTokenStore) GetUserIDByToken(ctx context.Context, token string) (int64, error) {

	token = strings.TrimSpace(token)

	tokenHash := security.HashToken(token)

	// Query the db for the user ID associated with the given refresh token, ensuring the token is not expired
	query := `
	SELECT user_id
	FROM refresh_tokens
	WHERE token_hash = $1 
	AND expires_at > NOW()
	`

	// Scan the result into a userID variable
	var userID int64
	err := s.db.QueryRowContext(
		ctx,
		query,
		tokenHash,
	).Scan(&userID)

	if err!= nil {
		if err == sql.ErrNoRows {
			return 0, ErrInvalidRefreshToken // No user found for the token
		}
		return 0, err // Some other database error occurred
	}
	log.Println("Incoming Token:", token)
	log.Println("Incoming Hash:", tokenHash)
	return userID, nil
}

// DeleteToken removes a refresh token from the database
func (s *PostgresRefreshTokenStore) DeleteToken(ctx context.Context, token string) error {

	token = strings.TrimSpace(token)

	tokenHash := security.HashToken(token)
	// Delete the refresh token from the database
	query := 
	`
	DELETE FROM refresh_tokens
	WHERE token_hash = $1
	`

	_, err := s.db.ExecContext(
		ctx,
		query,
		tokenHash,
	)

	return err
}

func (s *PostgresRefreshTokenStore) DeleteExpired(ctx context.Context) error {

	// Delete expired refresh tokens from database
	query := `
	DELETE FROM refresh_tokens
	WHERE expires_at <= NOW()
	`

	res, err := s.db.ExecContext(
		ctx,
		query,
	)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows > 0 {
		fmt.Printf("Deleted %d expired refresh tokens\n", rows)
	}
	
	return nil
}