package models

import (
	profilemodels "gesitr/internal/profile/models"
	"gesitr/internal/shared"
)

type EquipmentRelationshipEntity struct {
	shared.BaseModel
	RelationshipType EquipmentRelationshipType        `gorm:"not null;uniqueIndex:idx_equipment_relationship"`
	Strength         float64                          `gorm:"not null"`
	Owner            string                           `gorm:"not null;uniqueIndex:idx_equipment_relationship"`
	OwnerProfile     *profilemodels.UserProfileEntity `gorm:"foreignKey:Owner;references:ID;constraint:OnDelete:RESTRICT" json:"-"`
	FromEquipmentID  uint                             `gorm:"not null;uniqueIndex:idx_equipment_relationship"`
	ToEquipmentID    uint                             `gorm:"not null;uniqueIndex:idx_equipment_relationship"`
}

func (EquipmentRelationshipEntity) TableName() string { return "equipment_relationships" }

func (e *EquipmentRelationshipEntity) ToDTO() EquipmentRelationship {
	return EquipmentRelationship{
		BaseModel:        e.BaseModel,
		RelationshipType: e.RelationshipType,
		Strength:         e.Strength,
		Owner:            e.Owner,
		FromEquipmentID:  e.FromEquipmentID,
		ToEquipmentID:    e.ToEquipmentID,
	}
}

func EquipmentRelationshipFromDTO(dto EquipmentRelationship) EquipmentRelationshipEntity {
	return EquipmentRelationshipEntity{
		BaseModel:        dto.BaseModel,
		RelationshipType: dto.RelationshipType,
		Strength:         dto.Strength,
		Owner:            dto.Owner,
		FromEquipmentID:  dto.FromEquipmentID,
		ToEquipmentID:    dto.ToEquipmentID,
	}
}
