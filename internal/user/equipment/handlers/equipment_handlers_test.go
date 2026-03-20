package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"gesitr/internal/user/equipment/models"
)

func TestListUserEquipment(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	t.Run("empty list", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/equipment", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		var result []models.UserEquipment
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 0 {
			t.Errorf("expected 0, got %d", len(result))
		}
	})

	// Seed data
	doJSON(r, "POST", "/api/user/equipment", map[string]any{
		"owner": "alice", "compendiumEquipmentId": "barbell", "compendiumVersion": 1,
	})
	doJSON(r, "POST", "/api/user/equipment", map[string]any{
		"owner": "bob", "compendiumEquipmentId": "dumbbell", "compendiumVersion": 2,
	})

	t.Run("list all", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/equipment", nil)
		var result []models.UserEquipment
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 2 {
			t.Errorf("expected 2, got %d", len(result))
		}
	})

	t.Run("filter by owner", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/equipment?owner=alice", nil)
		var result []models.UserEquipment
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 1 || result[0].Owner != "alice" {
			t.Errorf("owner filter: got %d results", len(result))
		}
	})

	t.Run("filter by compendiumEquipmentId", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/equipment?compendiumEquipmentId=dumbbell", nil)
		var result []models.UserEquipment
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 1 || result[0].CompendiumEquipmentID != "dumbbell" {
			t.Errorf("templateId filter: got %d results", len(result))
		}
	})

	t.Run("db error", func(t *testing.T) {
		closeDB(t)
		w := doJSON(r, "GET", "/api/user/equipment", nil)
		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", w.Code)
		}
	})
}

func TestCreateUserEquipment(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	t.Run("success", func(t *testing.T) {
		w := doJSON(r, "POST", "/api/user/equipment", map[string]any{
			"owner": "alice", "compendiumEquipmentId": "kettlebell", "compendiumVersion": 1,
		})
		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
		}
		var result models.UserEquipment
		json.Unmarshal(w.Body.Bytes(), &result)
		if result.ID == 0 || result.Owner != "alice" || result.CompendiumEquipmentID != "kettlebell" || result.CompendiumVersion != 1 {
			t.Error("create response mismatch")
		}
	})

	t.Run("duplicate import", func(t *testing.T) {
		w := doJSON(r, "POST", "/api/user/equipment", map[string]any{
			"owner": "alice", "compendiumEquipmentId": "kettlebell", "compendiumVersion": 1,
		})
		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected 500 for duplicate, got %d", w.Code)
		}
	})

	t.Run("bad json", func(t *testing.T) {
		w := doRaw(r, "POST", "/api/user/equipment", "{invalid")
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("db error", func(t *testing.T) {
		closeDB(t)
		w := doJSON(r, "POST", "/api/user/equipment", map[string]any{
			"owner": "x", "compendiumEquipmentId": "x", "compendiumVersion": 0,
		})
		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", w.Code)
		}
	})
}

func TestGetUserEquipment(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	doJSON(r, "POST", "/api/user/equipment", map[string]any{
		"owner": "alice", "compendiumEquipmentId": "barbell", "compendiumVersion": 1,
	})

	t.Run("found", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/equipment/1", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		var result models.UserEquipment
		json.Unmarshal(w.Body.Bytes(), &result)
		if result.CompendiumEquipmentID != "barbell" {
			t.Error("get response mismatch")
		}
	})

	t.Run("not found", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/equipment/999", nil)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})
}

func TestDeleteUserEquipment(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	doJSON(r, "POST", "/api/user/equipment", map[string]any{
		"owner": "alice", "compendiumEquipmentId": "barbell", "compendiumVersion": 1,
	})

	t.Run("success", func(t *testing.T) {
		w := doJSON(r, "DELETE", "/api/user/equipment/1", nil)
		if w.Code != http.StatusNoContent {
			t.Errorf("expected 204, got %d", w.Code)
		}
	})

	t.Run("db error", func(t *testing.T) {
		closeDB(t)
		w := doJSON(r, "DELETE", "/api/user/equipment/1", nil)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404 (db closed), got %d", w.Code)
		}
	})
}
