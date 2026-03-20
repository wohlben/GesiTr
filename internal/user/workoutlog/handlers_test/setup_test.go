package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"os"
	"testing"

	"gesitr/internal/auth"
	"gesitr/internal/database"
	userexercise "gesitr/internal/user/exercise"
	"gesitr/internal/user/record"
	"gesitr/internal/user/workout"
	workoutloghandlers "gesitr/internal/user/workoutlog/handlers"
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
	workoutLogs.GET("", workoutloghandlers.ListWorkoutLogs)
	workoutLogs.POST("", workoutloghandlers.CreateWorkoutLog)
	workoutLogs.GET("/:id", workoutloghandlers.GetWorkoutLog)
	workoutLogs.PATCH("/:id", workoutloghandlers.UpdateWorkoutLog)
	workoutLogs.DELETE("/:id", workoutloghandlers.DeleteWorkoutLog)
	workoutLogs.POST("/:id/start", workoutloghandlers.StartWorkoutLog)
	workoutLogs.POST("/:id/abandon", workoutloghandlers.AbandonWorkoutLog)

	logSections := api.Group("/workout-log-sections")
	logSections.GET("", workoutloghandlers.ListWorkoutLogSections)
	logSections.POST("", workoutloghandlers.CreateWorkoutLogSection)
	logSections.GET("/:id", workoutloghandlers.GetWorkoutLogSection)
	logSections.DELETE("/:id", workoutloghandlers.DeleteWorkoutLogSection)

	logExercises := api.Group("/workout-log-exercises")
	logExercises.GET("", workoutloghandlers.ListWorkoutLogExercises)
	logExercises.POST("", workoutloghandlers.CreateWorkoutLogExercise)
	logExercises.PATCH("/:id", workoutloghandlers.UpdateWorkoutLogExercise)
	logExercises.DELETE("/:id", workoutloghandlers.DeleteWorkoutLogExercise)

	logExerciseSets := api.Group("/workout-log-exercise-sets")
	logExerciseSets.GET("", workoutloghandlers.ListWorkoutLogExerciseSets)
	logExerciseSets.POST("", workoutloghandlers.CreateWorkoutLogExerciseSet)
	logExerciseSets.PATCH("/:id", workoutloghandlers.UpdateWorkoutLogExerciseSet)
	logExerciseSets.DELETE("/:id", workoutloghandlers.DeleteWorkoutLogExerciseSet)

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

func itoa(id uint) string {
	return fmt.Sprintf("%d", id)
}

// doJSONLog wraps doJSON and logs the request body and response with pretty-printed JSON.
func doJSONLog(t *testing.T, r *gin.Engine, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()

	if body != nil {
		reqJSON, _ := json.MarshalIndent(body, "  ", "  ")
		t.Logf(">>> %s %s\n  Request body:\n  %s", method, path, reqJSON)
	} else {
		t.Logf(">>> %s %s (no body)", method, path)
	}

	w := doJSON(r, method, path, body)

	var pretty json.RawMessage
	if err := json.Unmarshal(w.Body.Bytes(), &pretty); err == nil {
		respJSON, _ := json.MarshalIndent(pretty, "  ", "  ")
		t.Logf("<<< %d\n  Response body:\n  %s", w.Code, respJSON)
	} else {
		t.Logf("<<< %d\n  Response body (raw): %s", w.Code, w.Body.String())
	}

	return w
}
