package handlers_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"gesitr/internal/database"
	"gesitr/internal/user/workoutlog/models"
)

func TestCreateProposedLog(t *testing.T) {
	setupTestDB(t)
	defer closeDB(t)
	r := newRouter()

	dueStart := time.Now().Add(24 * time.Hour)
	dueEnd := time.Now().Add(48 * time.Hour)

	w := doJSON(r, "POST", "/api/user/workout-logs", map[string]any{
		"name":     "Proposed Workout",
		"status":   "proposed",
		"dueStart": dueStart.Format(time.RFC3339),
		"dueEnd":   dueEnd.Format(time.RFC3339),
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var log models.WorkoutLog
	json.Unmarshal(w.Body.Bytes(), &log)
	if log.Status != models.WorkoutLogStatusProposed {
		t.Errorf("expected proposed status, got %s", log.Status)
	}
	if log.DueStart == nil || log.DueEnd == nil {
		t.Error("expected dueStart and dueEnd to be set")
	}
}

func TestCreateCommittedLog(t *testing.T) {
	setupTestDB(t)
	defer closeDB(t)
	r := newRouter()

	dueStart := time.Now().Add(24 * time.Hour)
	dueEnd := time.Now().Add(48 * time.Hour)

	w := doJSON(r, "POST", "/api/user/workout-logs", map[string]any{
		"name":     "Committed Workout",
		"status":   "committed",
		"dueStart": dueStart.Format(time.RFC3339),
		"dueEnd":   dueEnd.Format(time.RFC3339),
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var log models.WorkoutLog
	json.Unmarshal(w.Body.Bytes(), &log)
	if log.Status != models.WorkoutLogStatusCommitted {
		t.Errorf("expected committed status, got %s", log.Status)
	}
}

func TestCreateCommitmentRequiresDueWindow(t *testing.T) {
	setupTestDB(t)
	defer closeDB(t)
	r := newRouter()

	// Missing dueStart/dueEnd
	w := doJSON(r, "POST", "/api/user/workout-logs", map[string]any{
		"name":   "Bad Commitment",
		"status": "proposed",
	})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for missing due window, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCreateCommitmentInvalidDueWindow(t *testing.T) {
	setupTestDB(t)
	defer closeDB(t)
	r := newRouter()

	// dueEnd before dueStart
	dueStart := time.Now().Add(48 * time.Hour)
	dueEnd := time.Now().Add(24 * time.Hour)

	w := doJSON(r, "POST", "/api/user/workout-logs", map[string]any{
		"name":     "Bad Window",
		"status":   "proposed",
		"dueStart": dueStart.Format(time.RFC3339),
		"dueEnd":   dueEnd.Format(time.RFC3339),
	})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid due window, got %d: %s", w.Code, w.Body.String())
	}
}

func TestSkipProposedLog(t *testing.T) {
	setupTestDB(t)
	defer closeDB(t)
	r := newRouter()

	dueStart := time.Now().Add(24 * time.Hour)
	dueEnd := time.Now().Add(48 * time.Hour)

	doJSON(r, "POST", "/api/user/workout-logs", map[string]any{
		"name":     "To Skip",
		"status":   "proposed",
		"dueStart": dueStart.Format(time.RFC3339),
		"dueEnd":   dueEnd.Format(time.RFC3339),
	})

	w := doJSON(r, "POST", "/api/user/workout-logs/1/skip", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var log models.WorkoutLog
	json.Unmarshal(w.Body.Bytes(), &log)
	if log.Status != models.WorkoutLogStatusSkipped {
		t.Errorf("expected skipped status, got %s", log.Status)
	}
}

func TestSkipCommittedLogForbidden(t *testing.T) {
	setupTestDB(t)
	defer closeDB(t)
	r := newRouter()

	dueStart := time.Now().Add(24 * time.Hour)
	dueEnd := time.Now().Add(48 * time.Hour)

	doJSON(r, "POST", "/api/user/workout-logs", map[string]any{
		"name":     "Committed",
		"status":   "committed",
		"dueStart": dueStart.Format(time.RFC3339),
		"dueEnd":   dueEnd.Format(time.RFC3339),
	})

	w := doJSON(r, "POST", "/api/user/workout-logs/1/skip", nil)
	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409 (can't skip committed), got %d: %s", w.Code, w.Body.String())
	}
}

func TestCommitProposedLog(t *testing.T) {
	setupTestDB(t)
	defer closeDB(t)
	r := newRouter()

	dueStart := time.Now().Add(24 * time.Hour)
	dueEnd := time.Now().Add(48 * time.Hour)

	doJSON(r, "POST", "/api/user/workout-logs", map[string]any{
		"name":     "To Commit",
		"status":   "proposed",
		"dueStart": dueStart.Format(time.RFC3339),
		"dueEnd":   dueEnd.Format(time.RFC3339),
	})

	w := doJSON(r, "POST", "/api/user/workout-logs/1/commit", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var log models.WorkoutLog
	json.Unmarshal(w.Body.Bytes(), &log)
	if log.Status != models.WorkoutLogStatusCommitted {
		t.Errorf("expected committed status, got %s", log.Status)
	}
}

func TestStartCommittedLog(t *testing.T) {
	setupTestDB(t)
	defer closeDB(t)
	r := newRouter()

	dueStart := time.Now().Add(24 * time.Hour)
	dueEnd := time.Now().Add(48 * time.Hour)

	doJSON(r, "POST", "/api/user/workout-logs", map[string]any{
		"name":     "To Start",
		"status":   "committed",
		"dueStart": dueStart.Format(time.RFC3339),
		"dueEnd":   dueEnd.Format(time.RFC3339),
	})

	w := doJSON(r, "POST", "/api/user/workout-logs/1/start", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var log models.WorkoutLog
	json.Unmarshal(w.Body.Bytes(), &log)
	if log.Status != models.WorkoutLogStatusInProgress {
		t.Errorf("expected in_progress status, got %s", log.Status)
	}
}

func TestDeleteProposedLog(t *testing.T) {
	setupTestDB(t)
	defer closeDB(t)
	r := newRouter()

	dueStart := time.Now().Add(24 * time.Hour)
	dueEnd := time.Now().Add(48 * time.Hour)

	doJSON(r, "POST", "/api/user/workout-logs", map[string]any{
		"name":     "To Delete",
		"status":   "proposed",
		"dueStart": dueStart.Format(time.RFC3339),
		"dueEnd":   dueEnd.Format(time.RFC3339),
	})

	w := doJSON(r, "DELETE", "/api/user/workout-logs/1", nil)
	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d: %s", w.Code, w.Body.String())
	}
}

func TestDeleteCommittedLog(t *testing.T) {
	setupTestDB(t)
	defer closeDB(t)
	r := newRouter()

	dueStart := time.Now().Add(24 * time.Hour)
	dueEnd := time.Now().Add(48 * time.Hour)

	doJSON(r, "POST", "/api/user/workout-logs", map[string]any{
		"name":     "To Delete",
		"status":   "committed",
		"dueStart": dueStart.Format(time.RFC3339),
		"dueEnd":   dueEnd.Format(time.RFC3339),
	})

	w := doJSON(r, "DELETE", "/api/user/workout-logs/1", nil)
	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d: %s", w.Code, w.Body.String())
	}
}

func TestMultipleCommitmentsPerWorkout(t *testing.T) {
	setupTestDB(t)
	defer closeDB(t)
	r := newRouter()

	// Create a workout
	doJSON(r, "POST", "/api/user/workouts", map[string]any{"name": "My Workout"})

	dueStart1 := time.Now().Add(24 * time.Hour)
	dueEnd1 := time.Now().Add(48 * time.Hour)
	dueStart2 := time.Now().Add(72 * time.Hour)
	dueEnd2 := time.Now().Add(96 * time.Hour)

	// Create two committed logs for the same workout
	w1 := doJSON(r, "POST", "/api/user/workout-logs", map[string]any{
		"name":      "Commitment 1",
		"workoutId": 1,
		"status":    "committed",
		"dueStart":  dueStart1.Format(time.RFC3339),
		"dueEnd":    dueEnd1.Format(time.RFC3339),
	})
	if w1.Code != http.StatusCreated {
		t.Fatalf("first commitment: expected 201, got %d: %s", w1.Code, w1.Body.String())
	}

	w2 := doJSON(r, "POST", "/api/user/workout-logs", map[string]any{
		"name":      "Commitment 2",
		"workoutId": 1,
		"status":    "committed",
		"dueStart":  dueStart2.Format(time.RFC3339),
		"dueEnd":    dueEnd2.Format(time.RFC3339),
	})
	if w2.Code != http.StatusCreated {
		t.Fatalf("second commitment: expected 201, got %d: %s", w2.Code, w2.Body.String())
	}
}

func TestCommittedToBrokenByTicker(t *testing.T) {
	setupTestDB(t)
	defer closeDB(t)

	// Create a committed log with due_end in the past
	pastDueEnd := time.Now().Add(-1 * time.Hour)
	pastDueStart := time.Now().Add(-25 * time.Hour)
	entity := models.WorkoutLogEntity{
		Owner:    "alice",
		Name:     "Overdue",
		Status:   models.WorkoutLogStatusCommitted,
		DueStart: &pastDueStart,
		DueEnd:   &pastDueEnd,
	}
	database.DB.Create(&entity)

	// Simulate what the ticker does
	now := time.Now()
	database.DB.Model(&models.WorkoutLogEntity{}).
		Where("status = ? AND due_end < ?", models.WorkoutLogStatusCommitted, now).
		Updates(map[string]any{
			"status":            models.WorkoutLogStatusBroken,
			"status_changed_at": now,
		})

	// Verify status changed
	var updated models.WorkoutLogEntity
	database.DB.First(&updated, entity.ID)
	if updated.Status != models.WorkoutLogStatusBroken {
		t.Errorf("expected broken status, got %s", updated.Status)
	}
}

func TestCommittedNotBrokenWhenDueEndFuture(t *testing.T) {
	setupTestDB(t)
	defer closeDB(t)

	// Create a committed log with due_end in the future
	futureDueStart := time.Now().Add(24 * time.Hour)
	futureDueEnd := time.Now().Add(48 * time.Hour)
	entity := models.WorkoutLogEntity{
		Owner:    "alice",
		Name:     "Future",
		Status:   models.WorkoutLogStatusCommitted,
		DueStart: &futureDueStart,
		DueEnd:   &futureDueEnd,
	}
	database.DB.Create(&entity)

	// Run the ticker logic
	now := time.Now()
	database.DB.Model(&models.WorkoutLogEntity{}).
		Where("status = ? AND due_end < ?", models.WorkoutLogStatusCommitted, now).
		Updates(map[string]any{
			"status":            models.WorkoutLogStatusBroken,
			"status_changed_at": now,
		})

	// Verify status didn't change
	var updated models.WorkoutLogEntity
	database.DB.First(&updated, entity.ID)
	if updated.Status != models.WorkoutLogStatusCommitted {
		t.Errorf("expected committed status (future due_end), got %s", updated.Status)
	}
}

func TestFullCommitmentFlow(t *testing.T) {
	setupTestDB(t)
	defer closeDB(t)
	r := newRouter()

	dueStart := time.Now().Add(24 * time.Hour)
	dueEnd := time.Now().Add(48 * time.Hour)

	// 1. Create proposed
	w := doJSON(r, "POST", "/api/user/workout-logs", map[string]any{
		"name":     "Full Flow",
		"status":   "proposed",
		"dueStart": dueStart.Format(time.RFC3339),
		"dueEnd":   dueEnd.Format(time.RFC3339),
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create: expected 201, got %d", w.Code)
	}

	// 2. Commit
	w = doJSON(r, "POST", "/api/user/workout-logs/1/commit", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("commit: expected 200, got %d", w.Code)
	}

	// 3. Start
	w = doJSON(r, "POST", "/api/user/workout-logs/1/start", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("start: expected 200, got %d", w.Code)
	}

	var log models.WorkoutLog
	json.Unmarshal(w.Body.Bytes(), &log)
	if log.Status != models.WorkoutLogStatusInProgress {
		t.Errorf("expected in_progress, got %s", log.Status)
	}
	if log.DueStart == nil || log.DueEnd == nil {
		t.Error("due window should be preserved after start")
	}
}
