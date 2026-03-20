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
	userequipmentmodels "gesitr/internal/user/equipment/models"
	userexercisehandlers "gesitr/internal/user/exercise/handlers"
	userexercisemodels "gesitr/internal/user/exercise/models"
	"gesitr/internal/user/workout/models"

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
	t.Setenv("AUTH_FALLBACK_USER", "alice")
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatal(err)
	}
	db.AutoMigrate(
		&userexercisemodels.UserExerciseEntity{},
		&userexercisemodels.UserExerciseSchemeEntity{},
		&userequipmentmodels.UserEquipmentEntity{},
		&models.WorkoutEntity{},
		&models.WorkoutSectionEntity{},
		&models.WorkoutSectionExerciseEntity{},
	)
	database.DB = db
}

func newRouter() *gin.Engine {
	r := gin.New()
	api := r.Group("/api/user")
	api.Use(auth.UserID())

	exercises := api.Group("/exercises")
	exercises.GET("", userexercisehandlers.ListUserExercises)
	exercises.POST("", userexercisehandlers.CreateUserExercise)
	exercises.GET("/:id", userexercisehandlers.GetUserExercise)
	exercises.DELETE("/:id", userexercisehandlers.DeleteUserExercise)

	schemes := api.Group("/exercise-schemes")
	schemes.GET("", userexercisehandlers.ListUserExerciseSchemes)
	schemes.POST("", userexercisehandlers.CreateUserExerciseScheme)
	schemes.GET("/:id", userexercisehandlers.GetUserExerciseScheme)
	schemes.PUT("/:id", userexercisehandlers.UpdateUserExerciseScheme)
	schemes.DELETE("/:id", userexercisehandlers.DeleteUserExerciseScheme)

	workouts := api.Group("/workouts")
	workouts.GET("", ListWorkouts)
	workouts.POST("", CreateWorkout)
	workouts.GET("/:id", GetWorkout)
	workouts.PUT("/:id", UpdateWorkout)
	workouts.DELETE("/:id", DeleteWorkout)

	sections := api.Group("/workout-sections")
	sections.GET("", ListWorkoutSections)
	sections.POST("", CreateWorkoutSection)
	sections.GET("/:id", GetWorkoutSection)
	sections.DELETE("/:id", DeleteWorkoutSection)

	sectionExercises := api.Group("/workout-section-exercises")
	sectionExercises.GET("", ListWorkoutSectionExercises)
	sectionExercises.POST("", CreateWorkoutSectionExercise)
	sectionExercises.DELETE("/:id", DeleteWorkoutSectionExercise)

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
