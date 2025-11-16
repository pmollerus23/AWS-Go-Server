package auth

import (
	"context"
	"errors"
)

// Context keys for storing auth data
type contextKey string

const (
	userContextKey contextKey = "user"
)

// Errors
var (
	ErrNoUserInContext = errors.New("no user found in context")
)

// WithUser adds a user to the request context.
func WithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

// GetUser retrieves the user from the request context.
func GetUser(ctx context.Context) (*User, error) {
	user, ok := ctx.Value(userContextKey).(*User)
	if !ok || user == nil {
		return nil, ErrNoUserInContext
	}
	return user, nil
}

// MustGetUser retrieves the user from context, panics if not found.
// Use only in handlers where authentication is required.
func MustGetUser(ctx context.Context) *User {
	user, err := GetUser(ctx)
	if err != nil {
		panic(err)
	}
	return user
}

// GetUserID retrieves just the user ID from context.
func GetUserID(ctx context.Context) (string, error) {
	user, err := GetUser(ctx)
	if err != nil {
		return "", err
	}
	return user.ID, nil
}

// IsAuthenticated checks if a user is authenticated in the context.
func IsAuthenticated(ctx context.Context) bool {
	_, err := GetUser(ctx)
	return err == nil
}
