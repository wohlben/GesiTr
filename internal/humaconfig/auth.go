package humaconfig

import (
	"context"
	"os"
)

type contextKey string

// UserIDContextKey is used to propagate the authenticated user ID
// through stdlib context for huma handlers.
const UserIDContextKey contextKey = "userID"

// UsernameContextKey propagates the authenticated username through stdlib context.
const UsernameContextKey contextKey = "username"

// AuthHeader is the HTTP header used to identify the user.
// Defaults to "X-User-Id" but can be overridden via AUTH_HEADER env var.
var AuthHeader = func() string {
	if h := os.Getenv("AUTH_HEADER"); h != "" {
		return h
	}
	return "X-User-Id"
}()

// AuthUsernameHeader is the HTTP header used to read the display username.
// Defaults to "X-Auth-Request-User" but can be overridden via AUTH_USERNAME_HEADER env var.
var AuthUsernameHeader = func() string {
	if h := os.Getenv("AUTH_USERNAME_HEADER"); h != "" {
		return h
	}
	return "X-Auth-Request-User"
}()

// GetUserID reads the user ID from a stdlib context (for huma handlers).
func GetUserID(ctx context.Context) string {
	if v, ok := ctx.Value(UserIDContextKey).(string); ok {
		return v
	}
	return ""
}

// GetUsername reads the username from a stdlib context (for huma handlers).
func GetUsername(ctx context.Context) string {
	if v, ok := ctx.Value(UsernameContextKey).(string); ok {
		return v
	}
	return ""
}
