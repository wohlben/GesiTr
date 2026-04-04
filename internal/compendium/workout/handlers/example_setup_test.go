package handlers

import (
	"bytes"
	"net/http/httptest"
	"os"

	equipmentmodels "gesitr/internal/compendium/equipment/models"
	exercisemodels "gesitr/internal/compendium/exercise/models"
	ownershipgroupmodels "gesitr/internal/compendium/ownershipgroup/models"
	"gesitr/internal/compendium/workout/models"
	workoutgroupmodels "gesitr/internal/compendium/workoutgroup/models"
	"gesitr/internal/database"
	exerciseschememodels "gesitr/internal/user/exercisescheme/models"
	namePreferenceModels "gesitr/internal/user/namepreference/models"

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

func doRawAs(r *gin.Engine, method, path, body, userID string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Id", userID)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
