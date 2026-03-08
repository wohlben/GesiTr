package models

import (
	"encoding/json"
	"slices"
	"time"
)

type ExerciseHistoryEntity struct {
	ID         uint      `gorm:"primaryKey;autoIncrement"`
	ExerciseID uint      `gorm:"not null;index:idx_exercise_history"`
	Version    int       `gorm:"not null"`
	Snapshot   string    `gorm:"type:text;not null"`
	ChangedAt  time.Time `gorm:"not null"`
	ChangedBy  string    `gorm:"not null"`
}

func (ExerciseHistoryEntity) TableName() string { return "exercise_history" }

type EquipmentHistoryEntity struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	EquipmentID uint      `gorm:"not null;index:idx_equipment_history"`
	Version     int       `gorm:"not null"`
	Snapshot    string    `gorm:"type:text;not null"`
	ChangedAt   time.Time `gorm:"not null"`
	ChangedBy   string    `gorm:"not null"`
}

func (EquipmentHistoryEntity) TableName() string { return "equipment_history" }

// VersionEntry is the API response DTO for a single version in the history.
type VersionEntry struct {
	Version   int             `json:"version"`
	Snapshot  json.RawMessage `json:"snapshot"`
	ChangedAt time.Time       `json:"changedAt"`
	ChangedBy string          `json:"changedBy"`
}

func (e *ExerciseHistoryEntity) ToVersionEntry() VersionEntry {
	return VersionEntry{
		Version:   e.Version,
		Snapshot:  json.RawMessage(e.Snapshot),
		ChangedAt: e.ChangedAt,
		ChangedBy: e.ChangedBy,
	}
}

func (e *EquipmentHistoryEntity) ToVersionEntry() VersionEntry {
	return VersionEntry{
		Version:   e.Version,
		Snapshot:  json.RawMessage(e.Snapshot),
		ChangedAt: e.ChangedAt,
		ChangedBy: e.ChangedBy,
	}
}

// SnapshotJSON marshals a DTO to JSON for history storage.
func SnapshotJSON(dto any) string {
	data, _ := json.Marshal(dto)
	return string(data)
}

// ExerciseChanged compares two Exercise DTOs, ignoring BaseModel and Version.
func ExerciseChanged(old, new Exercise) bool {
	return normalizeExerciseJSON(old) != normalizeExerciseJSON(new)
}

// EquipmentChanged compares two Equipment DTOs, ignoring BaseModel and Version.
func EquipmentChanged(old, new Equipment) bool {
	return normalizeEquipmentJSON(old) != normalizeEquipmentJSON(new)
}

func normalizeExerciseJSON(dto Exercise) string {
	dto.BaseModel = BaseModel{}
	dto.Version = 0

	// Sort unordered slices for consistent comparison
	slices.Sort(dto.Force)
	slices.Sort(dto.PrimaryMuscles)
	slices.Sort(dto.SecondaryMuscles)
	slices.Sort(dto.SuggestedMeasurementParadigms)
	slices.Sort(dto.AlternativeNames)
	slices.Sort(dto.EquipmentIDs)
	// Instructions and Images are position-ordered — don't sort

	// Normalize nil vs empty slices
	if dto.Force == nil {
		dto.Force = []Force{}
	}
	if dto.PrimaryMuscles == nil {
		dto.PrimaryMuscles = []Muscle{}
	}
	if dto.SecondaryMuscles == nil {
		dto.SecondaryMuscles = []Muscle{}
	}
	if dto.SuggestedMeasurementParadigms == nil {
		dto.SuggestedMeasurementParadigms = []MeasurementParadigm{}
	}
	if dto.Instructions == nil {
		dto.Instructions = []string{}
	}
	if dto.Images == nil {
		dto.Images = []string{}
	}
	if dto.AlternativeNames == nil {
		dto.AlternativeNames = []string{}
	}
	if dto.EquipmentIDs == nil {
		dto.EquipmentIDs = []string{}
	}

	data, _ := json.Marshal(dto)
	return string(data)
}

func normalizeEquipmentJSON(dto Equipment) string {
	dto.BaseModel = BaseModel{}
	dto.Version = 0

	data, _ := json.Marshal(dto)
	return string(data)
}
