package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(UserID())
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, GetUserID(c))
	})
	return r
}

func TestUserID_WithHeader(t *testing.T) {
	r := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("X-User-Id", "alice")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if w.Body.String() != "alice" {
		t.Errorf("expected 'alice', got '%s'", w.Body.String())
	}
}

func TestUserID_WithFallback(t *testing.T) {
	t.Setenv("AUTH_FALLBACK_USER", "anon")
	r := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if w.Body.String() != "anon" {
		t.Errorf("expected 'anon', got '%s'", w.Body.String())
	}
}

func TestUserID_NoHeaderNoFallback(t *testing.T) {
	t.Setenv("AUTH_FALLBACK_USER", "")
	r := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}
