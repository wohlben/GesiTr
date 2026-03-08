package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"gesitr/internal/database"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	database.DB = db
}

func TestAutoMigrate(t *testing.T) {
	setupTestDB(t)
	autoMigrate()

	tables := []string{
		"exercises", "exercise_forces", "exercise_muscles",
		"exercise_measurement_paradigms", "exercise_instructions",
		"exercise_images", "exercise_alternative_names",
		"equipment", "exercise_equipments", "fulfillments",
		"exercise_relationships", "exercise_groups", "exercise_group_members",
	}
	for _, table := range tables {
		if !database.DB.Migrator().HasTable(table) {
			t.Errorf("table %q was not created", table)
		}
	}
}

func TestSetupRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB(t)
	autoMigrate()

	r := gin.New()
	setupRoutes(r)

	routes := r.Routes()
	expected := map[string]bool{
		"GET /api/exercises":                                         false,
		"POST /api/exercises":                                        false,
		"GET /api/exercises/:id":                                     false,
		"PUT /api/exercises/:id":                                     false,
		"DELETE /api/exercises/:id":                                  false,
		"GET /api/exercises/:id/versions":                            false,
		"GET /api/exercises/templates/:templateId/versions/:version": false,
		"GET /api/equipment":                                         false,
		"POST /api/equipment":                                        false,
		"GET /api/equipment/:id":                                     false,
		"PUT /api/equipment/:id":                                     false,
		"DELETE /api/equipment/:id":                                  false,
		"GET /api/equipment/:id/versions":                            false,
		"GET /api/equipment/templates/:templateId/versions/:version": false,
		"GET /api/fulfillments":                                      false,
		"POST /api/fulfillments":                                     false,
		"DELETE /api/fulfillments/:id":                               false,
		"GET /api/exercise-relationships":                            false,
		"POST /api/exercise-relationships":                           false,
		"DELETE /api/exercise-relationships/:id":                     false,
		"GET /api/exercise-groups":                                   false,
		"POST /api/exercise-groups":                                  false,
		"GET /api/exercise-groups/:id":                               false,
		"PUT /api/exercise-groups/:id":                               false,
		"DELETE /api/exercise-groups/:id":                            false,
		"GET /api/exercise-group-members":                            false,
		"POST /api/exercise-group-members":                           false,
		"DELETE /api/exercise-group-members/:id":                     false,
	}

	for _, route := range routes {
		key := route.Method + " " + route.Path
		if _, ok := expected[key]; ok {
			expected[key] = true
		}
	}

	for key, found := range expected {
		if !found {
			t.Errorf("route %q not registered", key)
		}
	}
}

func TestSetupSPA(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	setupSPA(r)

	t.Run("serves index.html for unknown routes", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/some/spa/route", nil)
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK && w.Code != http.StatusMovedPermanently {
			t.Errorf("expected 200 or 301 for SPA fallback, got %d", w.Code)
		}
	})

	t.Run("serves existing static file", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/favicon.ico", nil)
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("expected 200 for static file, got %d", w.Code)
		}
	})
}

func TestBuildApp(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// buildApp calls database.Init() which creates gesitr.db
	// Run in a temp directory to avoid polluting the project
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	// Set DEV=true to skip SPA setup (embedded files may not resolve from temp dir)
	origDev := os.Getenv("DEV")
	os.Setenv("DEV", "true")
	defer os.Setenv("DEV", origDev)

	r := buildApp()
	if r == nil {
		t.Fatal("buildApp returned nil")
	}

	// Verify routes are registered
	routes := r.Routes()
	if len(routes) == 0 {
		t.Error("no routes registered")
	}

	// Verify DB tables were created
	if !database.DB.Migrator().HasTable("exercises") {
		t.Error("exercises table not created")
	}
}

func TestBuildAppWithSPA(t *testing.T) {
	gin.SetMode(gin.TestMode)

	dir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	// Unset DEV to trigger SPA setup
	origDev := os.Getenv("DEV")
	os.Unsetenv("DEV")
	defer os.Setenv("DEV", origDev)

	r := buildApp()
	if r == nil {
		t.Fatal("buildApp returned nil")
	}

	// Verify SPA serves files
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/favicon.ico", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 for static file via SPA, got %d", w.Code)
	}
}
