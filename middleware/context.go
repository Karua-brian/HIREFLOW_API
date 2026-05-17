package middleware

import (
	"context"
	"job_board/domain"
)

type contextKey string

const userContextKey = contextKey("user")

// Helper to store a user (as an opaque value) in the context
func WithUser(ctx context.Context, user *domain.User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

// Helper to extract the user from context; returns the stored value and a bool indicating presence
func UserFromContext(ctx context.Context) (*domain.User, bool) {
	user, ok := ctx.Value(userContextKey).(*domain.User)
	return user, ok
}