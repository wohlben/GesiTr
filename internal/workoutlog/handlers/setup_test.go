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
	"gesitr/internal/humaconfig"
	profilemodels "gesitr/internal/profile/models"
	exerciseloghandlers "gesitr/internal/user/exerciselog/handlers"
	exerciselogmodels "gesitr/internal/user/exerciselog/models"
	workouthandlers "gesitr/internal/workout/handlers"
	workoutmodels "gesitr/internal/workout/models"
	"gesitr/internal/workoutlog/models"

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
		&workoutmodels.WorkoutHistoryEntity{},
		&workoutmodels.WorkoutSectionEntity{},
		&workoutmodels.WorkoutSectionItemEntity{},
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

	humaAPI := humaconfig.NewAPI(r, api)
	exercisehandlers.RegisterRoutes(humaAPI)
	equipmenthandlers.RegisterRoutes(humaAPI)
	workouthandlers.RegisterRoutes(humaAPI)
	RegisterRoutes(humaAPI)
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
