package humaconfig

import (
	"context"
	"os"
)

type contextKey string

// UserIDContextKey is used to propagate the authenticated user ID
// through stdlib context for huma handlers.
const UserIDContextKey contextKey = "userID"

// AuthHeader is the HTTP header used to identify the user.
// Defaults to "X-User-Id" but can be overridden via AUTH_HEADER env var.
var AuthHeader = func() string {
	if h := os.Getenv("AUTH_HEADER"); h != "" {
		return h
	}
	return "X-User-Id"
}()

// GetUserID reads the user ID from a stdlib context (for huma handlers).
func GetUserID(ctx context.Context) string {
	if v, ok := ctx.Value(UserIDContextKey).(string); ok {
		return v
	}
	return ""
}
