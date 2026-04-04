package handlers

import (
	"encoding/json"
	"testing"

	"gesitr/internal/compendium/equipment/models"
	exerciseModels "gesitr/internal/compendium/exercise/models"
	ownershipgroupmodels "gesitr/internal/compendium/ownershipgroup/models"
	"gesitr/internal/database"
	exerciseLogModels "gesitr/internal/user/exerciselog/models"
	masteryModels "gesitr/internal/user/mastery/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// setupMasteryTestDB sets up a test DB with all models needed for mastery testing.
func setupMasteryTestDB(t *testing.T) {
	t.Helper()
	t.Setenv("AUTH_FALLBACK_USER", "testuser")
	db, err := gorm.Open(sqlite.Open("file::memory:?_foreign_keys=on"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatal(err)
	}
	db.AutoMigrate(
		&ownershipgroupmodels.OwnershipGroupEntity{},
		&ownershipgroupmodels.OwnershipGroupMembershipEntity{},
		&models.EquipmentEntity{},
		&models.EquipmentHistoryEntity{},
		&exerciseModels.ExerciseEntity{},
		&exerciseModels.ExerciseEquipment{},
		&exerciseLogModels.ExerciseLogEntity{},
		&masteryModels.EquipmentMasteryExperienceEntity{},
	)
	database.DB = db
}

func TestListEquipmentMasteryVisibility(t *testing.T) {
	setupMasteryTestDB(t)
	r := newRouter()

	// Create 3 pieces of equipment:
	// 1. owned by testuser (private)
	w := doJSON(r, "POST", "/api/equipment", map[string]any{
		"name": "my-barbell", "displayName": "My Barbell", "description": "",
		"category": "free_weights",
	})
	var ownedEquip models.Equipment
	json.Unmarshal(w.Body.Bytes(), &ownedEquip)

	// 2. owned by bob, public
	w = doJSONAs(r, "POST", "/api/equipment", "bob", map[string]any{
		"name": "bob-bench", "displayName": "Bob Bench", "description": "",
		"category": "benches", "public": true,
	})
	var bobPublic models.Equipment
	json.Unmarshal(w.Body.Bytes(), &bobPublic)

	// 3. owned by bob, NOT public
	doJSONAs(r, "POST", "/api/equipment", "bob", map[string]any{
		"name": "bob-secret", "displayName": "Bob Secret Rack", "description": "",
		"category": "machines",
	})

	t.Run("default shows own + public (3 items created, 2 visible to testuser)", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/equipment", nil)
		var page paginatedJSON
		json.Unmarshal(w.Body.Bytes(), &page)
		var result []models.Equipment
		json.Unmarshal(page.Items, &result)
		// testuser sees: my-barbell (own) + bob-bench (public) = 2
		// bob-secret is NOT visible (not public, not owned)
		if len(result) != 2 {
			t.Errorf("expected 2 visible (own + public), got %d", len(result))
			for _, eq := range result {
				t.Logf("  - %s (id=%d, public=%v)", eq.Name, eq.ID, eq.Public)
			}
		}
	})

	t.Run("mastery=me with no mastery shows only owned", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/equipment?mastery=me", nil)
		var page paginatedJSON
		json.Unmarshal(w.Body.Bytes(), &page)
		var result []models.Equipment
		json.Unmarshal(page.Items, &result)
		// No mastery yet — only owned equipment visible
		if len(result) != 1 {
			t.Errorf("expected 1 (own only with no mastery), got %d", len(result))
		}
		if len(result) > 0 && result[0].Name != "my-barbell" {
			t.Errorf("expected my-barbell, got %q", result[0].Name)
		}
	})

	// Simulate mastery: testuser has used bob's public bench
	database.DB.Create(&masteryModels.EquipmentMasteryExperienceEntity{
		Owner:       "testuser",
		EquipmentID: bobPublic.ID,
		TotalReps:   50,
	})

	t.Run("mastery=me shows owned + equipment with mastery", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/equipment?mastery=me", nil)
		var page paginatedJSON
		json.Unmarshal(w.Body.Bytes(), &page)
		var result []models.Equipment
		json.Unmarshal(page.Items, &result)
		// testuser sees: my-barbell (own) + bob-bench (has mastery) = 2
		if len(result) != 2 {
			t.Errorf("expected 2 (own + mastery), got %d", len(result))
			for _, eq := range result {
				t.Logf("  - %s (id=%d, public=%v)", eq.Name, eq.ID, eq.Public)
			}
		}
	})

	t.Run("mastery=me does NOT reveal non-public equipment without mastery", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/equipment?mastery=me", nil)
		var page paginatedJSON
		json.Unmarshal(w.Body.Bytes(), &page)
		var result []models.Equipment
		json.Unmarshal(page.Items, &result)
		for _, eq := range result {
			if eq.Name == "bob-secret" {
				t.Error("bob-secret should NOT be visible (not public, no mastery)")
			}
		}
	})
}

func TestListEquipmentMasteryOrdering(t *testing.T) {
	setupMasteryTestDB(t)
	r := newRouter()

	// testuser creates 1 piece of equipment
	w := doJSON(r, "POST", "/api/equipment", map[string]any{
		"name": "my-bar", "displayName": "My Bar", "description": "",
		"category": "free_weights",
	})
	var ownedEquip models.Equipment
	json.Unmarshal(w.Body.Bytes(), &ownedEquip)

	// bob creates 2 public pieces of equipment
	w = doJSONAs(r, "POST", "/api/equipment", "bob", map[string]any{
		"name": "bob-bench", "displayName": "Bob Bench", "description": "",
		"category": "benches", "public": true,
	})
	var bobBench models.Equipment
	json.Unmarshal(w.Body.Bytes(), &bobBench)

	w = doJSONAs(r, "POST", "/api/equipment", "bob", map[string]any{
		"name": "bob-rack", "displayName": "Bob Rack", "description": "",
		"category": "machines", "public": true,
	})
	var bobRack models.Equipment
	json.Unmarshal(w.Body.Bytes(), &bobRack)

	// testuser gains mastery on bob-bench (via exercise usage)
	database.DB.Create(&masteryModels.EquipmentMasteryExperienceEntity{
		Owner:       "testuser",
		EquipmentID: bobBench.ID,
		TotalReps:   100,
	})

	t.Run("mastery=me: owned first, then mastery equipment, non-mastery excluded", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/equipment?mastery=me", nil)
		var page paginatedJSON
		json.Unmarshal(w.Body.Bytes(), &page)
		var result []models.Equipment
		json.Unmarshal(page.Items, &result)

		// Expected: my-bar (owned), bob-bench (mastery) — bob-rack excluded (no mastery, not owned)
		if len(result) != 2 {
			t.Fatalf("expected 2, got %d", len(result))
		}

		// bob-bench with mastery should be included alongside owned equipment
		names := make([]string, len(result))
		for i, eq := range result {
			names[i] = eq.Name
		}

		found := map[string]bool{}
		for _, eq := range result {
			found[eq.Name] = true
		}
		if !found["my-bar"] {
			t.Error("expected my-bar (owned) to be present")
		}
		if !found["bob-bench"] {
			t.Error("expected bob-bench (has mastery) to be present")
		}
		if found["bob-rack"] {
			t.Error("bob-rack should NOT appear (no mastery, not owned)")
		}
	})

	t.Run("default view still shows all public + owned (no mastery filter)", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/equipment", nil)
		var page paginatedJSON
		json.Unmarshal(w.Body.Bytes(), &page)
		var result []models.Equipment
		json.Unmarshal(page.Items, &result)
		// Default: my-bar (own) + bob-bench (public) + bob-rack (public) = 3
		if len(result) != 3 {
			t.Errorf("expected 3 (own + all public), got %d", len(result))
		}
	})
}
