package auth

import (
	"context"
	"net/http"
	"os"

	"gesitr/internal/humaconfig"

	"github.com/gin-gonic/gin"
)

const userIDKey = "userID"

// UserID returns a Gin middleware that reads the X-User-Id header
// and stores it in the context. If the header is missing, it falls
// back to AUTH_FALLBACK_USER. If that env var is also unset, the
// request is rejected with 401.
func UserID() gin.HandlerFunc {
	fallback := os.Getenv("AUTH_FALLBACK_USER")

	return func(c *gin.Context) {
		id := c.GetHeader("X-User-Id")
		if id == "" {
			id = fallback
		}
		if id == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing X-User-Id header"})
			return
		}
		c.Set(userIDKey, id)
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), humaconfig.UserIDContextKey, id))
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
