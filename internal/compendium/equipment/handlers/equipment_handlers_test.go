package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"gesitr/internal/compendium/equipment/models"
	"gesitr/internal/database"
	"gesitr/internal/shared"
)

func TestListEquipment(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	t.Run("empty list", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/equipment", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		var page paginatedJSON
		json.Unmarshal(w.Body.Bytes(), &page)
		var result []models.Equipment
		json.Unmarshal(page.Items, &result)
		if len(result) != 0 {
			t.Errorf("expected 0, got %d", len(result))
		}
	})

	// Seed data for filter tests
	doJSON(r, "POST", "/api/equipment", map[string]any{
		"name": "barbell", "displayName": "Barbell", "description": "A bar",
		"category": "free_weights"})
	doJSON(r, "POST", "/api/equipment", map[string]any{
		"name": "bench", "displayName": "Flat Bench", "description": "A bench",
		"category": "benches"})

	t.Run("list all own equipment", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/equipment", nil)
		var page paginatedJSON
		json.Unmarshal(w.Body.Bytes(), &page)
		var result []models.Equipment
		json.Unmarshal(page.Items, &result)
		if len(result) != 2 {
			t.Errorf("expected 2, got %d", len(result))
		}
	})

	t.Run("filter by q name", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/equipment?q=bar", nil)
		var page paginatedJSON
		json.Unmarshal(w.Body.Bytes(), &page)
		var result []models.Equipment
		json.Unmarshal(page.Items, &result)
		if len(result) != 1 || result[0].Name != "barbell" {
			t.Errorf("q filter: got %d results", len(result))
		}
	})

	t.Run("filter by q displayName", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/equipment?q=Flat", nil)
		var page paginatedJSON
		json.Unmarshal(w.Body.Bytes(), &page)
		var result []models.Equipment
		json.Unmarshal(page.Items, &result)
		if len(result) != 1 || result[0].Name != "bench" {
			t.Errorf("q displayName filter: got %d results", len(result))
		}
	})

	t.Run("filter by category", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/equipment?category=benches", nil)
		var page paginatedJSON
		json.Unmarshal(w.Body.Bytes(), &page)
		var result []models.Equipment
		json.Unmarshal(page.Items, &result)
		if len(result) != 1 || result[0].Name != "bench" {
			t.Errorf("category filter: got %d results", len(result))
		}
	})

	t.Run("db error", func(t *testing.T) {
		closeDB(t)
		w := doJSON(r, "GET", "/api/equipment", nil)
		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", w.Code)
		}
	})
}

func TestListEquipmentVisibility(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	// testuser creates private equipment
	doJSON(r, "POST", "/api/equipment", map[string]any{
		"name": "private-bar", "displayName": "Private Bar", "description": "",
		"category": "free_weights"})
	// testuser creates public equipment
	doJSON(r, "POST", "/api/equipment", map[string]any{
		"name": "public-bar", "displayName": "Public Bar", "description": "",
		"category": "free_weights", "public": true,
	})

	// otheruser creates private equipment
	doJSONAs(r, "POST", "/api/equipment", "otheruser", map[string]any{
		"name": "other-private", "displayName": "Other Private", "description": "",
		"category": "free_weights"})
	// otheruser creates public equipment
	doJSONAs(r, "POST", "/api/equipment", "otheruser", map[string]any{
		"name": "other-public", "displayName": "Other Public", "description": "",
		"category": "free_weights", "public": true,
	})

	t.Run("default shows own + public", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/equipment", nil)
		var page paginatedJSON
		json.Unmarshal(w.Body.Bytes(), &page)
		var result []models.Equipment
		json.Unmarshal(page.Items, &result)
		// testuser sees: private-bar (own), public-bar (own), other-public (public)
		if len(result) != 3 {
			t.Errorf("expected 3 (own + public), got %d", len(result))
		}
	})

	t.Run("owner=me shows only own", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/equipment?owner=me", nil)
		var page paginatedJSON
		json.Unmarshal(w.Body.Bytes(), &page)
		var result []models.Equipment
		json.Unmarshal(page.Items, &result)
		if len(result) != 2 {
			t.Errorf("expected 2 (own only), got %d", len(result))
		}
	})

	t.Run("owner=otheruser shows only their public", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/equipment?owner=otheruser", nil)
		var page paginatedJSON
		json.Unmarshal(w.Body.Bytes(), &page)
		var result []models.Equipment
		json.Unmarshal(page.Items, &result)
		if len(result) != 1 {
			t.Errorf("expected 1 (other public only), got %d", len(result))
		}
		if len(result) > 0 && result[0].Name != "other-public" {
			t.Errorf("expected other-public, got %q", result[0].Name)
		}
	})

	t.Run("public=true shows only public", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/equipment?public=true", nil)
		var page paginatedJSON
		json.Unmarshal(w.Body.Bytes(), &page)
		var result []models.Equipment
		json.Unmarshal(page.Items, &result)
		if len(result) != 2 {
			t.Errorf("expected 2 (all public), got %d", len(result))
		}
	})

	t.Run("otheruser default view", func(t *testing.T) {
		w := doJSONAs(r, "GET", "/api/equipment", "otheruser", nil)
		var page paginatedJSON
		json.Unmarshal(w.Body.Bytes(), &page)
		var result []models.Equipment
		json.Unmarshal(page.Items, &result)
		// otheruser sees: other-private (own), other-public (own), public-bar (public)
		if len(result) != 3 {
			t.Errorf("expected 3 (own + public), got %d", len(result))
		}
	})
}

func TestCreateEquipment(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	t.Run("success", func(t *testing.T) {
		w := doJSON(r, "POST", "/api/equipment", map[string]any{
			"name": "dumbbell", "displayName": "Dumbbell", "description": "A weight",
			"category": "free_weights"})
		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
		}
		var result models.Equipment
		json.Unmarshal(w.Body.Bytes(), &result)
		if result.ID == 0 || result.Name != "dumbbell" {
			t.Error("create response mismatch")
		}
	})

	t.Run("create returns valid equipment", func(t *testing.T) {
		w := doJSON(r, "POST", "/api/equipment", map[string]any{
			"name": "cable", "displayName": "Cable", "description": "",
			"category": "machines"})
		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
		}
		var result models.Equipment
		json.Unmarshal(w.Body.Bytes(), &result)
		if result.ID == 0 {
			t.Errorf("expected non-zero ID")
		}
	})

	t.Run("bad json", func(t *testing.T) {
		w := doRaw(r, "POST", "/api/equipment", "{invalid")
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("db error", func(t *testing.T) {
		closeDB(t)
		w := doJSON(r, "POST", "/api/equipment", map[string]any{
			"name": "x", "displayName": "X", "description": "",
			"category": "other"})
		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", w.Code)
		}
	})
}

func TestGetEquipment(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	// Create one
	w := doJSON(r, "POST", "/api/equipment", map[string]any{
		"name": "kettlebell", "displayName": "Kettlebell", "description": "",
		"category": "free_weights"})
	var created models.Equipment
	json.Unmarshal(w.Body.Bytes(), &created)

	t.Run("found", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/equipment/1", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		var result models.Equipment
		json.Unmarshal(w.Body.Bytes(), &result)
		if result.Name != "kettlebell" {
			t.Error("get response mismatch")
		}
	})

	t.Run("not found", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/equipment/999", nil)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})
}

func TestUpdateEquipment(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	doJSON(r, "POST", "/api/equipment", map[string]any{
		"name": "band", "displayName": "Band", "description": "",
		"category": "accessories"})

	t.Run("success", func(t *testing.T) {
		w := doJSON(r, "PUT", "/api/equipment/1", map[string]any{
			"name": "resistance-band", "displayName": "Resistance Band", "description": "elastic",
			"category": "accessories"})
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
		}
		var result models.Equipment
		json.Unmarshal(w.Body.Bytes(), &result)
		if result.Name != "resistance-band" || result.Version != 1 {
			t.Errorf("update response: name=%q version=%d", result.Name, result.Version)
		}
	})

	t.Run("no version bump when unchanged", func(t *testing.T) {
		w := doJSON(r, "PUT", "/api/equipment/1", map[string]any{
			"name": "resistance-band", "displayName": "Resistance Band", "description": "elastic",
			"category": "accessories"})
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		var result models.Equipment
		json.Unmarshal(w.Body.Bytes(), &result)
		if result.Version != 1 {
			t.Errorf("Version = %d, want 1 (should not have bumped)", result.Version)
		}
	})

	t.Run("no extra history on unchanged update", func(t *testing.T) {
		var count int64
		database.DB.Model(&models.EquipmentHistoryEntity{}).Where("equipment_id = ?", 1).Count(&count)
		if count != 2 {
			t.Errorf("expected 2 history records (no extra from no-op), got %d", count)
		}
	})

	t.Run("successive updates accumulate history", func(t *testing.T) {
		w := doJSON(r, "PUT", "/api/equipment/1", map[string]any{
			"name": "super-band", "displayName": "Super Band", "description": "v2",
			"category": "accessories"})
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		var result models.Equipment
		json.Unmarshal(w.Body.Bytes(), &result)
		if result.Version != 2 {
			t.Errorf("Version = %d, want 2", result.Version)
		}
		var count int64
		database.DB.Model(&models.EquipmentHistoryEntity{}).Where("equipment_id = ?", 1).Count(&count)
		if count != 3 {
			t.Errorf("expected 3 history records (v0, v1, v2), got %d", count)
		}
	})

	t.Run("forbidden for non-owner", func(t *testing.T) {
		w := doJSONAs(r, "PUT", "/api/equipment/1", "otheruser", map[string]any{
			"name": "hijack", "displayName": "Hijack", "description": "",
			"category": "accessories"})
		if w.Code != http.StatusForbidden {
			t.Errorf("expected 403, got %d", w.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		w := doJSON(r, "PUT", "/api/equipment/999", map[string]any{
			"name": "x", "displayName": "X", "description": "",
			"category": "other"})
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})

	t.Run("bad json", func(t *testing.T) {
		w := doRaw(r, "PUT", "/api/equipment/1", "{bad")
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("db error first lookup", func(t *testing.T) {
		closeDB(t)
		w := doJSON(r, "PUT", "/api/equipment/1", map[string]any{
			"name": "x", "displayName": "X", "description": "",
			"category": "other"})
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404 (db closed), got %d", w.Code)
		}
	})
}

func TestListEquipmentVersions(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	// Create equipment (v0) and update it (v1)
	doJSON(r, "POST", "/api/equipment", map[string]any{
		"name": "plate", "displayName": "Plate", "description": "A weight plate",
		"category": "free_weights"})
	doJSON(r, "PUT", "/api/equipment/1", map[string]any{
		"name": "bumper-plate", "displayName": "Bumper Plate", "description": "Rubber coated",
		"category": "free_weights"})

	t.Run("returns all versions ordered", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/equipment/1/versions", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
		}
		var entries []shared.VersionEntry
		json.Unmarshal(w.Body.Bytes(), &entries)
		if len(entries) != 2 {
			t.Fatalf("expected 2 versions, got %d", len(entries))
		}
		if entries[0].Version != 0 || entries[1].Version != 1 {
			t.Errorf("versions = %d, %d", entries[0].Version, entries[1].Version)
		}
	})

	t.Run("snapshot contains correct data", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/equipment/1/versions", nil)
		var entries []shared.VersionEntry
		json.Unmarshal(w.Body.Bytes(), &entries)

		var v0 models.Equipment
		json.Unmarshal(entries[0].Snapshot, &v0)
		if v0.Name != "plate" || v0.DisplayName != "Plate" {
			t.Errorf("v0 snapshot: name=%q displayName=%q", v0.Name, v0.DisplayName)
		}

		var v1 models.Equipment
		json.Unmarshal(entries[1].Snapshot, &v1)
		if v1.Name != "bumper-plate" || v1.Description != "Rubber coated" {
			t.Errorf("v1 snapshot: name=%q desc=%q", v1.Name, v1.Description)
		}
	})

	t.Run("not found", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/equipment/999/versions", nil)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})

	t.Run("db error", func(t *testing.T) {
		closeDB(t)
		w := doJSON(r, "GET", "/api/equipment/1/versions", nil)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404 (db closed), got %d", w.Code)
		}
	})
}

func TestGetEquipmentVersion(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	// Create equipment (v0) and update it (v1)
	doJSON(r, "POST", "/api/equipment", map[string]any{
		"name": "plate", "displayName": "Plate", "description": "A weight plate",
		"category": "free_weights"})
	doJSON(r, "PUT", "/api/equipment/1", map[string]any{
		"name": "bumper-plate", "displayName": "Bumper Plate", "description": "Rubber coated",
		"category": "free_weights"})

	t.Run("returns specific version", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/equipment/1/versions/0", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
		}
		var entry shared.VersionEntry
		json.Unmarshal(w.Body.Bytes(), &entry)
		if entry.Version != 0 {
			t.Errorf("version = %d, want 0", entry.Version)
		}
		var snapshot models.Equipment
		json.Unmarshal(entry.Snapshot, &snapshot)
		if snapshot.Name != "plate" {
			t.Errorf("snapshot name = %q, want plate", snapshot.Name)
		}
	})

	t.Run("returns v1", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/equipment/1/versions/1", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
		}
		var entry shared.VersionEntry
		json.Unmarshal(w.Body.Bytes(), &entry)
		if entry.Version != 1 {
			t.Errorf("version = %d, want 1", entry.Version)
		}
	})

	t.Run("equipment not found", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/equipment/999/versions/0", nil)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})

	t.Run("version not found", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/equipment/1/versions/999", nil)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})

	t.Run("invalid version", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/equipment/1/versions/abc", nil)
		if w.Code != http.StatusUnprocessableEntity {
			t.Errorf("expected 422, got %d", w.Code)
		}
	})
}

func TestDeleteEquipment(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	doJSON(r, "POST", "/api/equipment", map[string]any{
		"name": "rope", "displayName": "Rope", "description": "",
		"category": "accessories"})

	t.Run("success", func(t *testing.T) {
		w := doJSON(r, "DELETE", "/api/equipment/1", nil)
		if w.Code != http.StatusNoContent {
			t.Errorf("expected 204, got %d", w.Code)
		}
	})

	t.Run("forbidden for non-owner", func(t *testing.T) {
		setupTestDB(t)
		r := newRouter()
		doJSON(r, "POST", "/api/equipment", map[string]any{
			"name": "rope2", "displayName": "Rope2", "description": "",
			"category": "accessories"})
		w := doJSONAs(r, "DELETE", "/api/equipment/1", "otheruser", nil)
		if w.Code != http.StatusForbidden {
			t.Errorf("expected 403, got %d", w.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		w := doJSON(r, "DELETE", "/api/equipment/999", nil)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})

	t.Run("db error", func(t *testing.T) {
		closeDB(t)
		w := doJSON(r, "DELETE", "/api/equipment/1", nil)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404 (db closed), got %d", w.Code)
		}
	})
}

func TestGetEquipmentPermissions(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	// Create a public equipment owned by testuser
	doJSON(r, "POST", "/api/equipment", map[string]any{
		"name": "barbell", "displayName": "Barbell", "description": "Standard barbell",
		"category": "free_weights", "public": true,
	})

	// Create a private equipment owned by testuser
	doJSON(r, "POST", "/api/equipment", map[string]any{
		"name": "custom-band", "displayName": "Custom Band", "description": "Private",
		"category": "accessories", "public": false,
	})

	t.Run("owner gets full permissions", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/equipment/1/permissions", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
		}
		var resp shared.PermissionsResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		if len(resp.Permissions) != 3 {
			t.Fatalf("expected 3 permissions, got %d: %v", len(resp.Permissions), resp.Permissions)
		}
	})

	t.Run("non-owner on public gets READ only", func(t *testing.T) {
		w := doJSONAs(r, "GET", "/api/equipment/1/permissions", "system", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
		}
		var resp shared.PermissionsResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		if len(resp.Permissions) != 1 || resp.Permissions[0] != "READ" {
			t.Fatalf("expected [READ], got %v", resp.Permissions)
		}
	})

	t.Run("non-owner on private gets empty permissions", func(t *testing.T) {
		w := doJSONAs(r, "GET", "/api/equipment/2/permissions", "system", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
		var resp shared.PermissionsResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		if len(resp.Permissions) != 0 {
			t.Errorf("expected empty permissions, got %v", resp.Permissions)
		}
	})

	t.Run("owner gets READ, MODIFY, DELETE", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/equipment/1/permissions", nil)
		var resp shared.PermissionsResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		expected := map[shared.Permission]bool{
			shared.PermissionRead:   false,
			shared.PermissionModify: false,
			shared.PermissionDelete: false,
		}
		for _, p := range resp.Permissions {
			expected[p] = true
		}
		for perm, found := range expected {
			if !found {
				t.Errorf("missing permission: %s", perm)
			}
		}
	})

	t.Run("owner on private still gets full permissions", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/equipment/2/permissions", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
		}
		var resp shared.PermissionsResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		if len(resp.Permissions) != 3 {
			t.Fatalf("expected 3 permissions for owner on private, got %d: %v", len(resp.Permissions), resp.Permissions)
		}
	})

	t.Run("not found", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/equipment/999/permissions", nil)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})

	t.Run("db error", func(t *testing.T) {
		closeDB(t)
		w := doJSON(r, "GET", "/api/equipment/1/permissions", nil)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404 (db closed), got %d", w.Code)
		}
	})
}
