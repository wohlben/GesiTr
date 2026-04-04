package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"os"
	"testing"

	"gesitr/internal/auth"
	equipmentmodels "gesitr/internal/compendium/equipment/models"
	exercisehandlers "gesitr/internal/compendium/exercise/handlers"
	exercisemodels "gesitr/internal/compendium/exercise/models"
	"gesitr/internal/compendium/workout/models"
	workoutgrouphandlers "gesitr/internal/compendium/workoutgroup/handlers"
	workoutgroupmodels "gesitr/internal/compendium/workoutgroup/models"
	"gesitr/internal/database"
	"gesitr/internal/humaconfig"
	ownershipgroupmodels "gesitr/internal/compendium/ownershipgroup/models"
	exerciseschemehandlers "gesitr/internal/user/exercisescheme/handlers"
	exerciseschememodels "gesitr/internal/user/exercisescheme/models"
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
	t.Setenv("AUTH_FALLBACK_USER", "alice")
	db, err := gorm.Open(sqlite.Open("file::memory:?_foreign_keys=on"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatal(err)
	}
	db.AutoMigrate(
		&ownershipgroupmodels.OwnershipGroupEntity{},
		&ownershipgroupmodels.OwnershipGroupMembershipEntity{},
		&exercisemodels.ExerciseEntity{},
		&exercisemodels.ExerciseForce{},
		&exercisemodels.ExerciseMuscle{},
		&exercisemodels.ExerciseMeasurementParadigm{},
		&exercisemodels.ExerciseInstruction{},
		&exercisemodels.ExerciseImage{},
		&exercisemodels.ExerciseName{},
		&exercisemodels.ExerciseEquipment{},
		&exercisemodels.ExerciseHistoryEntity{},
		&exerciseschememodels.ExerciseSchemeEntity{},
		&exerciseschememodels.ExerciseSchemeSectionItemEntity{},
		&namePreferenceModels.ExerciseNamePreference{},
		&equipmentmodels.EquipmentEntity{},
		&models.ExerciseGroupEntity{},
		&models.ExerciseGroupMemberEntity{},
		&models.WorkoutEntity{},
		&models.WorkoutHistoryEntity{},
		&models.WorkoutSectionEntity{},
		&models.WorkoutSectionItemEntity{},
		&models.WorkoutRelationshipEntity{},
		&workoutgroupmodels.WorkoutGroupEntity{},
		&workoutgroupmodels.WorkoutGroupMembershipEntity{},
	)
	database.DB = db
}

func newRouter() *gin.Engine {
	r := gin.New()
	api := r.Group("/api")
	api.Use(auth.UserID())

	humaAPI := humaconfig.NewAPI(r, api)
	exercisehandlers.RegisterRoutes(humaAPI)
	exerciseschemehandlers.RegisterRoutes(humaAPI)
	RegisterRoutes(humaAPI)
	workoutgrouphandlers.RegisterRoutes(humaAPI)

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
