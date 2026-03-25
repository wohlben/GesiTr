package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"gesitr/internal/auth"
	"gesitr/internal/database"
	"gesitr/internal/docs"
	equipmentHandlers "gesitr/internal/equipment/handlers"
	equipmentModels "gesitr/internal/equipment/models"
	fulfillmentHandlers "gesitr/internal/equipmentfulfillment/handlers"
	fulfillmentModels "gesitr/internal/equipmentfulfillment/models"
	equipmentRelHandlers "gesitr/internal/equipmentrelationship/handlers"
	equipmentRelModels "gesitr/internal/equipmentrelationship/models"
	exerciseHandlers "gesitr/internal/exercise/handlers"
	exerciseModels "gesitr/internal/exercise/models"
	exerciseGroupHandlers "gesitr/internal/exercisegroup/handlers"
	exerciseGroupModels "gesitr/internal/exercisegroup/models"
	exerciseRelHandlers "gesitr/internal/exerciserelationship/handlers"
	exerciseRelModels "gesitr/internal/exerciserelationship/models"
	"gesitr/internal/humaconfig"
	"gesitr/internal/profile"
	profileHandlers "gesitr/internal/profile/handlers"
	profileModels "gesitr/internal/profile/models"
	exerciseLogHandlers "gesitr/internal/user/exerciselog/handlers"
	exerciseLogModels "gesitr/internal/user/exerciselog/models"
	workoutHandlers "gesitr/internal/user/workout/handlers"
	workoutModels "gesitr/internal/user/workout/models"
	workoutloghandlers "gesitr/internal/user/workoutlog/handlers"
	workoutlogmodels "gesitr/internal/user/workoutlog/models"

	"github.com/gin-gonic/gin"
)

//go:embed web/dist/browser/*
var staticFiles embed.FS

//go:embed docs/generated/*
var docsFiles embed.FS

// runMigrations handles manual schema changes that AutoMigrate cannot
// (e.g., dropping columns or indexes).
func runMigrations() {
	// Drop the slug column from exercises if it exists (removed in resource-permissions feature).
	// AutoMigrate does not drop columns or indexes, so we do it manually.
	var count int64
	database.DB.Raw("SELECT COUNT(*) FROM pragma_table_info('exercises') WHERE name = 'slug'").Scan(&count)
	if count > 0 {
		database.DB.Exec("DROP INDEX IF EXISTS idx_owner_slug")
		database.DB.Exec("ALTER TABLE exercises DROP COLUMN slug")
	}

	// Drop the template_id column from exercises if it exists (replaced by forked relationships).
	var templateCount int64
	database.DB.Raw("SELECT COUNT(*) FROM pragma_table_info('exercises') WHERE name = 'template_id'").Scan(&templateCount)
	if templateCount > 0 {
		database.DB.Exec("DROP INDEX IF EXISTS idx_owner_template")
		database.DB.Exec("ALTER TABLE exercises DROP COLUMN template_id")
	}
}

func autoMigrate() {
	database.DB.AutoMigrate(
		&profileModels.UserProfileEntity{},
		&exerciseModels.ExerciseEntity{},
		&exerciseModels.ExerciseForce{},
		&exerciseModels.ExerciseMuscle{},
		&exerciseModels.ExerciseMeasurementParadigm{},
		&exerciseModels.ExerciseInstruction{},
		&exerciseModels.ExerciseImage{},
		&exerciseModels.ExerciseAlternativeName{},
		&equipmentModels.EquipmentEntity{},
		&exerciseModels.ExerciseEquipment{},
		&fulfillmentModels.FulfillmentEntity{},
		&exerciseRelModels.ExerciseRelationshipEntity{},
		&exerciseGroupModels.ExerciseGroupEntity{},
		&exerciseGroupModels.ExerciseGroupMemberEntity{},
		&exerciseModels.ExerciseHistoryEntity{},
		&equipmentModels.EquipmentHistoryEntity{},
		&exerciseModels.ExerciseSchemeEntity{},
		&equipmentRelModels.EquipmentRelationshipEntity{},
		&workoutModels.WorkoutEntity{},
		&workoutModels.WorkoutSectionEntity{},
		&workoutModels.WorkoutSectionExerciseEntity{},
		&workoutlogmodels.WorkoutLogEntity{},
		&workoutlogmodels.WorkoutLogSectionEntity{},
		&workoutlogmodels.WorkoutLogExerciseEntity{},
		&workoutlogmodels.WorkoutLogExerciseSetEntity{},
		&exerciseLogModels.ExerciseLogEntity{},
	)
}

func setupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	api.Use(auth.UserID())
	api.Use(profile.EnsureProfile())

	// Huma API — shares the /api group so Gin auth/profile middleware applies.
	humaAPI := humaconfig.NewAPI(r, api)
	exerciseHandlers.RegisterRoutes(humaAPI)
	equipmentHandlers.RegisterRoutes(humaAPI)
	fulfillmentHandlers.RegisterRoutes(humaAPI)
	exerciseRelHandlers.RegisterRoutes(humaAPI)
	equipmentRelHandlers.RegisterRoutes(humaAPI)
	exerciseGroupHandlers.RegisterRoutes(humaAPI)
	profileHandlers.RegisterRoutes(humaAPI)
	workoutHandlers.RegisterRoutes(humaAPI)
	workoutloghandlers.RegisterRoutes(humaAPI)
	exerciseLogHandlers.RegisterRoutes(humaAPI)
}

func setupSPA(r *gin.Engine) {
	distFS, err := fs.Sub(staticFiles, "web/dist/browser")
	if err != nil {
		log.Fatal("Failed to load embedded files:", err)
	}
	indexHTML, err := fs.ReadFile(distFS, "index.html")
	if err != nil {
		log.Fatal("Failed to read index.html:", err)
	}
	r.NoRoute(func(c *gin.Context) {
		f, err := http.FS(distFS).Open(c.Request.URL.Path)
		if err == nil {
			f.Close()
			c.FileFromFS(c.Request.URL.Path, http.FS(distFS))
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", indexHTML)
	})
}

func buildApp() *gin.Engine {
	database.Init()
	autoMigrate()
	runMigrations()

	r := gin.Default()
	setupRoutes(r)
	docsFS, _ := fs.Sub(docsFiles, "docs/generated")
	docs.SetupRoutes(r, docsFS)

	if os.Getenv("DEV") == "true" {
		r.POST("/api/ci/reset-db", func(c *gin.Context) {
			database.DB.Exec("PRAGMA foreign_keys = OFF")
			var tables []struct{ Name string }
			database.DB.Raw("SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'").Scan(&tables)
			for _, t := range tables {
				database.DB.Exec("DELETE FROM " + t.Name)
			}
			database.DB.Exec("DELETE FROM sqlite_sequence")
			database.DB.Exec("PRAGMA foreign_keys = ON")
			c.JSON(http.StatusOK, gin.H{"status": "reset"})
		})
	} else {
		setupSPA(r)
	}
	return r
}

func main() {
	r := buildApp()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
