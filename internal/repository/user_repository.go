package repository

import (
	"context"
	"database/sql"
	"job_board/internal/domain"

	"github.com/google/uuid"
)

// UserStore defines how the service interacts with persistance for user data
type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error)
}

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepo(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

// CreateUser inserts a new user into the database
func (s *PostgresUserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	// Query the new user into the database
	query := `
	INSERT INTO users (email, password_hash, role)
	VALUES ($1, $2, $3)
	RETURNING id, created_at, updated_at
	`
	// Execute the query and scan the generated ID back into the user struct
	err := s.db.QueryRowContext(
		ctx,
		query,
		user.Email,
		user.Password,
		user.Role,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt) // Get the generated ID and set it on the user struct

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresUserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	// Query the database for a user with the given email
	query := `
		SELECT id, email, password_hash, role, created_at, updated_at
		FROM users
		WHERE email = $1
		`

	user :=  &domain.User{}

	err := s.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	// If no user is found, return nil without an error
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return user, nil
}

// GetUserByID retrieves a user by their ID
func (s *PostgresUserRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error) {

	// Query the database for a user with the given ID
	query := `
	SELECT id, created_at, updated_at
	FROM users
	WHERE id = $1
	`
	row := s.db.QueryRowContext(ctx, query, userID)

	// Scan the result into a User struct
	user := &domain.User{}

	err := row.Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	// If no user is found, return nil without an error
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	
	return user, nil
}


