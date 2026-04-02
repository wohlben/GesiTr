package models

import (
	"encoding/json"
	"slices"
	"time"

	profilemodels "gesitr/internal/profile/models"
	"gesitr/internal/shared"
)

type ExerciseHistoryEntity struct {
	ID               uint                             `gorm:"primaryKey;autoIncrement"`
	ExerciseID       uint                             `gorm:"not null;index:idx_exercise_history"`
	Version          int                              `gorm:"not null"`
	Snapshot         string                           `gorm:"type:text;not null"`
	ChangedAt        time.Time                        `gorm:"not null"`
	ChangedBy        string                           `gorm:"not null"`
	ChangedByProfile *profilemodels.UserProfileEntity `gorm:"foreignKey:ChangedBy;references:ID;constraint:OnDelete:RESTRICT" json:"-"`
}

func (ExerciseHistoryEntity) TableName() string { return "exercise_history" }

func (e *ExerciseHistoryEntity) ToVersionEntry() shared.VersionEntry {
	return shared.VersionEntry{
		Version:   e.Version,
		Snapshot:  json.RawMessage(e.Snapshot),
		ChangedAt: e.ChangedAt,
		ChangedBy: e.ChangedBy,
	}
}

// ExerciseChanged compares two Exercise DTOs, ignoring BaseModel and Version.
func ExerciseChanged(old, new Exercise) bool {
	return normalizeExerciseJSON(old) != normalizeExerciseJSON(new)
}

func normalizeExerciseJSON(dto Exercise) string {
	dto.BaseModel = shared.BaseModel{}
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
		dto.EquipmentIDs = []uint{}
	}

	data, _ := json.Marshal(dto)
	return string(data)
}
