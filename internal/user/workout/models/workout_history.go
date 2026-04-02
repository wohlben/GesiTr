package models

import (
	"encoding/json"
	"time"

	profilemodels "gesitr/internal/profile/models"
	"gesitr/internal/shared"
)

type WorkoutHistoryEntity struct {
	ID               uint                             `gorm:"primaryKey;autoIncrement"`
	WorkoutID        uint                             `gorm:"not null;index:idx_workout_history"`
	Version          int                              `gorm:"not null"`
	Snapshot         string                           `gorm:"type:text;not null"`
	ChangedAt        time.Time                        `gorm:"not null"`
	ChangedBy        string                           `gorm:"not null"`
	ChangedByProfile *profilemodels.UserProfileEntity `gorm:"foreignKey:ChangedBy;references:ID;constraint:OnDelete:RESTRICT" json:"-"`
}

func (WorkoutHistoryEntity) TableName() string { return "workout_history" }

func (e *WorkoutHistoryEntity) ToVersionEntry() shared.VersionEntry {
	return shared.VersionEntry{
		Version:   e.Version,
		Snapshot:  json.RawMessage(e.Snapshot),
		ChangedAt: e.ChangedAt,
		ChangedBy: e.ChangedBy,
	}
}

// WorkoutChanged compares two Workout DTOs, ignoring BaseModel and Version.
func WorkoutChanged(old, new Workout) bool {
	return normalizeWorkoutJSON(old) != normalizeWorkoutJSON(new)
}

func normalizeWorkoutJSON(dto Workout) string {
	dto.BaseModel = shared.BaseModel{}
	dto.Version = 0
	dto.WorkoutGroup = nil

	if dto.Sections == nil {
		dto.Sections = []WorkoutSection{}
	}

	data, _ := json.Marshal(dto)
	return string(data)
}
