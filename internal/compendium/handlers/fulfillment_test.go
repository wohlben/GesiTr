package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"gesitr/internal/compendium/models"
)

func TestListFulfillments(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	t.Run("empty", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/fulfillments", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
	})

	doJSON(r, "POST", "/api/fulfillments", map[string]any{
		"equipmentTemplateId": "eq1", "fulfillsEquipmentTemplateId": "eq2", "createdBy": "system",
	})
	doJSON(r, "POST", "/api/fulfillments", map[string]any{
		"equipmentTemplateId": "eq3", "fulfillsEquipmentTemplateId": "eq4", "createdBy": "system",
	})

	t.Run("all", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/fulfillments", nil)
		var result []models.Fulfillment
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 2 {
			t.Errorf("expected 2, got %d", len(result))
		}
	})

	t.Run("filter equipmentTemplateId", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/fulfillments?equipmentTemplateId=eq1", nil)
		var result []models.Fulfillment
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 1 {
			t.Errorf("expected 1, got %d", len(result))
		}
	})

	t.Run("filter fulfillsEquipmentTemplateId", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/fulfillments?fulfillsEquipmentTemplateId=eq4", nil)
		var result []models.Fulfillment
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 1 {
			t.Errorf("expected 1, got %d", len(result))
		}
	})

	t.Run("db error", func(t *testing.T) {
		closeDB(t)
		w := doJSON(r, "GET", "/api/fulfillments", nil)
		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", w.Code)
		}
	})
}

func TestCreateFulfillment(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	t.Run("success", func(t *testing.T) {
		w := doJSON(r, "POST", "/api/fulfillments", map[string]any{
			"equipmentTemplateId": "a", "fulfillsEquipmentTemplateId": "b", "createdBy": "system",
		})
		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d", w.Code)
		}
		var result models.Fulfillment
		json.Unmarshal(w.Body.Bytes(), &result)
		if result.ID == 0 {
			t.Error("expected non-zero ID")
		}
	})

	t.Run("bad json", func(t *testing.T) {
		w := doRaw(r, "POST", "/api/fulfillments", "{bad")
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("db error", func(t *testing.T) {
		closeDB(t)
		w := doJSON(r, "POST", "/api/fulfillments", map[string]any{
			"equipmentTemplateId": "x", "fulfillsEquipmentTemplateId": "y", "createdBy": "s",
		})
		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", w.Code)
		}
	})
}

func TestDeleteFulfillment(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	doJSON(r, "POST", "/api/fulfillments", map[string]any{
		"equipmentTemplateId": "a", "fulfillsEquipmentTemplateId": "b", "createdBy": "system",
	})

	t.Run("success", func(t *testing.T) {
		w := doJSON(r, "DELETE", "/api/fulfillments/1", nil)
		if w.Code != http.StatusNoContent {
			t.Errorf("expected 204, got %d", w.Code)
		}
	})

	t.Run("db error", func(t *testing.T) {
		closeDB(t)
		w := doJSON(r, "DELETE", "/api/fulfillments/1", nil)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404 (db closed), got %d", w.Code)
		}
	})
}
