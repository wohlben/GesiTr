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
	userequipmenthandlers "gesitr/internal/user/equipment/handlers"
	userequipmentmodels "gesitr/internal/user/equipment/models"
	userexercisehandlers "gesitr/internal/user/exercise/handlers"
	userexercisemodels "gesitr/internal/user/exercise/models"
	recordhandlers "gesitr/internal/user/record/handlers"
	recordmodels "gesitr/internal/user/record/models"
	workouthandlers "gesitr/internal/user/workout/handlers"
	workoutmodels "gesitr/internal/user/workout/models"
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
		&userexercisemodels.UserExerciseEntity{},
		&userexercisemodels.UserExerciseSchemeEntity{},
		&userequipmentmodels.UserEquipmentEntity{},
		&workoutmodels.WorkoutEntity{},
		&workoutmodels.WorkoutSectionEntity{},
		&workoutmodels.WorkoutSectionExerciseEntity{},
		&models.WorkoutLogEntity{},
		&models.WorkoutLogSectionEntity{},
		&models.WorkoutLogExerciseEntity{},
		&models.WorkoutLogExerciseSetEntity{},
		&recordmodels.UserRecordEntity{},
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

	equipment := api.Group("/equipment")
	equipment.GET("", userequipmenthandlers.ListUserEquipment)
	equipment.POST("", userequipmenthandlers.CreateUserEquipment)
	equipment.GET("/:id", userequipmenthandlers.GetUserEquipment)
	equipment.DELETE("/:id", userequipmenthandlers.DeleteUserEquipment)

	schemes := api.Group("/exercise-schemes")
	schemes.GET("", userexercisehandlers.ListUserExerciseSchemes)
	schemes.POST("", userexercisehandlers.CreateUserExerciseScheme)
	schemes.GET("/:id", userexercisehandlers.GetUserExerciseScheme)
	schemes.PUT("/:id", userexercisehandlers.UpdateUserExerciseScheme)
	schemes.DELETE("/:id", userexercisehandlers.DeleteUserExerciseScheme)

	workouts := api.Group("/workouts")
	workouts.GET("", workouthandlers.ListWorkouts)
	workouts.POST("", workouthandlers.CreateWorkout)
	workouts.GET("/:id", workouthandlers.GetWorkout)
	workouts.PUT("/:id", workouthandlers.UpdateWorkout)
	workouts.DELETE("/:id", workouthandlers.DeleteWorkout)

	sections := api.Group("/workout-sections")
	sections.GET("", workouthandlers.ListWorkoutSections)
	sections.POST("", workouthandlers.CreateWorkoutSection)
	sections.GET("/:id", workouthandlers.GetWorkoutSection)
	sections.DELETE("/:id", workouthandlers.DeleteWorkoutSection)

	sectionExercises := api.Group("/workout-section-exercises")
	sectionExercises.GET("", workouthandlers.ListWorkoutSectionExercises)
	sectionExercises.POST("", workouthandlers.CreateWorkoutSectionExercise)
	sectionExercises.DELETE("/:id", workouthandlers.DeleteWorkoutSectionExercise)

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
	records.GET("", recordhandlers.ListUserRecords)
	records.GET("/:id", recordhandlers.GetUserRecord)

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
