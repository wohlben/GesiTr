package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"gesitr/internal/compendium/models"
)

func TestListEquipment(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	t.Run("empty list", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/equipment", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		var result []models.Equipment
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 0 {
			t.Errorf("expected 0, got %d", len(result))
		}
	})

	// Seed data for filter tests
	doJSON(r, "POST", "/api/equipment", map[string]any{
		"name": "barbell", "displayName": "Barbell", "description": "A bar",
		"category": "free_weights", "templateId": "barbell", "createdBy": "system",
	})
	doJSON(r, "POST", "/api/equipment", map[string]any{
		"name": "bench", "displayName": "Flat Bench", "description": "A bench",
		"category": "benches", "templateId": "bench", "createdBy": "system",
	})

	t.Run("list all", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/equipment", nil)
		var result []models.Equipment
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 2 {
			t.Errorf("expected 2, got %d", len(result))
		}
	})

	t.Run("filter by q name", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/equipment?q=bar", nil)
		var result []models.Equipment
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 1 || result[0].Name != "barbell" {
			t.Errorf("q filter: got %d results", len(result))
		}
	})

	t.Run("filter by q displayName", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/equipment?q=Flat", nil)
		var result []models.Equipment
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 1 || result[0].Name != "bench" {
			t.Errorf("q displayName filter: got %d results", len(result))
		}
	})

	t.Run("filter by category", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/equipment?category=benches", nil)
		var result []models.Equipment
		json.Unmarshal(w.Body.Bytes(), &result)
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

func TestCreateEquipment(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	t.Run("success", func(t *testing.T) {
		w := doJSON(r, "POST", "/api/equipment", map[string]any{
			"name": "dumbbell", "displayName": "Dumbbell", "description": "A weight",
			"category": "free_weights", "templateId": "dumbbell", "createdBy": "system",
		})
		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
		}
		var result models.Equipment
		json.Unmarshal(w.Body.Bytes(), &result)
		if result.ID == 0 || result.Name != "dumbbell" {
			t.Error("create response mismatch")
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
			"category": "other", "templateId": "x", "createdBy": "s",
		})
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
		"category": "free_weights", "templateId": "kb", "createdBy": "system",
	})
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
		"category": "accessories", "templateId": "band", "createdBy": "system",
	})

	t.Run("success", func(t *testing.T) {
		w := doJSON(r, "PUT", "/api/equipment/1", map[string]any{
			"name": "resistance-band", "displayName": "Resistance Band", "description": "elastic",
			"category": "accessories", "templateId": "band", "createdBy": "system",
		})
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
		}
		var result models.Equipment
		json.Unmarshal(w.Body.Bytes(), &result)
		if result.Name != "resistance-band" || result.Version != 1 {
			t.Errorf("update response: name=%q version=%d", result.Name, result.Version)
		}
	})

	t.Run("not found", func(t *testing.T) {
		w := doJSON(r, "PUT", "/api/equipment/999", map[string]any{
			"name": "x", "displayName": "X", "description": "",
			"category": "other", "templateId": "x", "createdBy": "s",
		})
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

	t.Run("save error unique constraint", func(t *testing.T) {
		// Create a second equipment
		doJSON(r, "POST", "/api/equipment", map[string]any{
			"name": "other", "displayName": "Other", "description": "",
			"category": "other", "templateId": "other-tid", "createdBy": "system",
		})
		// Update second equipment with first's templateId → unique violation on Save
		w := doJSON(r, "PUT", "/api/equipment/2", map[string]any{
			"name": "conflict", "displayName": "Conflict", "description": "",
			"category": "other", "templateId": "band", "createdBy": "system",
		})
		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected 500 for unique violation, got %d", w.Code)
		}
	})

	t.Run("db error first lookup", func(t *testing.T) {
		closeDB(t)
		w := doJSON(r, "PUT", "/api/equipment/1", map[string]any{
			"name": "x", "displayName": "X", "description": "",
			"category": "other", "templateId": "x", "createdBy": "s",
		})
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404 (db closed), got %d", w.Code)
		}
	})
}

func TestDeleteEquipment(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	doJSON(r, "POST", "/api/equipment", map[string]any{
		"name": "rope", "displayName": "Rope", "description": "",
		"category": "accessories", "templateId": "rope", "createdBy": "system",
	})

	t.Run("success", func(t *testing.T) {
		w := doJSON(r, "DELETE", "/api/equipment/1", nil)
		if w.Code != http.StatusNoContent {
			t.Errorf("expected 204, got %d", w.Code)
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
