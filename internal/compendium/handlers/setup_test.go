package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"os"
	"testing"

	"gesitr/internal/compendium/models"
	"gesitr/internal/database"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func setupTestDB(t *testing.T) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	db.AutoMigrate(
		&models.ExerciseEntity{},
		&models.ExerciseForce{},
		&models.ExerciseMuscle{},
		&models.ExerciseMeasurementParadigm{},
		&models.ExerciseInstruction{},
		&models.ExerciseImage{},
		&models.ExerciseAlternativeName{},
		&models.EquipmentEntity{},
		&models.ExerciseEquipment{},
		&models.FulfillmentEntity{},
		&models.ExerciseRelationshipEntity{},
		&models.ExerciseGroupEntity{},
		&models.ExerciseGroupMemberEntity{},
	)
	database.DB = db
}

func newRouter() *gin.Engine {
	r := gin.New()
	api := r.Group("/api")

	exercises := api.Group("/exercises")
	exercises.GET("", ListExercises)
	exercises.POST("", CreateExercise)
	exercises.GET("/:id", GetExercise)
	exercises.PUT("/:id", UpdateExercise)
	exercises.DELETE("/:id", DeleteExercise)

	equipment := api.Group("/equipment")
	equipment.GET("", ListEquipment)
	equipment.POST("", CreateEquipment)
	equipment.GET("/:id", GetEquipment)
	equipment.PUT("/:id", UpdateEquipment)
	equipment.DELETE("/:id", DeleteEquipment)

	fulfillments := api.Group("/fulfillments")
	fulfillments.GET("", ListFulfillments)
	fulfillments.POST("", CreateFulfillment)
	fulfillments.DELETE("/:id", DeleteFulfillment)

	rels := api.Group("/exercise-relationships")
	rels.GET("", ListExerciseRelationships)
	rels.POST("", CreateExerciseRelationship)
	rels.DELETE("/:id", DeleteExerciseRelationship)

	groups := api.Group("/exercise-groups")
	groups.GET("", ListExerciseGroups)
	groups.POST("", CreateExerciseGroup)
	groups.GET("/:id", GetExerciseGroup)
	groups.DELETE("/:id", DeleteExerciseGroup)

	members := api.Group("/exercise-group-members")
	members.GET("", ListExerciseGroupMembers)
	members.POST("", CreateExerciseGroupMember)
	members.DELETE("/:id", DeleteExerciseGroupMember)

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
