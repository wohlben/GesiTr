package models

import (
	"encoding/json"
	"time"

	profilemodels "gesitr/internal/profile/models"
	"gesitr/internal/shared"
)

type EquipmentHistoryEntity struct {
	ID               uint                             `gorm:"primaryKey;autoIncrement"`
	EquipmentID      uint                             `gorm:"not null;index:idx_equipment_history"`
	Version          int                              `gorm:"not null"`
	Snapshot         string                           `gorm:"type:text;not null"`
	ChangedAt        time.Time                        `gorm:"not null"`
	ChangedBy        string                           `gorm:"not null"`
	ChangedByProfile *profilemodels.UserProfileEntity `gorm:"foreignKey:ChangedBy;references:ID;constraint:OnDelete:RESTRICT" json:"-"`
}

func (EquipmentHistoryEntity) TableName() string { return "equipment_history" }

func (e *EquipmentHistoryEntity) ToVersionEntry() shared.VersionEntry {
	return shared.VersionEntry{
		Version:   e.Version,
		Snapshot:  json.RawMessage(e.Snapshot),
		ChangedAt: e.ChangedAt,
		ChangedBy: e.ChangedBy,
	}
}

// EquipmentChanged compares two Equipment DTOs, ignoring BaseModel and Version.
func EquipmentChanged(old, new Equipment) bool {
	return normalizeEquipmentJSON(old) != normalizeEquipmentJSON(new)
}

func normalizeEquipmentJSON(dto Equipment) string {
	dto.BaseModel = shared.BaseModel{}
	dto.Version = 0

	data, _ := json.Marshal(dto)
	return string(data)
}
