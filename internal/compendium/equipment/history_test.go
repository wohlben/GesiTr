package equipment

import (
	"encoding/json"
	"testing"
	"time"

	"gesitr/internal/shared"
)

func TestEquipmentHistoryEntityTableName(t *testing.T) {
	if got := (EquipmentHistoryEntity{}).TableName(); got != "equipment_history" {
		t.Errorf("TableName() = %q, want %q", got, "equipment_history")
	}
}

func TestEquipmentChanged(t *testing.T) {
	base := Equipment{
		Name: "barbell", DisplayName: "Barbell", Description: "A bar",
		Category: EquipmentCategoryFreeWeights, TemplateID: "barbell", CreatedBy: "system",
	}

	t.Run("identical", func(t *testing.T) {
		if EquipmentChanged(base, base) {
			t.Error("expected no change for identical DTOs")
		}
	})

	t.Run("version ignored", func(t *testing.T) {
		other := base
		other.Version = 10
		if EquipmentChanged(base, other) {
			t.Error("expected no change when only version differs")
		}
	})

	t.Run("name changed", func(t *testing.T) {
		other := base
		other.Name = "ez-bar"
		if !EquipmentChanged(base, other) {
			t.Error("expected change for different name")
		}
	})

	t.Run("category changed", func(t *testing.T) {
		other := base
		other.Category = EquipmentCategoryMachines
		if !EquipmentChanged(base, other) {
			t.Error("expected change for different category")
		}
	})
}

func TestSnapshotJSON(t *testing.T) {
	dto := Equipment{Name: "test", TemplateID: "t1"}
	s := shared.SnapshotJSON(dto)
	if s == "" || s == "{}" {
		t.Error("expected non-empty snapshot")
	}
	// Verify it's valid JSON that round-trips correctly
	var roundTrip Equipment
	if err := json.Unmarshal([]byte(s), &roundTrip); err != nil {
		t.Fatalf("snapshot is not valid JSON: %v", err)
	}
	if roundTrip.Name != "test" || roundTrip.TemplateID != "t1" {
		t.Errorf("round-trip mismatch: name=%q templateId=%q", roundTrip.Name, roundTrip.TemplateID)
	}
}

func TestEquipmentHistoryToVersionEntry(t *testing.T) {
	now := time.Now()
	snapshot := `{"name":"barbell","version":0}`
	h := &EquipmentHistoryEntity{
		ID: 1, EquipmentID: 10, Version: 0,
		Snapshot: snapshot, ChangedAt: now, ChangedBy: "admin",
	}
	entry := h.ToVersionEntry()
	if entry.Version != 0 {
		t.Errorf("Version = %d", entry.Version)
	}
	if entry.ChangedBy != "admin" {
		t.Errorf("ChangedBy = %q", entry.ChangedBy)
	}
	if string(entry.Snapshot) != snapshot {
		t.Errorf("Snapshot = %q, want %q", string(entry.Snapshot), snapshot)
	}
}

func TestEquipmentChangedOptionalFields(t *testing.T) {
	base := Equipment{
		Name: "test", DisplayName: "Test", TemplateID: "t1", CreatedBy: "system",
	}

	t.Run("imageUrl changed", func(t *testing.T) {
		url := "http://example.com/img.png"
		other := base
		other.ImageUrl = &url
		if !EquipmentChanged(base, other) {
			t.Error("expected change when imageUrl added")
		}
	})

	t.Run("description changed", func(t *testing.T) {
		other := base
		other.Description = "new"
		if !EquipmentChanged(base, other) {
			t.Error("expected change for different description")
		}
	})

	t.Run("displayName changed", func(t *testing.T) {
		other := base
		other.DisplayName = "Different"
		if !EquipmentChanged(base, other) {
			t.Error("expected change for different displayName")
		}
	})

	t.Run("createdBy changed", func(t *testing.T) {
		other := base
		other.CreatedBy = "other-user"
		if !EquipmentChanged(base, other) {
			t.Error("expected change for different createdBy")
		}
	})
}
