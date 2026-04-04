package models

import (
	"gesitr/internal/shared"
)

type LocalityAvailabilityEntity struct {
	shared.BaseModel
	LocalityID       uint `gorm:"not null;uniqueIndex:idx_locality_equipment"`
	EquipmentID      uint `gorm:"not null;uniqueIndex:idx_locality_equipment"`
	Available        bool `gorm:"not null;default:true"`
	OwnershipGroupID uint
	EquipmentName    string `gorm:"->"` // populated via join, not stored
}

func (LocalityAvailabilityEntity) TableName() string { return "locality_availabilities" }

func (e *LocalityAvailabilityEntity) ToDTO() LocalityAvailability {
	return LocalityAvailability{
		BaseModel:        e.BaseModel,
		LocalityID:       e.LocalityID,
		EquipmentID:      e.EquipmentID,
		Available:        e.Available,
		OwnershipGroupID: e.OwnershipGroupID,
		EquipmentName:    e.EquipmentName,
	}
}

func LocalityAvailabilityFromDTO(dto LocalityAvailability) LocalityAvailabilityEntity {
	return LocalityAvailabilityEntity{
		BaseModel:        dto.BaseModel,
		LocalityID:       dto.LocalityID,
		EquipmentID:      dto.EquipmentID,
		Available:        dto.Available,
		OwnershipGroupID: dto.OwnershipGroupID,
	}
}
