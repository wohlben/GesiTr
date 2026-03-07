package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"os"
	"testing"

	"gesitr/internal/database"
	"gesitr/internal/user/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func setupTestDB(t *testing.T) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	db.AutoMigrate(
		&models.UserExerciseEntity{},
		&models.UserEquipmentEntity{},
		&models.UserExerciseSchemeEntity{},
	)
	database.DB = db
}

func newRouter() *gin.Engine {
	r := gin.New()
	api := r.Group("/api/user")

	exercises := api.Group("/exercises")
	exercises.GET("", ListUserExercises)
	exercises.POST("", CreateUserExercise)
	exercises.GET("/:id", GetUserExercise)
	exercises.DELETE("/:id", DeleteUserExercise)

	equipment := api.Group("/equipment")
	equipment.GET("", ListUserEquipment)
	equipment.POST("", CreateUserEquipment)
	equipment.GET("/:id", GetUserEquipment)
	equipment.DELETE("/:id", DeleteUserEquipment)

	schemes := api.Group("/exercise-schemes")
	schemes.GET("", ListUserExerciseSchemes)
	schemes.POST("", CreateUserExerciseScheme)
	schemes.GET("/:id", GetUserExerciseScheme)
	schemes.PUT("/:id", UpdateUserExerciseScheme)
	schemes.DELETE("/:id", DeleteUserExerciseScheme)

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
