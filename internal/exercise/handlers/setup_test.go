package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"os"
	"testing"

	"gesitr/internal/auth"
	"gesitr/internal/database"
	"gesitr/internal/exercise/models"
	profilemodels "gesitr/internal/profile/models"

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
	t.Setenv("AUTH_FALLBACK_USER", "testuser")
	db, err := gorm.Open(sqlite.Open("file::memory:?_foreign_keys=on"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatal(err)
	}
	db.AutoMigrate(
		&profilemodels.UserProfileEntity{},
		&models.ExerciseEntity{},
		&models.ExerciseForce{},
		&models.ExerciseMuscle{},
		&models.ExerciseMeasurementParadigm{},
		&models.ExerciseInstruction{},
		&models.ExerciseImage{},
		&models.ExerciseAlternativeName{},
		&models.ExerciseEquipment{},
		&models.ExerciseHistoryEntity{},
		&models.ExerciseSchemeEntity{},
	)
	db.Create(&profilemodels.UserProfileEntity{ID: "testuser", Name: "testuser"})
	db.Create(&profilemodels.UserProfileEntity{ID: "bob", Name: "bob"})
	database.DB = db
}

func newRouter() *gin.Engine {
	r := gin.New()
	api := r.Group("/api")
	api.Use(auth.UserID())

	exercises := api.Group("/exercises")
	exercises.GET("", ListExercises)
	exercises.POST("", CreateExercise)
	exercises.GET("/:id", GetExercise)
	exercises.PUT("/:id", UpdateExercise)
	exercises.DELETE("/:id", DeleteExercise)
	exercises.GET("/:id/versions", ListExerciseVersions)
	exercises.GET("/templates/:templateId/versions/:version", GetExerciseVersion)

	schemes := api.Group("/exercise-schemes")
	schemes.GET("", ListExerciseSchemes)
	schemes.POST("", CreateExerciseScheme)
	schemes.GET("/:id", GetExerciseScheme)
	schemes.PUT("/:id", UpdateExerciseScheme)
	schemes.DELETE("/:id", DeleteExerciseScheme)

	return r
}

func doJSON(r *gin.Engine, method, path string, body any) *httptest.ResponseRecorder {
	var reader io.Reader
	if body != nil {
		data, _ := json.Marshal(body)
		reader = bytes.NewReader(data)
	}
	req := httptest.NewRequest(method, path, reader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func doJSONAs(r *gin.Engine, method, path string, body any, userID string) *httptest.ResponseRecorder {
	var reader io.Reader
	if body != nil {
		data, _ := json.Marshal(body)
		reader = bytes.NewReader(data)
	}
	req := httptest.NewRequest(method, path, reader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("X-User-Id", userID)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func doRaw(r *gin.Engine, method, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

type paginatedJSON struct {
	Items  json.RawMessage `json:"items"`
	Total  int             `json:"total"`
	Limit  int             `json:"limit"`
	Offset int             `json:"offset"`
}

func closeDB(t *testing.T) {
	t.Helper()
	sqlDB, err := database.DB.DB()
	if err != nil {
		t.Fatal(err)
	}
	sqlDB.Close()
}
