package profile

import (
	"net/http"
	"os"

	"gesitr/internal/humaconfig"
	"gesitr/internal/profile/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const registrationHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Welcome to GesiTr</title>
<style>
  *, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }
  body {
    font-family: system-ui, -apple-system, sans-serif;
    display: flex; align-items: center; justify-content: center;
    min-height: 100vh;
    background: #0f172a; color: #e2e8f0;
  }
  .card {
    background: #1e293b; border-radius: 12px; padding: 2.5rem;
    max-width: 420px; width: 100%; text-align: center;
    box-shadow: 0 4px 24px rgba(0,0,0,.4);
  }
  h1 { font-size: 1.5rem; margin-bottom: .75rem; }
  p { color: #94a3b8; margin-bottom: 1.5rem; line-height: 1.5; }
  form { display: inline; }
  button {
    background: #3b82f6; color: #fff; border: none; border-radius: 8px;
    padding: .75rem 2rem; font-size: 1rem; cursor: pointer;
    transition: background .15s;
  }
  button:hover { background: #2563eb; }
</style>
</head>
<body>
<div class="card">
  <h1>Welcome to GesiTr</h1>
  <p>This seems to be your first visit. Create your profile to get started.</p>
  <button onclick="register()" id="btn">Create GesiTr Profile</button>
  <script>
    async function register() {
      const btn = document.getElementById('btn');
      btn.disabled = true;
      btn.textContent = 'Creating\u2026';
      const res = await fetch('/api/profile', { method: 'POST' });
      if (res.ok || res.status === 409) {
        window.location.href = '/';
      } else {
        btn.textContent = 'Something went wrong. Try again.';
        btn.disabled = false;
      }
    }
  </script>
</div>
</body>
</html>`

// RequireProfile returns a Gin handler that checks whether the current user
// has a profile. If not, it serves a simple registration page. If the user ID
// header is absent (e.g. in dev mode without oauth2-proxy), it falls through
// to the next handler without blocking.
func RequireProfile(db *gorm.DB, next gin.HandlerFunc) gin.HandlerFunc {
	header := humaconfig.AuthHeader
	usernameHeader := humaconfig.AuthUsernameHeader
	fallback := os.Getenv("AUTH_FALLBACK_USER")

	return func(c *gin.Context) {
		userID := c.GetHeader(header)
		if userID == "" {
			userID = fallback
		}
		if userID == "" {
			// No auth at all — skip profile check.
			next(c)
			return
		}

		var existing models.ProfileEntity
		if err := db.Where("user_id = ?", userID).First(&existing).Error; err == nil {
			// Profile exists — check if upstream username changed.
			if username := c.GetHeader(usernameHeader); username != "" && username != existing.Username {
				db.Model(&existing).Update("username", username)
			}
			next(c)
			return
		}

		// No profile — serve registration page.
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(registrationHTML))
	}
}
