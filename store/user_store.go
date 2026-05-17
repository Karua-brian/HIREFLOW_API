package store

import (
	"context"
	"job_board/domain"
)

// UserStore defines how the service interacts with persistance for user data
type UserStore interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserByID(ctx context.Context, id int64) (*domain.User, error)
}
