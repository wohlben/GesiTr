package exercise

import (
	"encoding/json"
	"testing"
	"time"

	"gesitr/internal/shared"
)

func TestExerciseHistoryEntityTableName(t *testing.T) {
	if got := (ExerciseHistoryEntity{}).TableName(); got != "exercise_history" {
		t.Errorf("TableName() = %q, want %q", got, "exercise_history")
	}
}

func TestExerciseChanged(t *testing.T) {
	base := Exercise{
		Name: "Bench Press", Slug: "bench-press", Type: ExerciseTypeStrength,
		TechnicalDifficulty: DifficultyIntermediate, CreatedBy: "system",
		Force: []Force{ForcePush}, PrimaryMuscles: []Muscle{MuscleChest},
	}

	t.Run("identical", func(t *testing.T) {
		if ExerciseChanged(base, base) {
			t.Error("expected no change for identical DTOs")
		}
	})

	t.Run("different version ignored", func(t *testing.T) {
		other := base
		other.Version = 5
		if ExerciseChanged(base, other) {
			t.Error("expected no change when only version differs")
		}
	})

	t.Run("different base model ignored", func(t *testing.T) {
		other := base
		other.ID = 99
		if ExerciseChanged(base, other) {
			t.Error("expected no change when only BaseModel differs")
		}
	})

	t.Run("name changed", func(t *testing.T) {
		other := base
		other.Name = "Incline Bench Press"
		if !ExerciseChanged(base, other) {
			t.Error("expected change for different name")
		}
	})

	t.Run("muscles changed", func(t *testing.T) {
		other := base
		other.PrimaryMuscles = []Muscle{MuscleFrontDelts}
		if !ExerciseChanged(base, other) {
			t.Error("expected change for different muscles")
		}
	})

	t.Run("slice order ignored for unordered fields", func(t *testing.T) {
		a := base
		a.Force = []Force{ForcePush, ForcePull}
		b := base
		b.Force = []Force{ForcePull, ForcePush}
		if ExerciseChanged(a, b) {
			t.Error("expected no change when only slice order differs for unordered field")
		}
	})

	t.Run("instruction order matters", func(t *testing.T) {
		a := base
		a.Instructions = []string{"Step 1", "Step 2"}
		b := base
		b.Instructions = []string{"Step 2", "Step 1"}
		if !ExerciseChanged(a, b) {
			t.Error("expected change when instruction order differs")
		}
	})

	t.Run("nil vs empty slice treated equal", func(t *testing.T) {
		a := base
		a.Instructions = nil
		b := base
		b.Instructions = []string{}
		if ExerciseChanged(a, b) {
			t.Error("expected no change for nil vs empty slice")
		}
	})
}

func TestSnapshotJSON(t *testing.T) {
	dto := Exercise{Name: "test", TemplateID: "t1"}
	s := shared.SnapshotJSON(dto)
	if s == "" || s == "{}" {
		t.Error("expected non-empty snapshot")
	}
	// Verify it's valid JSON that round-trips correctly
	var roundTrip Exercise
	if err := json.Unmarshal([]byte(s), &roundTrip); err != nil {
		t.Fatalf("snapshot is not valid JSON: %v", err)
	}
	if roundTrip.Name != "test" || roundTrip.TemplateID != "t1" {
		t.Errorf("round-trip mismatch: name=%q templateId=%q", roundTrip.Name, roundTrip.TemplateID)
	}
}

func TestExerciseHistoryToVersionEntry(t *testing.T) {
	now := time.Now()
	snapshot := `{"name":"Bench Press","version":1}`
	h := &ExerciseHistoryEntity{
		ID: 1, ExerciseID: 42, Version: 1,
		Snapshot: snapshot, ChangedAt: now, ChangedBy: "user1",
	}
	entry := h.ToVersionEntry()
	if entry.Version != 1 {
		t.Errorf("Version = %d", entry.Version)
	}
	if entry.ChangedBy != "user1" {
		t.Errorf("ChangedBy = %q", entry.ChangedBy)
	}
	if !entry.ChangedAt.Equal(now) {
		t.Errorf("ChangedAt mismatch")
	}
	// Snapshot should be raw JSON (not double-encoded)
	if string(entry.Snapshot) != snapshot {
		t.Errorf("Snapshot = %q, want %q", string(entry.Snapshot), snapshot)
	}
}

func TestExerciseChangedOptionalFields(t *testing.T) {
	base := Exercise{
		Name: "Test", Slug: "test", Type: ExerciseTypeStrength,
		CreatedBy: "system",
	}

	t.Run("authorName changed", func(t *testing.T) {
		author := "John"
		other := base
		other.AuthorName = &author
		if !ExerciseChanged(base, other) {
			t.Error("expected change when authorName added")
		}
	})

	t.Run("description changed", func(t *testing.T) {
		other := base
		other.Description = "new desc"
		if !ExerciseChanged(base, other) {
			t.Error("expected change for different description")
		}
	})

	t.Run("templateId changed", func(t *testing.T) {
		other := base
		other.TemplateID = "tmpl-1"
		if !ExerciseChanged(base, other) {
			t.Error("expected change when templateId changed")
		}
	})

	t.Run("equipmentIds changed", func(t *testing.T) {
		other := base
		other.EquipmentIDs = []string{"barbell"}
		if !ExerciseChanged(base, other) {
			t.Error("expected change when equipmentIds added")
		}
	})

	t.Run("equipmentIds order ignored", func(t *testing.T) {
		a := base
		a.EquipmentIDs = []string{"barbell", "bench"}
		b := base
		b.EquipmentIDs = []string{"bench", "barbell"}
		if ExerciseChanged(a, b) {
			t.Error("expected no change when only equipmentIds order differs")
		}
	})
}
