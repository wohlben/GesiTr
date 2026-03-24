package humaconfig

import "context"

type contextKey string

// UserIDContextKey is used to propagate the authenticated user ID
// through stdlib context for huma handlers.
const UserIDContextKey contextKey = "userID"

// GetUserID reads the user ID from a stdlib context (for huma handlers).
func GetUserID(ctx context.Context) string {
	if v, ok := ctx.Value(UserIDContextKey).(string); ok {
		return v
	}
	return ""
}
