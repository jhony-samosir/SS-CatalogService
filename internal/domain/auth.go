package domain

import (
	"context"
)

type contextKey string

const userContextKey contextKey = "user_context"

// UserContext holds information about the authenticated user.
type UserContext struct {
	UserID   string
	FullName string
	SellerID *int
	Roles    []string
}

// ContextWithUser returns a new context with the UserContext value.
func ContextWithUser(ctx context.Context, user UserContext) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

// UserFromContext extracts UserContext from the context.
func UserFromContext(ctx context.Context) (UserContext, bool) {
	u, ok := ctx.Value(userContextKey).(UserContext)
	return u, ok
}
