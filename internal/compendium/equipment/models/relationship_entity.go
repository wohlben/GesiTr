package models

import (
	"gesitr/internal/shared"
)

type EquipmentRelationshipEntity struct {
	shared.BaseModel
	RelationshipType EquipmentRelationshipType `gorm:"not null;uniqueIndex:idx_equipment_relationship"`
	Strength         float64                   `gorm:"not null"`
	OwnershipGroupID uint                      `gorm:"uniqueIndex:idx_equipment_relationship"`
	FromEquipmentID  uint                      `gorm:"not null;uniqueIndex:idx_equipment_relationship"`
	ToEquipmentID    uint                      `gorm:"not null;uniqueIndex:idx_equipment_relationship"`
}

func (EquipmentRelationshipEntity) TableName() string { return "equipment_relationships" }

func (e *EquipmentRelationshipEntity) ToDTO() EquipmentRelationship {
	return EquipmentRelationship{
		BaseModel:        e.BaseModel,
		RelationshipType: e.RelationshipType,
		Strength:         e.Strength,
		OwnershipGroupID: e.OwnershipGroupID,
		FromEquipmentID:  e.FromEquipmentID,
		ToEquipmentID:    e.ToEquipmentID,
	}
}

func EquipmentRelationshipFromDTO(dto EquipmentRelationship) EquipmentRelationshipEntity {
	return EquipmentRelationshipEntity{
		BaseModel:        dto.BaseModel,
		RelationshipType: dto.RelationshipType,
		Strength:         dto.Strength,
		OwnershipGroupID: dto.OwnershipGroupID,
		FromEquipmentID:  dto.FromEquipmentID,
		ToEquipmentID:    dto.ToEquipmentID,
	}
}
