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
	equipmenthandlers "gesitr/internal/compendium/equipment/handlers"
	equipmentmodels "gesitr/internal/compendium/equipment/models"
	exercisehandlers "gesitr/internal/compendium/exercise/handlers"
	exercisemodels "gesitr/internal/compendium/exercise/models"
	workouthandlers "gesitr/internal/compendium/workout/handlers"
	workoutmodels "gesitr/internal/compendium/workout/models"
	workoutgrouphandlers "gesitr/internal/compendium/workoutgroup/handlers"
	workoutgroupmodels "gesitr/internal/compendium/workoutgroup/models"
	"gesitr/internal/database"
	"gesitr/internal/humaconfig"
	profilemodels "gesitr/internal/profile/models"
	exerciseloghandlers "gesitr/internal/user/exerciselog/handlers"
	exerciselogmodels "gesitr/internal/user/exerciselog/models"
	workoutloghandlers "gesitr/internal/user/workoutlog/handlers"
	"gesitr/internal/user/workoutlog/models"
	workoutschedulehandlers "gesitr/internal/user/workoutschedule/handlers"
	workoutschedulemodels "gesitr/internal/user/workoutschedule/models"

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
		&workoutgroupmodels.WorkoutGroupEntity{},
		&workoutgroupmodels.WorkoutGroupMembershipEntity{},
		&workoutschedulemodels.WorkoutScheduleEntity{},
		&workoutschedulemodels.SchedulePeriodEntity{},
		&workoutschedulemodels.ScheduleCommitmentEntity{},
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
	workoutloghandlers.RegisterRoutes(humaAPI)
	exerciseloghandlers.RegisterRoutes(humaAPI)
	workoutgrouphandlers.RegisterRoutes(humaAPI)
	workoutschedulehandlers.RegisterRoutes(humaAPI)

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

func doJSONLogAs(t *testing.T, r *gin.Engine, method, path string, body any, userID string) *httptest.ResponseRecorder {
	t.Helper()
	if body != nil {
		reqJSON, _ := json.MarshalIndent(body, "  ", "  ")
		t.Logf(">>> %s %s (as %s)\n  Request body:\n  %s", method, path, userID, reqJSON)
	} else {
		t.Logf(">>> %s %s (as %s, no body)", method, path, userID)
	}
	w := doJSONAs(r, method, path, body, userID)
	var pretty json.RawMessage
	if err := json.Unmarshal(w.Body.Bytes(), &pretty); err == nil {
		respJSON, _ := json.MarshalIndent(pretty, "  ", "  ")
		t.Logf("<<< %d\n  Response body:\n  %s", w.Code, respJSON)
	} else {
		t.Logf("<<< %d\n  Response body (raw): %s", w.Code, w.Body.String())
	}
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
