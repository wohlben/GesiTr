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
	workoutGroupHandlers "gesitr/internal/user/workoutgroup/handlers"
	workoutGroupModels "gesitr/internal/user/workoutgroup/models"
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

	// Drop the template_id column from equipment if it exists.
	var equipTemplateCount int64
	database.DB.Raw("SELECT COUNT(*) FROM pragma_table_info('equipment') WHERE name = 'template_id'").Scan(&equipTemplateCount)
	if equipTemplateCount > 0 {
		database.DB.Exec("DROP INDEX IF EXISTS idx_equip_owner_template")
		database.DB.Exec("ALTER TABLE equipment DROP COLUMN template_id")
	}

	// Drop the template_id column from exercise_groups if it exists.
	var groupTemplateCount int64
	database.DB.Raw("SELECT COUNT(*) FROM pragma_table_info('exercise_groups') WHERE name = 'template_id'").Scan(&groupTemplateCount)
	if groupTemplateCount > 0 {
		database.DB.Exec("ALTER TABLE exercise_groups DROP COLUMN template_id")
	}

	// Drop the description column from exercise_groups if it exists.
	var groupDescCount int64
	database.DB.Raw("SELECT COUNT(*) FROM pragma_table_info('exercise_groups') WHERE name = 'description'").Scan(&groupDescCount)
	if groupDescCount > 0 {
		database.DB.Exec("ALTER TABLE exercise_groups DROP COLUMN description")
	}

	// Backfill exercise_schemes.workout_section_item_id from workout_section_items.exercise_scheme_id.
	// This populates the reverse FK for existing data after adding the new column.
	database.DB.Exec(`UPDATE exercise_schemes SET workout_section_item_id = (
		SELECT id FROM workout_section_items WHERE exercise_scheme_id = exercise_schemes.id AND deleted_at IS NULL LIMIT 1
	) WHERE workout_section_item_id IS NULL AND id IN (
		SELECT exercise_scheme_id FROM workout_section_items WHERE exercise_scheme_id IS NOT NULL AND deleted_at IS NULL
	)`)

	// Migrate workout_section_exercises → workout_section_items if old table still exists.
	var oldTableCount int64
	database.DB.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='workout_section_exercises'").Scan(&oldTableCount)
	if oldTableCount > 0 {
		database.DB.Exec(`INSERT INTO workout_section_items (id, created_at, updated_at, deleted_at, workout_section_id, type, exercise_scheme_id, position)
			SELECT id, created_at, updated_at, deleted_at, workout_section_id, 'exercise', exercise_scheme_id, position FROM workout_section_exercises`)
		database.DB.Exec("DROP TABLE workout_section_exercises")
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
		&workoutModels.WorkoutSectionItemEntity{},
		&workoutlogmodels.WorkoutLogEntity{},
		&workoutlogmodels.WorkoutLogSectionEntity{},
		&workoutlogmodels.WorkoutLogExerciseEntity{},
		&workoutlogmodels.WorkoutLogExerciseSetEntity{},
		&exerciseLogModels.ExerciseLogEntity{},
		&workoutGroupModels.WorkoutGroupEntity{},
		&workoutGroupModels.WorkoutGroupMembershipEntity{},
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
	workoutGroupHandlers.RegisterRoutes(humaAPI)
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
			profile.ResetProfileCache()
			c.JSON(http.StatusOK, gin.H{"status": "reset"})
		})
	}
	setupSPA(r)
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
