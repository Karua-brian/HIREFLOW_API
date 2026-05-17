package store

import (
	"context"
	"errors"
	"time"
)

// RefreshTokenStore defines how the service interacts with persistance for refresh tokens
type RefreshTokenStore interface {
	SaveToken(ctx context.Context, userID int64, token string, expires time.Time) error
	GetUserIDByToken(ctx context.Context, token string) (int64, error)
	DeleteToken(ctx context.Context, token string) error
	DeleteExpired(ctx context.Context) error
}

var ErrInvalidRefreshToken = errors.New("invalid refresh token")