package store

import (
	"context"
	"database/sql"
	"errors"
	"job_board/domain"
)

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{db: db}
}

// CreateUser inserts a new user into the database
func (s *PostgresUserStore) CreateUser(ctx context.Context, user *domain.User) error {
	// Query the new user into the database
	query := `
	INSERT INTO users (email, password, role)
	VALUES ($1, $2, $3)
	RETURNING id, created_at, updated_at
	`
	// Execute the query and scan the generated ID back into the user struct
	return s.db.QueryRowContext(
		ctx,
		query,
		user.Email,
		user.Password,
		user.Role,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt) // Get the generated ID and set it on the user struct
}

func (s *PostgresUserStore) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	// Query the database for a user with the given email
	row := s.db.QueryRowContext( 
		ctx, 
		`SELECT id, email, password, role
		FROM users
		WHERE email = $1`,
		email,
	)
	// Scan the result into a User struct
	user := &domain.User{}

	// Handle the case where no user is found
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Role,
	)

	// If no user is found, return nil without an error 
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	// Return the found user
	return user, nil
}	

// GetUserByID retrieves a user by their ID
func (s *PostgresUserStore) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {

	// Query the database for a user with the given ID
	query := `
	SELECT id, email, password, role
	FROM users
	WHERE id = $1
	`
	row := s.db.QueryRowContext(ctx, query, id)

	// Scan the result into a User struct
	user := &domain.User{}
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Role,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}
