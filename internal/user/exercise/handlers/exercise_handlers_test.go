package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"gesitr/internal/database"
	"gesitr/internal/user/exercise/models"
)

func TestListUserExercises(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	t.Run("empty list", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/exercises", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		var result []models.UserExercise
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 0 {
			t.Errorf("expected 0, got %d", len(result))
		}
	})

	// Seed data: one via handler (owner set from auth = alice), one directly for bob
	doJSON(r, "POST", "/api/user/exercises", map[string]any{
		"compendiumExerciseId": "bench-press", "compendiumVersion": 1,
	})
	database.DB.Create(&models.UserExerciseEntity{Owner: "bob", CompendiumExerciseID: "squat", CompendiumVersion: 2})

	t.Run("list all", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/exercises", nil)
		var result []models.UserExercise
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 2 {
			t.Errorf("expected 2, got %d", len(result))
		}
	})

	t.Run("filter by owner", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/exercises?owner=alice", nil)
		var result []models.UserExercise
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 1 || result[0].Owner != "alice" {
			t.Errorf("owner filter: got %d results", len(result))
		}
	})

	t.Run("filter by compendiumExerciseId", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/exercises?compendiumExerciseId=squat", nil)
		var result []models.UserExercise
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 1 || result[0].CompendiumExerciseID != "squat" {
			t.Errorf("templateId filter: got %d results", len(result))
		}
	})

	t.Run("db error", func(t *testing.T) {
		closeDB(t)
		w := doJSON(r, "GET", "/api/user/exercises", nil)
		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", w.Code)
		}
	})
}

func TestCreateUserExercise(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	t.Run("success", func(t *testing.T) {
		w := doJSON(r, "POST", "/api/user/exercises", map[string]any{
			"owner": "alice", "compendiumExerciseId": "deadlift", "compendiumVersion": 3,
		})
		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
		}
		var result models.UserExercise
		json.Unmarshal(w.Body.Bytes(), &result)
		if result.ID == 0 || result.Owner != "alice" || result.CompendiumExerciseID != "deadlift" || result.CompendiumVersion != 3 {
			t.Error("create response mismatch")
		}
	})

	t.Run("duplicate import", func(t *testing.T) {
		w := doJSON(r, "POST", "/api/user/exercises", map[string]any{
			"owner": "alice", "compendiumExerciseId": "deadlift", "compendiumVersion": 3,
		})
		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected 500 for duplicate, got %d", w.Code)
		}
	})

	t.Run("bad json", func(t *testing.T) {
		w := doRaw(r, "POST", "/api/user/exercises", "{invalid")
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("db error", func(t *testing.T) {
		closeDB(t)
		w := doJSON(r, "POST", "/api/user/exercises", map[string]any{
			"owner": "x", "compendiumExerciseId": "x", "compendiumVersion": 0,
		})
		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", w.Code)
		}
	})
}

func TestGetUserExercise(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	doJSON(r, "POST", "/api/user/exercises", map[string]any{
		"owner": "alice", "compendiumExerciseId": "bench-press", "compendiumVersion": 1,
	})

	t.Run("found", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/exercises/1", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		var result models.UserExercise
		json.Unmarshal(w.Body.Bytes(), &result)
		if result.CompendiumExerciseID != "bench-press" {
			t.Error("get response mismatch")
		}
	})

	t.Run("not found", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/exercises/999", nil)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})
}

func TestDeleteUserExercise(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	doJSON(r, "POST", "/api/user/exercises", map[string]any{
		"owner": "alice", "compendiumExerciseId": "bench-press", "compendiumVersion": 1,
	})

	t.Run("success", func(t *testing.T) {
		w := doJSON(r, "DELETE", "/api/user/exercises/1", nil)
		if w.Code != http.StatusNoContent {
			t.Errorf("expected 204, got %d", w.Code)
		}
	})

	t.Run("db error", func(t *testing.T) {
		closeDB(t)
		w := doJSON(r, "DELETE", "/api/user/exercises/1", nil)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404 (db closed), got %d", w.Code)
		}
	})
}
