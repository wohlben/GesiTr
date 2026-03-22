package profile

import (
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"

	"gesitr/internal/auth"
	"gesitr/internal/database"
	"gesitr/internal/profile/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func setupTestDB(t *testing.T) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?_foreign_keys=on"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatal(err)
	}
	db.AutoMigrate(&models.UserProfileEntity{})
	database.DB = db

	// Reset the cache between tests
	knownProfiles = sync.Map{}
}

func closeDB(t *testing.T) {
	t.Helper()
	sqlDB, err := database.DB.DB()
	if err != nil {
		t.Fatal(err)
	}
	sqlDB.Close()
}

func TestEnsureProfile_CreatesOnFirstRequest(t *testing.T) {
	setupTestDB(t)
	defer closeDB(t)

	r := gin.New()
	r.Use(auth.UserID())
	r.Use(EnsureProfile())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-User-Id", "newuser")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var profile models.UserProfileEntity
	if err := database.DB.First(&profile, "id = ?", "newuser").Error; err != nil {
		t.Fatal("profile was not created:", err)
	}
	if profile.Name != "newuser" {
		t.Errorf("expected default name newuser, got %s", profile.Name)
	}
}

func TestEnsureProfile_SkipsOnSubsequentRequests(t *testing.T) {
	setupTestDB(t)
	defer closeDB(t)

	r := gin.New()
	r.Use(auth.UserID())
	r.Use(EnsureProfile())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	// First request creates profile
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-User-Id", "repeat")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Second request should use cache
	req = httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-User-Id", "repeat")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// Verify only one profile exists
	var count int64
	database.DB.Model(&models.UserProfileEntity{}).Where("id = ?", "repeat").Count(&count)
	if count != 1 {
		t.Errorf("expected 1 profile, got %d", count)
	}
}

func TestEnsureProfile_PreservesExistingProfile(t *testing.T) {
	setupTestDB(t)
	defer closeDB(t)

	// Pre-create profile with custom name
	database.DB.Create(&models.UserProfileEntity{ID: "existing", Name: "Custom Name"})

	r := gin.New()
	r.Use(auth.UserID())
	r.Use(EnsureProfile())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-User-Id", "existing")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var profile models.UserProfileEntity
	database.DB.First(&profile, "id = ?", "existing")
	if profile.Name != "Custom Name" {
		t.Errorf("expected preserved name Custom Name, got %s", profile.Name)
	}
}
