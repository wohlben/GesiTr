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
	equipmenthandlers "gesitr/internal/equipment/handlers"
	equipmentmodels "gesitr/internal/equipment/models"
	exercisehandlers "gesitr/internal/exercise/handlers"
	exercisemodels "gesitr/internal/exercise/models"
	profilemodels "gesitr/internal/profile/models"
	exerciseloghandlers "gesitr/internal/user/exerciselog/handlers"
	exerciselogmodels "gesitr/internal/user/exerciselog/models"
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
	db, err := gorm.Open(sqlite.Open("file::memory:?_foreign_keys=on"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatal(err)
	}
	db.AutoMigrate(
		&profilemodels.UserProfileEntity{},
		&exercisemodels.ExerciseEntity{},
		&exercisemodels.ExerciseForce{},
		&exercisemodels.ExerciseMuscle{},
		&exercisemodels.ExerciseMeasurementParadigm{},
		&exercisemodels.ExerciseInstruction{},
		&exercisemodels.ExerciseImage{},
		&exercisemodels.ExerciseAlternativeName{},
		&exercisemodels.ExerciseEquipment{},
		&exercisemodels.ExerciseHistoryEntity{},
		&exercisemodels.ExerciseSchemeEntity{},
		&equipmentmodels.EquipmentEntity{},
		&workoutmodels.WorkoutEntity{},
		&workoutmodels.WorkoutSectionEntity{},
		&workoutmodels.WorkoutSectionExerciseEntity{},
		&models.WorkoutLogEntity{},
		&models.WorkoutLogSectionEntity{},
		&models.WorkoutLogExerciseEntity{},
		&models.WorkoutLogExerciseSetEntity{},
		&exerciselogmodels.ExerciseLogEntity{},
	)
	db.Create(&profilemodels.UserProfileEntity{ID: "alice", Name: "alice"})
	db.Create(&profilemodels.UserProfileEntity{ID: "bob", Name: "bob"})
	database.DB = db
}

func newRouter() *gin.Engine {
	r := gin.New()
	api := r.Group("/api")
	api.Use(auth.UserID())

	exercises := api.Group("/exercises")
	exercises.GET("", exercisehandlers.ListExercises)
	exercises.POST("", exercisehandlers.CreateExercise)
	exercises.GET("/:id", exercisehandlers.GetExercise)
	exercises.DELETE("/:id", exercisehandlers.DeleteExercise)

	equipment := api.Group("/equipment")
	equipment.GET("", equipmenthandlers.ListEquipment)
	equipment.POST("", equipmenthandlers.CreateEquipment)
	equipment.GET("/:id", equipmenthandlers.GetEquipment)
	equipment.DELETE("/:id", equipmenthandlers.DeleteEquipment)

	schemes := api.Group("/exercise-schemes")
	schemes.GET("", exercisehandlers.ListExerciseSchemes)
	schemes.POST("", exercisehandlers.CreateExerciseScheme)
	schemes.GET("/:id", exercisehandlers.GetExerciseScheme)
	schemes.PUT("/:id", exercisehandlers.UpdateExerciseScheme)
	schemes.DELETE("/:id", exercisehandlers.DeleteExerciseScheme)

	user := api.Group("/user")

	workouts := user.Group("/workouts")
	workouts.GET("", workouthandlers.ListWorkouts)
	workouts.POST("", workouthandlers.CreateWorkout)
	workouts.GET("/:id", workouthandlers.GetWorkout)
	workouts.PUT("/:id", workouthandlers.UpdateWorkout)
	workouts.DELETE("/:id", workouthandlers.DeleteWorkout)

	sections := user.Group("/workout-sections")
	sections.GET("", workouthandlers.ListWorkoutSections)
	sections.POST("", workouthandlers.CreateWorkoutSection)
	sections.GET("/:id", workouthandlers.GetWorkoutSection)
	sections.DELETE("/:id", workouthandlers.DeleteWorkoutSection)

	sectionExercises := user.Group("/workout-section-exercises")
	sectionExercises.GET("", workouthandlers.ListWorkoutSectionExercises)
	sectionExercises.POST("", workouthandlers.CreateWorkoutSectionExercise)
	sectionExercises.DELETE("/:id", workouthandlers.DeleteWorkoutSectionExercise)

	workoutLogs := user.Group("/workout-logs")
	workoutLogs.GET("", ListWorkoutLogs)
	workoutLogs.POST("", CreateWorkoutLog)
	workoutLogs.GET("/:id", GetWorkoutLog)
	workoutLogs.PATCH("/:id", UpdateWorkoutLog)
	workoutLogs.DELETE("/:id", DeleteWorkoutLog)
	workoutLogs.POST("/adhoc", StartAdhocWorkoutLog)
	workoutLogs.POST("/:id/start", StartWorkoutLog)
	workoutLogs.POST("/:id/abandon", AbandonWorkoutLog)
	workoutLogs.POST("/:id/finish", FinishWorkoutLog)

	logSections := user.Group("/workout-log-sections")
	logSections.GET("", ListWorkoutLogSections)
	logSections.POST("", CreateWorkoutLogSection)
	logSections.GET("/:id", GetWorkoutLogSection)
	logSections.DELETE("/:id", DeleteWorkoutLogSection)

	logExercises := user.Group("/workout-log-exercises")
	logExercises.GET("", ListWorkoutLogExercises)
	logExercises.POST("", CreateWorkoutLogExercise)
	logExercises.PATCH("/:id", UpdateWorkoutLogExercise)
	logExercises.DELETE("/:id", DeleteWorkoutLogExercise)

	logExerciseSets := user.Group("/workout-log-exercise-sets")
	logExerciseSets.GET("", ListWorkoutLogExerciseSets)
	logExerciseSets.POST("", CreateWorkoutLogExerciseSet)
	logExerciseSets.PATCH("/:id", UpdateWorkoutLogExerciseSet)
	logExerciseSets.DELETE("/:id", DeleteWorkoutLogExerciseSet)

	exerciseLogs := user.Group("/exercise-logs")
	exerciseLogs.GET("", exerciseloghandlers.ListExerciseLogs)
	exerciseLogs.POST("", exerciseloghandlers.CreateExerciseLog)
	exerciseLogs.GET("/:id", exerciseloghandlers.GetExerciseLog)
	exerciseLogs.PATCH("/:id", exerciseloghandlers.UpdateExerciseLog)
	exerciseLogs.DELETE("/:id", exerciseloghandlers.DeleteExerciseLog)

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
