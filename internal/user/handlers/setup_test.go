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
	t.Setenv("AUTH_FALLBACK_USER", "alice")
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	db.AutoMigrate(
		&models.UserExerciseEntity{},
		&models.UserEquipmentEntity{},
		&models.UserExerciseSchemeEntity{},
		&models.WorkoutEntity{},
		&models.WorkoutSectionEntity{},
		&models.WorkoutSectionExerciseEntity{},
		&models.WorkoutLogEntity{},
		&models.WorkoutLogSectionEntity{},
		&models.WorkoutLogExerciseEntity{},
		&models.WorkoutLogExerciseSetEntity{},
		&models.UserRecordEntity{},
	)
	database.DB = db
}

func newRouter() *gin.Engine {
	r := gin.New()
	api := r.Group("/api/user")
	api.Use(auth.UserID())

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

	workoutLogs := api.Group("/workout-logs")
	workoutLogs.GET("", ListWorkoutLogs)
	workoutLogs.POST("", CreateWorkoutLog)
	workoutLogs.GET("/:id", GetWorkoutLog)
	workoutLogs.PUT("/:id", UpdateWorkoutLog)
	workoutLogs.DELETE("/:id", DeleteWorkoutLog)
	workoutLogs.POST("/:id/start", StartWorkoutLog)
	workoutLogs.POST("/:id/abandon", AbandonWorkoutLog)

	logSections := api.Group("/workout-log-sections")
	logSections.GET("", ListWorkoutLogSections)
	logSections.POST("", CreateWorkoutLogSection)
	logSections.GET("/:id", GetWorkoutLogSection)
	logSections.DELETE("/:id", DeleteWorkoutLogSection)

	logExercises := api.Group("/workout-log-exercises")
	logExercises.GET("", ListWorkoutLogExercises)
	logExercises.POST("", CreateWorkoutLogExercise)
	logExercises.PUT("/:id", UpdateWorkoutLogExercise)
	logExercises.DELETE("/:id", DeleteWorkoutLogExercise)

	logExerciseSets := api.Group("/workout-log-exercise-sets")
	logExerciseSets.GET("", ListWorkoutLogExerciseSets)
	logExerciseSets.POST("", CreateWorkoutLogExerciseSet)
	logExerciseSets.PUT("/:id", UpdateWorkoutLogExerciseSet)
	logExerciseSets.DELETE("/:id", DeleteWorkoutLogExerciseSet)

	records := api.Group("/records")
	records.GET("", ListUserRecords)
	records.GET("/:id", GetUserRecord)

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
