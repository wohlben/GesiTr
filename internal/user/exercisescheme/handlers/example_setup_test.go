package handlers

import (
	"bytes"
	"net/http/httptest"
	"os"

	"gesitr/internal/auth"
	exerciseHandlers "gesitr/internal/compendium/exercise/handlers"
	exerciseModels "gesitr/internal/compendium/exercise/models"
	"gesitr/internal/database"
	"gesitr/internal/humaconfig"
	"gesitr/internal/user/exercisescheme/models"
	namePreferenceModels "gesitr/internal/user/namepreference/models"

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
		&exerciseModels.ExerciseEntity{},
		&exerciseModels.ExerciseForce{},
		&exerciseModels.ExerciseMuscle{},
		&exerciseModels.ExerciseMeasurementParadigm{},
		&exerciseModels.ExerciseInstruction{},
		&exerciseModels.ExerciseImage{},
		&exerciseModels.ExerciseName{},
		&exerciseModels.ExerciseEquipment{},
		&exerciseModels.ExerciseHistoryEntity{},
		&models.ExerciseSchemeEntity{},
		&namePreferenceModels.ExerciseNamePreference{},
	)
	database.DB = db
}

func newExampleRouter() *gin.Engine {
	r := gin.New()
	api := r.Group("/api")
	api.Use(auth.UserID())

	humaAPI := humaconfig.NewAPI(r, api)
	exerciseHandlers.RegisterRoutes(humaAPI)
	RegisterRoutes(humaAPI)

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
