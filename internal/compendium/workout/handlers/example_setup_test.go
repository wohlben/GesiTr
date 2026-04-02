package handlers

import (
	"bytes"
	"net/http/httptest"
	"os"

	equipmentmodels "gesitr/internal/compendium/equipment/models"
	exercisemodels "gesitr/internal/compendium/exercise/models"
	"gesitr/internal/compendium/workout/models"
	workoutgroupmodels "gesitr/internal/compendium/workoutgroup/models"
	"gesitr/internal/database"
	profilemodels "gesitr/internal/profile/models"

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
		&models.ExerciseGroupEntity{},
		&models.ExerciseGroupMemberEntity{},
		&models.WorkoutEntity{},
		&models.WorkoutHistoryEntity{},
		&models.WorkoutSectionEntity{},
		&models.WorkoutSectionItemEntity{},
		&workoutgroupmodels.WorkoutGroupEntity{},
		&workoutgroupmodels.WorkoutGroupMembershipEntity{},
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
