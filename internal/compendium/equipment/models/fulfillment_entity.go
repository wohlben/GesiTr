package models

import (
	"gesitr/internal/shared"
)

type FulfillmentEntity struct {
	shared.BaseModel
	EquipmentID         uint `gorm:"not null;uniqueIndex:idx_fulfillment_pair"`
	FulfillsEquipmentID uint `gorm:"not null;uniqueIndex:idx_fulfillment_pair"`
	OwnershipGroupID    uint
}

func (FulfillmentEntity) TableName() string { return "fulfillments" }

func (e *FulfillmentEntity) ToDTO() Fulfillment {
	return Fulfillment{
		BaseModel:           e.BaseModel,
		EquipmentID:         e.EquipmentID,
		FulfillsEquipmentID: e.FulfillsEquipmentID,
		OwnershipGroupID:    e.OwnershipGroupID,
	}
}

func FulfillmentFromDTO(dto Fulfillment) FulfillmentEntity {
	return FulfillmentEntity{
		BaseModel:           dto.BaseModel,
		EquipmentID:         dto.EquipmentID,
		FulfillsEquipmentID: dto.FulfillsEquipmentID,
		OwnershipGroupID:    dto.OwnershipGroupID,
	}
}
