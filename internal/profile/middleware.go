package profile

import (
	"sync"

	"gesitr/internal/auth"
	"gesitr/internal/database"
	"gesitr/internal/profile/models"

	"github.com/gin-gonic/gin"
)

var knownProfiles sync.Map

// EnsureProfile returns a Gin middleware that creates a user profile
// on first request from an unknown user ID. It caches known IDs in
// memory to avoid a DB query on every request.
func EnsureProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := auth.GetUserID(c)
		if userID == "" {
			c.Next()
			return
		}

		if _, ok := knownProfiles.Load(userID); ok {
			c.Next()
			return
		}

		profile := models.UserProfileEntity{ID: userID, Name: userID}
		database.DB.FirstOrCreate(&profile, models.UserProfileEntity{ID: userID})
		knownProfiles.Store(userID, true)

		c.Next()
	}
}
