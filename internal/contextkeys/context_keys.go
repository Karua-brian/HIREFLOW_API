package contextkeys

import (
	"context"
	"job_board/internal/domain"
)

// Define custom types for context keys to avoid collisions and ensure type safety.
type key string

const UserKey key = "user"
const RequestIDKey key = "request_id"

// Helper to store a user in the context (testing purposes)
func WithUser(ctx context.Context, user *domain.User) context.Context {
	return context.WithValue(ctx, UserKey, user)
}

// Helper to extract the user from context; returns the stored value and a bool indicating presence
func UserFromContext(ctx context.Context) (*domain.User, bool) {
	user, ok := ctx.Value(UserKey).(*domain.User)
	return user, ok
}

