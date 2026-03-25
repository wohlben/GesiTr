package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"gesitr/internal/shared"
)

func newGroupPayload(name string) map[string]any {
	return map[string]any{
		"name":        name,
		"description": "test group",
	}
}

func TestGetExerciseGroupPermissions(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	// Create an exercise group owned by testuser
	doJSON(r, "POST", "/api/exercise-groups", newGroupPayload("Push"))

	t.Run("owner gets full permissions", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/exercise-groups/1/permissions", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
		}
		var resp shared.PermissionsResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		if len(resp.Permissions) != 3 {
			t.Fatalf("expected 3 permissions, got %d: %v", len(resp.Permissions), resp.Permissions)
		}
	})

	t.Run("non-owner gets empty permissions", func(t *testing.T) {
		w := doJSONAs(r, "GET", "/api/exercise-groups/1/permissions", nil, "bob")
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d, body = %s", w.Code, w.Body.String())
		}
		var resp shared.PermissionsResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		if len(resp.Permissions) != 0 {
			t.Errorf("expected empty permissions, got %v", resp.Permissions)
		}
	})

	t.Run("not found", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/exercise-groups/999/permissions", nil)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})
}

func TestGetExerciseGroupPermissions_VerifyValues(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	doJSON(r, "POST", "/api/exercise-groups", newGroupPayload("Pull"))

	w := doJSON(r, "GET", "/api/exercise-groups/1/permissions", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
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
}
