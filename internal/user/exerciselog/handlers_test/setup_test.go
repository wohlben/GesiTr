package handlers_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"os"
	"testing"

	"gesitr/internal/auth"
	equipmenthandlers "gesitr/internal/compendium/equipment/handlers"
	equipmentmodels "gesitr/internal/compendium/equipment/models"
	exercisehandlers "gesitr/internal/compendium/exercise/handlers"
	exercisemodels "gesitr/internal/compendium/exercise/models"
	workouthandlers "gesitr/internal/compendium/workout/handlers"
	workoutmodels "gesitr/internal/compendium/workout/models"
	"gesitr/internal/database"
	"gesitr/internal/humaconfig"
	exerciseloghandlers "gesitr/internal/user/exerciselog/handlers"
	exerciselogmodels "gesitr/internal/user/exerciselog/models"
	namePreferenceModels "gesitr/internal/user/namepreference/models"
	workoutloghandlers "gesitr/internal/user/workoutlog/handlers"
	workoutlogmodels "gesitr/internal/user/workoutlog/models"

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
		&exercisemodels.ExerciseEntity{},
		&exercisemodels.ExerciseForce{},
		&exercisemodels.ExerciseMuscle{},
		&exercisemodels.ExerciseMeasurementParadigm{},
		&exercisemodels.ExerciseInstruction{},
		&exercisemodels.ExerciseImage{},
		&exercisemodels.ExerciseName{},
		&exercisemodels.ExerciseEquipment{},
		&exercisemodels.ExerciseHistoryEntity{},
		&exercisemodels.ExerciseSchemeEntity{},
		&namePreferenceModels.ExerciseNamePreference{},
		&equipmentmodels.EquipmentEntity{},
		&workoutmodels.WorkoutEntity{},
		&workoutmodels.WorkoutHistoryEntity{},
		&workoutmodels.WorkoutSectionEntity{},
		&workoutmodels.WorkoutSectionItemEntity{},
		&workoutlogmodels.WorkoutLogEntity{},
		&workoutlogmodels.WorkoutLogSectionEntity{},
		&workoutlogmodels.WorkoutLogExerciseEntity{},
		&workoutlogmodels.WorkoutLogExerciseSetEntity{},
		&exerciselogmodels.ExerciseLogEntity{},
	)
	database.DB = db
}

func newRouter() *gin.Engine {
	r := gin.New()
	api := r.Group("/api")
	api.Use(auth.UserID())

	humaAPI := humaconfig.NewAPI(r, api)
	exercisehandlers.RegisterRoutes(humaAPI)
	equipmenthandlers.RegisterRoutes(humaAPI)
	workouthandlers.RegisterRoutes(humaAPI)
	workoutloghandlers.RegisterRoutes(humaAPI)
	exerciseloghandlers.RegisterRoutes(humaAPI)

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

func closeDB(t *testing.T) {
	t.Helper()
	sqlDB, err := database.DB.DB()
	if err != nil {
		t.Fatal(err)
	}
	sqlDB.Close()
}
