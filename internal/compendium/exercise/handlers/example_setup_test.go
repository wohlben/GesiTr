package handlers

import (
	"bytes"
	"net/http/httptest"
	"os"

	"gesitr/internal/auth"
	equipmenthandlers "gesitr/internal/compendium/equipment/handlers"
	equipmentmodels "gesitr/internal/compendium/equipment/models"
	"gesitr/internal/compendium/exercise/models"
	"gesitr/internal/database"
	"gesitr/internal/humaconfig"
	profilemodels "gesitr/internal/profile/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupExampleDB() {
	os.Setenv("AUTH_FALLBACK_USER", "testuser")
	db, err := gorm.Open(sqlite.Open("file::memory:?_foreign_keys=on"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(
		&profilemodels.UserProfileEntity{},
		&models.ExerciseEntity{},
		&models.ExerciseForce{},
		&models.ExerciseMuscle{},
		&models.ExerciseMeasurementParadigm{},
		&models.ExerciseInstruction{},
		&models.ExerciseImage{},
		&models.ExerciseAlternativeName{},
		&models.ExerciseEquipment{},
		&models.ExerciseHistoryEntity{},
		&models.ExerciseSchemeEntity{},
		&equipmentmodels.EquipmentEntity{},
		&equipmentmodels.EquipmentHistoryEntity{},
	)
	db.Create(&profilemodels.UserProfileEntity{ID: "testuser", Name: "testuser"})
	db.Create(&profilemodels.UserProfileEntity{ID: "other", Name: "other"})
	database.DB = db
}

// newExampleRouter registers both exercise and equipment routes, so examples
// can demonstrate cross-API flows.
func newExampleRouter() *gin.Engine {
	r := gin.New()
	api := r.Group("/api")
	api.Use(auth.UserID())

	humaAPI := humaconfig.NewAPI(r, api)
	RegisterRoutes(humaAPI)
	equipmenthandlers.RegisterRoutes(humaAPI)

	return r
}

func doRawAs(r *gin.Engine, method, path, body, userID string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Id", userID)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
