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
	userexercise "gesitr/internal/user/exercise"
	"gesitr/internal/user/record"
	"gesitr/internal/user/workout"
	"gesitr/internal/user/workoutlog/models"

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
		&userexercise.UserExerciseEntity{},
		&userexercise.UserEquipmentEntity{},
		&userexercise.UserExerciseSchemeEntity{},
		&workout.WorkoutEntity{},
		&workout.WorkoutSectionEntity{},
		&workout.WorkoutSectionExerciseEntity{},
		&models.WorkoutLogEntity{},
		&models.WorkoutLogSectionEntity{},
		&models.WorkoutLogExerciseEntity{},
		&models.WorkoutLogExerciseSetEntity{},
		&record.UserRecordEntity{},
	)
	database.DB = db
}

func newRouter() *gin.Engine {
	r := gin.New()
	api := r.Group("/api/user")
	api.Use(auth.UserID())

	exercises := api.Group("/exercises")
	exercises.GET("", userexercise.ListUserExercises)
	exercises.POST("", userexercise.CreateUserExercise)
	exercises.GET("/:id", userexercise.GetUserExercise)
	exercises.DELETE("/:id", userexercise.DeleteUserExercise)

	equipment := api.Group("/equipment")
	equipment.GET("", userexercise.ListUserEquipment)
	equipment.POST("", userexercise.CreateUserEquipment)
	equipment.GET("/:id", userexercise.GetUserEquipment)
	equipment.DELETE("/:id", userexercise.DeleteUserEquipment)

	schemes := api.Group("/exercise-schemes")
	schemes.GET("", userexercise.ListUserExerciseSchemes)
	schemes.POST("", userexercise.CreateUserExerciseScheme)
	schemes.GET("/:id", userexercise.GetUserExerciseScheme)
	schemes.PUT("/:id", userexercise.UpdateUserExerciseScheme)
	schemes.DELETE("/:id", userexercise.DeleteUserExerciseScheme)

	workouts := api.Group("/workouts")
	workouts.GET("", workout.ListWorkouts)
	workouts.POST("", workout.CreateWorkout)
	workouts.GET("/:id", workout.GetWorkout)
	workouts.PUT("/:id", workout.UpdateWorkout)
	workouts.DELETE("/:id", workout.DeleteWorkout)

	sections := api.Group("/workout-sections")
	sections.GET("", workout.ListWorkoutSections)
	sections.POST("", workout.CreateWorkoutSection)
	sections.GET("/:id", workout.GetWorkoutSection)
	sections.DELETE("/:id", workout.DeleteWorkoutSection)

	sectionExercises := api.Group("/workout-section-exercises")
	sectionExercises.GET("", workout.ListWorkoutSectionExercises)
	sectionExercises.POST("", workout.CreateWorkoutSectionExercise)
	sectionExercises.DELETE("/:id", workout.DeleteWorkoutSectionExercise)

	workoutLogs := api.Group("/workout-logs")
	workoutLogs.GET("", ListWorkoutLogs)
	workoutLogs.POST("", CreateWorkoutLog)
	workoutLogs.GET("/:id", GetWorkoutLog)
	workoutLogs.PATCH("/:id", UpdateWorkoutLog)
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
	logExercises.PATCH("/:id", UpdateWorkoutLogExercise)
	logExercises.DELETE("/:id", DeleteWorkoutLogExercise)

	logExerciseSets := api.Group("/workout-log-exercise-sets")
	logExerciseSets.GET("", ListWorkoutLogExerciseSets)
	logExerciseSets.POST("", CreateWorkoutLogExerciseSet)
	logExerciseSets.PATCH("/:id", UpdateWorkoutLogExerciseSet)
	logExerciseSets.DELETE("/:id", DeleteWorkoutLogExerciseSet)

	records := api.Group("/records")
	records.GET("", record.ListUserRecords)
	records.GET("/:id", record.GetUserRecord)

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
