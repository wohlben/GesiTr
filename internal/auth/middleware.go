package auth

import (
	"context"
	"net/http"
	"os"

	"gesitr/internal/humaconfig"

	"github.com/gin-gonic/gin"
)

const userIDKey = "userID"

// UserID returns a Gin middleware that reads the auth header (configurable
// via AUTH_HEADER env var, default "X-User-Id") and stores it in the context.
// If the header is missing, it falls back to AUTH_FALLBACK_USER. If that env
// var is also unset, the request is rejected with 401.
func UserID() gin.HandlerFunc {
	fallback := os.Getenv("AUTH_FALLBACK_USER")

	header := humaconfig.AuthHeader

	usernameHeader := humaconfig.AuthUsernameHeader

	return func(c *gin.Context) {
		id := c.GetHeader(header)
		if id == "" {
			id = fallback
		}
		if id == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing " + header + " header"})
			return
		}
		c.Set(userIDKey, id)
		ctx := context.WithValue(c.Request.Context(), humaconfig.UserIDContextKey, id)

		if username := c.GetHeader(usernameHeader); username != "" {
			c.Set("username", username)
			ctx = context.WithValue(ctx, humaconfig.UsernameContextKey, username)
		}

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// GetUserID retrieves the user ID from the Gin context.
func GetUserID(c *gin.Context) string {
	if id, ok := c.Get(userIDKey); ok {
		return id.(string)
	}
	return ""
}
