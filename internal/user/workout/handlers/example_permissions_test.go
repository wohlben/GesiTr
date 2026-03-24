package handlers

import (
	"bytes"
	"fmt"
	"net/http/httptest"
	"os"

	"gesitr/internal/database"
	equipmentmodels "gesitr/internal/equipment/models"
	exercisemodels "gesitr/internal/exercise/models"
	profilemodels "gesitr/internal/profile/models"
	"gesitr/internal/user/workout/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupExampleDB() {
	os.Setenv("AUTH_FALLBACK_USER", "alice")
	db, err := gorm.Open(sqlite.Open("file::memory:?_foreign_keys=on"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		panic(err)
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
		&models.WorkoutEntity{},
		&models.WorkoutSectionEntity{},
		&models.WorkoutSectionExerciseEntity{},
	)
	db.Create(&profilemodels.UserProfileEntity{ID: "alice", Name: "alice"})
	db.Create(&profilemodels.UserProfileEntity{ID: "bob", Name: "bob"})
	database.DB = db
}

func doRawAs(r *gin.Engine, method, path, body, userID string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Id", userID)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// Owner can retrieve their own workout.
func ExampleGetWorkout_ownerAccess() {
	setupExampleDB()
	r := newRouter()

	// Create a workout as alice (the AUTH_FALLBACK_USER).
	doRaw(r, "POST", "/api/user/workouts", `{
		"name": "Push Day"
	}`)

	// Retrieve it as the owner.
	w := doJSON(r, "GET", "/api/user/workouts/1", nil)
	fmt.Println(w.Code)
	// Output: 200
}

// Non-owner is denied access to another user's workout with 403 Forbidden.
func ExampleGetWorkout_nonOwnerDenied() {
	setupExampleDB()
	r := newRouter()

	// Create a workout as alice.
	doRaw(r, "POST", "/api/user/workouts", `{
		"name": "Push Day"
	}`)

	// Try to retrieve it as bob.
	w := doRawAs(r, "GET", "/api/user/workouts/1", "", "bob")
	fmt.Println(w.Code)
	// Output: 403
}
