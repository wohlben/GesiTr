package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"os"

	"gesitr/internal/compendium/equipment/models"
	"gesitr/internal/database"
	ownershipgroupmodels "gesitr/internal/compendium/ownershipgroup/models"
	"gesitr/internal/shared"

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
		&models.EquipmentEntity{},
		&models.EquipmentHistoryEntity{},
		&ownershipgroupmodels.OwnershipGroupEntity{},
		&ownershipgroupmodels.OwnershipGroupMembershipEntity{},
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

// Owner gets full permissions (READ, MODIFY, DELETE) on their own equipment.
func ExampleGetEquipmentPermissions_owner() {
	setupExampleDB()
	r := newRouter()

	// Create private equipment owned by testuser.
	doRaw(r, "POST", "/api/equipment", `{
		"name": "barbell",
		"displayName": "Barbell",
		"description": "Standard barbell",
		"category": "free_weights"
	}`)

	// Query permissions as the owner.
	w := doJSON(r, "GET", "/api/equipment/1/permissions", nil)

	var resp shared.PermissionsResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	fmt.Println(resp.Permissions)
	// Output: [READ MODIFY DELETE]
}

// Non-owner can read public equipment but cannot modify or delete it.
func ExampleGetEquipmentPermissions_nonOwnerPublic() {
	setupExampleDB()
	r := newRouter()

	// Create public equipment owned by testuser.
	doRaw(r, "POST", "/api/equipment", `{
		"name": "dumbbell",
		"displayName": "Dumbbell",
		"description": "Adjustable dumbbell",
		"category": "free_weights",
		"public": true
	}`)

	// Query permissions as a different user.
	w := doRawAs(r, "GET", "/api/equipment/1/permissions", "", "other")

	var resp shared.PermissionsResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	fmt.Println(resp.Permissions)
	// Output: [READ]
}

// Non-owner has no permissions on private equipment.
func ExampleGetEquipmentPermissions_nonOwnerPrivate() {
	setupExampleDB()
	r := newRouter()

	// Create private equipment owned by testuser.
	doRaw(r, "POST", "/api/equipment", `{
		"name": "custom-bar",
		"displayName": "Custom Bar",
		"description": "Private",
		"category": "free_weights"
	}`)

	// Query permissions as a different user.
	w := doRawAs(r, "GET", "/api/equipment/1/permissions", "", "other")

	var resp shared.PermissionsResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	fmt.Println(resp.Permissions)
	// Output: []
}
