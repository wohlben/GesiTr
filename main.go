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
	"regexp"
	"sync"
	"time"

	"gesitr/internal/auth"
	equipmentHandlers "gesitr/internal/compendium/equipment/handlers"
	equipmentModels "gesitr/internal/compendium/equipment/models"
	exerciseHandlers "gesitr/internal/compendium/exercise/handlers"
	exerciseModels "gesitr/internal/compendium/exercise/models"
	localityHandlers "gesitr/internal/compendium/locality/handlers"
	localityModels "gesitr/internal/compendium/locality/models"
	"gesitr/internal/compendium/ownershipgroup"
	ownershipGroupHandlers "gesitr/internal/compendium/ownershipgroup/handlers"
	ownershipGroupModels "gesitr/internal/compendium/ownershipgroup/models"
	workoutHandlers "gesitr/internal/compendium/workout/handlers"
	workoutModels "gesitr/internal/compendium/workout/models"
	workoutGroupHandlers "gesitr/internal/compendium/workoutgroup/handlers"
	workoutGroupModels "gesitr/internal/compendium/workoutgroup/models"
	"gesitr/internal/database"
	"gesitr/internal/docs"
	"gesitr/internal/humaconfig"
	"gesitr/internal/profile"
	profileHandlers "gesitr/internal/profile/handlers"
	profileModels "gesitr/internal/profile/models"
	exerciseLogHandlers "gesitr/internal/user/exerciselog/handlers"
	exerciseLogModels "gesitr/internal/user/exerciselog/models"
	exerciseSchemeHandlers "gesitr/internal/user/exercisescheme/handlers"
	exerciseSchemeModels "gesitr/internal/user/exercisescheme/models"
	masteryHandlers "gesitr/internal/user/mastery/handlers"
	masteryModels "gesitr/internal/user/mastery/models"
	namePreferenceHandlers "gesitr/internal/user/namepreference/handlers"
	namePreferenceModels "gesitr/internal/user/namepreference/models"
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

	// Backfill equipment_mastery_contributions from existing equipment relationships + fulfillments.
	var eqContribCount int64
	database.DB.Raw("SELECT COUNT(*) FROM equipment_mastery_contributions").Scan(&eqContribCount)
	if eqContribCount == 0 {
		var eqRelCount int64
		database.DB.Raw("SELECT COUNT(*) FROM equipment_relationships WHERE deleted_at IS NULL").Scan(&eqRelCount)
		var fulCount int64
		database.DB.Raw("SELECT COUNT(*) FROM fulfillments").Scan(&fulCount)
		if eqRelCount > 0 || fulCount > 0 {
			masteryHandlers.BackfillEquipmentContributions(database.DB)
		}
	}

	// Backfill equipment_mastery_experience from existing exercise_logs + exercise_equipments.
	var eqExpCount int64
	database.DB.Raw("SELECT COUNT(*) FROM equipment_mastery_experience").Scan(&eqExpCount)
	if eqExpCount == 0 {
		var logCount int64
		database.DB.Raw("SELECT COUNT(*) FROM exercise_logs WHERE deleted_at IS NULL").Scan(&logCount)
		if logCount > 0 {
			database.DB.Exec(`INSERT INTO equipment_mastery_experience (owner, equipment_id, total_reps)
				SELECT el.owner, ee.equipment_id, SUM(COALESCE(el.reps, 1))
				FROM exercise_logs el
				JOIN exercise_equipments ee ON ee.exercise_id = el.exercise_id
				WHERE el.deleted_at IS NULL
				GROUP BY el.owner, ee.equipment_id`)
		}
	}

	// Drop unique index on workout_groups.workout_id to allow multiple groups per workout.
	database.DB.Exec("DROP INDEX IF EXISTS idx_workout_groups_workout_id")
	database.DB.Exec("DROP INDEX IF EXISTS uni_workout_groups_workout_id")

	// Migrate exercise names: move exercises.name + exercise_alternative_names → exercise_names table.
	var nameColExists int64
	database.DB.Raw("SELECT COUNT(*) FROM pragma_table_info('exercises') WHERE name = 'name'").Scan(&nameColExists)
	var altTableExists int64
	database.DB.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='exercise_alternative_names'").Scan(&altTableExists)
	if nameColExists > 0 {
		// Insert the primary name at position 0.
		database.DB.Exec(`INSERT INTO exercise_names (exercise_id, position, name)
			SELECT id, 0, name FROM exercises WHERE name IS NOT NULL AND name != ''`)
		if altTableExists > 0 {
			// Insert alternative names at incrementing positions, skipping duplicates of the primary name.
			database.DB.Exec(`INSERT INTO exercise_names (exercise_id, position, name)
				SELECT ean.exercise_id,
					(SELECT COALESCE(MAX(position), -1) + 1 FROM exercise_names en WHERE en.exercise_id = ean.exercise_id),
					ean.name
				FROM exercise_alternative_names ean
				WHERE NOT EXISTS (
					SELECT 1 FROM exercise_names en
					WHERE en.exercise_id = ean.exercise_id AND en.name = ean.name
				)`)
			database.DB.Exec("DROP TABLE exercise_alternative_names")
		}
		database.DB.Exec("ALTER TABLE exercises DROP COLUMN name")
	} else if altTableExists > 0 {
		// Column already dropped but old table lingers — just migrate remaining data.
		database.DB.Exec(`INSERT INTO exercise_names (exercise_id, position, name)
			SELECT ean.exercise_id,
				(SELECT COALESCE(MAX(position), -1) + 1 FROM exercise_names en WHERE en.exercise_id = ean.exercise_id),
				ean.name
			FROM exercise_alternative_names ean
			WHERE NOT EXISTS (
				SELECT 1 FROM exercise_names en
				WHERE en.exercise_id = ean.exercise_id AND en.name = ean.name
			)`)
		database.DB.Exec("DROP TABLE exercise_alternative_names")
	}

	// Remove user_profiles table and FK constraints referencing it.
	removeProfileForeignKeys()
}

// removeProfileForeignKeys drops FK constraints referencing the user_profiles table
// and then drops the table itself. This uses the standard SQLite table-recreation
// approach since SQLite doesn't support ALTER TABLE DROP CONSTRAINT.
func removeProfileForeignKeys() {
	var tableExists int64
	database.DB.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='user_profiles'").Scan(&tableExists)
	if tableExists == 0 {
		return
	}

	database.DB.Exec("PRAGMA foreign_keys = OFF")
	defer database.DB.Exec("PRAGMA foreign_keys = ON")

	// Find all tables with FK constraints referencing user_profiles.
	var tables []struct {
		Name string
		SQL  string
	}
	database.DB.Raw("SELECT name, sql FROM sqlite_master WHERE type='table' AND sql LIKE '%user_profiles%'").Scan(&tables)

	re := regexp.MustCompile(",\\s*CONSTRAINT\\s+`[^`]+`\\s+FOREIGN KEY\\s*\\([^)]+\\)\\s*REFERENCES\\s*`user_profiles`\\s*\\([^)]+\\)(\\s+ON DELETE\\s+\\w+)?")

	for _, t := range tables {
		if t.Name == "user_profiles" {
			continue
		}
		newDDL := re.ReplaceAllString(t.SQL, "")
		if newDDL == t.SQL {
			continue
		}

		// Collect indexes to recreate after table swap.
		var indexes []struct{ SQL *string }
		database.DB.Raw("SELECT sql FROM sqlite_master WHERE type='index' AND tbl_name=? AND sql IS NOT NULL", t.Name).Scan(&indexes)

		tempName := "_old_" + t.Name
		database.DB.Exec("ALTER TABLE `" + t.Name + "` RENAME TO `" + tempName + "`")
		database.DB.Exec(newDDL)
		database.DB.Exec("INSERT INTO `" + t.Name + "` SELECT * FROM `" + tempName + "`")
		database.DB.Exec("DROP TABLE `" + tempName + "`")

		for _, idx := range indexes {
			if idx.SQL != nil {
				database.DB.Exec(*idx.SQL)
			}
		}
	}

	database.DB.Exec("DROP TABLE IF EXISTS user_profiles")
}

func autoMigrate() {
	database.DB.AutoMigrate(
		&exerciseModels.ExerciseEntity{},
		&exerciseModels.ExerciseForce{},
		&exerciseModels.ExerciseMuscle{},
		&exerciseModels.ExerciseMeasurementParadigm{},
		&exerciseModels.ExerciseInstruction{},
		&exerciseModels.ExerciseImage{},
		&exerciseModels.ExerciseName{},
		&equipmentModels.EquipmentEntity{},
		&exerciseModels.ExerciseEquipment{},
		&equipmentModels.FulfillmentEntity{},
		&exerciseModels.ExerciseRelationshipEntity{},
		&workoutModels.ExerciseGroupEntity{},
		&workoutModels.ExerciseGroupMemberEntity{},
		&exerciseModels.ExerciseHistoryEntity{},
		&equipmentModels.EquipmentHistoryEntity{},
		&exerciseSchemeModels.ExerciseSchemeEntity{},
		&exerciseSchemeModels.ExerciseSchemeSectionItemEntity{},
		&equipmentModels.EquipmentRelationshipEntity{},
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
		&masteryModels.EquipmentMasteryContributionEntity{},
		&masteryModels.EquipmentMasteryExperienceEntity{},
		&workoutModels.WorkoutHistoryEntity{},
		&workoutModels.WorkoutRelationshipEntity{},
		&namePreferenceModels.ExerciseNamePreference{},
		&localityModels.LocalityEntity{},
		&localityModels.LocalityAvailabilityEntity{},
		&ownershipGroupModels.OwnershipGroupEntity{},
		&ownershipGroupModels.OwnershipGroupMembershipEntity{},
		&profileModels.ProfileEntity{},
	)
}

func setupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	api.Use(auth.RequestTrace())
	api.Use(auth.UserID())

	// Huma API — shares the /api group so Gin auth middleware applies.
	humaAPI := humaconfig.NewAPI(r, api)
	exerciseHandlers.RegisterRoutes(humaAPI)
	equipmentHandlers.RegisterRoutes(humaAPI)
	workoutHandlers.RegisterRoutes(humaAPI)
	workoutloghandlers.RegisterRoutes(humaAPI)
	exerciseLogHandlers.RegisterRoutes(humaAPI)
	workoutGroupHandlers.RegisterRoutes(humaAPI)
	workoutScheduleHandlers.RegisterRoutes(humaAPI)
	masteryHandlers.RegisterRoutes(humaAPI)
	namePreferenceHandlers.RegisterRoutes(humaAPI)
	exerciseSchemeHandlers.RegisterRoutes(humaAPI)
	localityHandlers.RegisterRoutes(humaAPI)
	ownershipGroupHandlers.RegisterRoutes(humaAPI)
	profileHandlers.RegisterRoutes(humaAPI)
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
	spaHandler := func(c *gin.Context) {
		f, err := http.FS(distFS).Open(c.Request.URL.Path)
		if err == nil {
			f.Close()
			c.FileFromFS(c.Request.URL.Path, http.FS(distFS))
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", indexHTML)
	}
	r.NoRoute(profile.RequireProfile(database.DB, spaHandler))
}

func buildApp() *gin.Engine {
	database.Init()
	autoMigrate()
	ownershipgroup.MigrateExistingOwners(database.DB)
	runMigrations()
	workoutlog.StartCommitmentTicker(database.DB, 15*time.Minute)

	var r *gin.Engine
	if level := auth.RequestLogLevel(); level == "info" || level == "trace" {
		r = gin.New()
		r.Use(auth.APIOnlyLogger(), gin.Recovery())
	} else {
		r = gin.Default()
	}
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

			// Re-create profile for the fallback user so e2e tests can load the SPA.
			if fallback := os.Getenv("AUTH_FALLBACK_USER"); fallback != "" {
				database.DB.Create(&profileModels.ProfileEntity{
					UserID:   fallback,
					Username: fallback,
				})
			}

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

		// Dokploy may return a raw array or a wrapped object — try both.
		var records []json.RawMessage
		if err := json.Unmarshal(raw, &records); err != nil {
			// Try wrapped: { "data": [...] } or similar
			var wrapped map[string]json.RawMessage
			if err2 := json.Unmarshal(raw, &wrapped); err2 != nil {
				c.JSON(http.StatusBadGateway, gin.H{"error": "bad response from dokploy", "raw": string(raw[:min(len(raw), 500)])})
				return
			}
			// Try common wrapper keys
			for _, key := range []string{"data", "deployments", "items"} {
				if v, ok := wrapped[key]; ok {
					if json.Unmarshal(v, &records) == nil {
						break
					}
				}
			}
		}

		type deployment struct {
			Status      string `json:"status"`
			Title       string `json:"title"`
			Description string `json:"description"`
			CreatedAt   string `json:"createdAt"`
		}

		result := gin.H{"status": "idle"}
		for _, rec := range records {
			var d deployment
			if json.Unmarshal(rec, &d) == nil && d.Status != "" {
				title := d.Title
				if title == "" {
					title = d.Description
				}
				result = gin.H{"status": d.Status, "title": title, "createdAt": d.CreatedAt}
				break // First (most recent) deployment with a status
			}
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
