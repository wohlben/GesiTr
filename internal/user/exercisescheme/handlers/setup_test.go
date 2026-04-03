package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"os"
	"testing"

	"gesitr/internal/auth"
	exerciseHandlers "gesitr/internal/compendium/exercise/handlers"
	exerciseModels "gesitr/internal/compendium/exercise/models"
	"gesitr/internal/database"
	"gesitr/internal/humaconfig"
	"gesitr/internal/user/exercisescheme/models"
	namePreferenceModels "gesitr/internal/user/namepreference/models"

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
		&exerciseModels.ExerciseEntity{},
		&exerciseModels.ExerciseForce{},
		&exerciseModels.ExerciseMuscle{},
		&exerciseModels.ExerciseMeasurementParadigm{},
		&exerciseModels.ExerciseInstruction{},
		&exerciseModels.ExerciseImage{},
		&exerciseModels.ExerciseName{},
		&exerciseModels.ExerciseEquipment{},
		&exerciseModels.ExerciseHistoryEntity{},
		&models.ExerciseSchemeEntity{},
		&models.ExerciseSchemeSectionItemEntity{},
		&namePreferenceModels.ExerciseNamePreference{},
	)
	database.DB = db
}

func newRouter() *gin.Engine {
	r := gin.New()
	api := r.Group("/api")
	api.Use(auth.UserID())

	humaAPI := humaconfig.NewAPI(r, api)
	exerciseHandlers.RegisterRoutes(humaAPI)
	RegisterRoutes(humaAPI)

	return r
}

func newExercisePayload(name string) map[string]any {
	return map[string]any{
		"names": []string{name, "Alt Name"}, "type": "STRENGTH",
		"technicalDifficulty": "beginner", "bodyWeightScaling": 0.5,
		"description": "test",
		"force":       []string{"PUSH"}, "primaryMuscles": []string{"CHEST"},
		"secondaryMuscles":              []string{"TRICEPS"},
		"suggestedMeasurementParadigms": []string{"REP_BASED"},
		"instructions":                  []string{"Step 1", "Step 2"},
		"images":                        []string{"/img/a.jpg"},
	}
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

func closeDB(t *testing.T) {
	t.Helper()
	sqlDB, err := database.DB.DB()
	if err != nil {
		t.Fatal(err)
	}
	sqlDB.Close()
}
