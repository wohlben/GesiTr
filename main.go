package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

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
	masteryHandlers "gesitr/internal/user/mastery/handlers"
	masteryModels "gesitr/internal/user/mastery/models"
	workoutHandlers "gesitr/internal/user/workout/handlers"
	workoutModels "gesitr/internal/user/workout/models"
	workoutGroupHandlers "gesitr/internal/user/workoutgroup/handlers"
	workoutGroupModels "gesitr/internal/user/workoutgroup/models"
	workoutlog "gesitr/internal/user/workoutlog"
	workoutloghandlers "gesitr/internal/user/workoutlog/handlers"
	workoutlogmodels "gesitr/internal/user/workoutlog/models"
	workoutScheduleHandlers "gesitr/internal/user/workoutschedule/handlers"
	workoutScheduleModels "gesitr/internal/user/workoutschedule/models"

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

	// Drop old schedule config columns (replaced by period+commitment model).
	for _, col := range []string{"active", "days_of_week", "interval_weeks", "count", "period_days", "interval_days", "required_count", "type"} {
		var colCount int64
		database.DB.Raw("SELECT COUNT(*) FROM pragma_table_info('workout_schedules') WHERE name = ?", col).Scan(&colCount)
		if colCount > 0 {
			database.DB.Exec("ALTER TABLE workout_schedules DROP COLUMN " + col)
		}
	}
	// Drop required_count from schedule_periods if it exists.
	var reqCountCol int64
	database.DB.Raw("SELECT COUNT(*) FROM pragma_table_info('schedule_periods') WHERE name = 'required_count'").Scan(&reqCountCol)
	if reqCountCol > 0 {
		database.DB.Exec("ALTER TABLE schedule_periods DROP COLUMN required_count")
	}

	// Migrate workout_section_exercises → workout_section_items if old table still exists.
	var oldTableCount int64
	database.DB.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='workout_section_exercises'").Scan(&oldTableCount)
	if oldTableCount > 0 {
		database.DB.Exec(`INSERT INTO workout_section_items (id, created_at, updated_at, deleted_at, workout_section_id, type, exercise_scheme_id, position)
			SELECT id, created_at, updated_at, deleted_at, workout_section_id, 'exercise', exercise_scheme_id, position FROM workout_section_exercises`)
		database.DB.Exec("DROP TABLE workout_section_exercises")
	}

	// Backfill mastery_contributions from existing exercise_relationships.
	var contribCount int64
	database.DB.Raw("SELECT COUNT(*) FROM mastery_contributions").Scan(&contribCount)
	if contribCount == 0 {
		var relCount int64
		database.DB.Raw("SELECT COUNT(*) FROM exercise_relationships WHERE deleted_at IS NULL").Scan(&relCount)
		if relCount > 0 {
			masteryHandlers.BackfillContributions(database.DB)
		}
	}

	// Backfill mastery_experience from existing exercise_logs.
	var expCount int64
	database.DB.Raw("SELECT COUNT(*) FROM mastery_experience").Scan(&expCount)
	if expCount == 0 {
		var logCount int64
		database.DB.Raw("SELECT COUNT(*) FROM exercise_logs WHERE deleted_at IS NULL").Scan(&logCount)
		if logCount > 0 {
			database.DB.Exec(`INSERT INTO mastery_experience (owner, exercise_id, total_reps)
				SELECT owner, exercise_id, SUM(COALESCE(reps, 1))
				FROM exercise_logs WHERE deleted_at IS NULL
				GROUP BY owner, exercise_id`)
		}
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
		&workoutScheduleModels.WorkoutScheduleEntity{},
		&workoutScheduleModels.SchedulePeriodEntity{},
		&workoutScheduleModels.ScheduleCommitmentEntity{},
		&masteryModels.MasteryContributionEntity{},
		&masteryModels.MasteryExperienceEntity{},
	)
}

func setupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	api.Use(auth.RequestTrace())
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
	workoutScheduleHandlers.RegisterRoutes(humaAPI)
	masteryHandlers.RegisterRoutes(humaAPI)
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
	workoutlog.StartCommitmentTicker(database.DB, 15*time.Minute)

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
	setupDeployStatus(r)
	setupSPA(r)
	return r
}

func setupDeployStatus(r *gin.Engine) {
	dokployURL := os.Getenv("DOKPLOY_URL")
	dokployKey := os.Getenv("DOKPLOY_API_KEY")
	dokployAppID := os.Getenv("DOKPLOY_APP_ID")
	if dokployURL == "" || dokployKey == "" || dokployAppID == "" {
		return
	}

	type cachedStatus struct {
		body      []byte
		fetchedAt time.Time
	}
	var (
		mu    sync.Mutex
		cache *cachedStatus
	)
	client := &http.Client{Timeout: 10 * time.Second}

	r.GET("/api/deploy-status", func(c *gin.Context) {
		mu.Lock()
		if cache != nil && time.Since(cache.fetchedAt) < 30*time.Second {
			body := cache.body
			mu.Unlock()
			c.Data(http.StatusOK, "application/json", body)
			return
		}
		mu.Unlock()

		url := fmt.Sprintf("%s/api/deployment.all?applicationId=%s", dokployURL, dokployAppID)
		req, err := http.NewRequestWithContext(c.Request.Context(), http.MethodGet, url, nil)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "failed to build request"})
			return
		}
		req.Header.Set("x-api-key", dokployKey)

		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "dokploy unreachable"})
			return
		}
		defer resp.Body.Close()

		raw, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
		if err != nil || resp.StatusCode != http.StatusOK {
			c.JSON(http.StatusBadGateway, gin.H{"error": "dokploy error", "upstream": resp.StatusCode})
			return
		}

		var deployments []struct {
			Status    string `json:"status"`
			Title     string `json:"title"`
			CreatedAt string `json:"createdAt"`
		}
		if err := json.Unmarshal(raw, &deployments); err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "bad response from dokploy"})
			return
		}

		result := gin.H{"status": "idle"}
		if len(deployments) > 0 {
			d := deployments[0]
			result = gin.H{"status": d.Status, "title": d.Title, "createdAt": d.CreatedAt}
		}

		body, _ := json.Marshal(result)
		mu.Lock()
		cache = &cachedStatus{body: body, fetchedAt: time.Now()}
		mu.Unlock()

		c.Data(http.StatusOK, "application/json", body)
	})
}

func main() {
	r := buildApp()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
