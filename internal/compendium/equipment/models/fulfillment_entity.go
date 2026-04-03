package models

import (
	"gesitr/internal/shared"
)

type FulfillmentEntity struct {
	shared.BaseModel
	EquipmentID         uint   `gorm:"not null;uniqueIndex:idx_fulfillment_pair"`
	FulfillsEquipmentID uint   `gorm:"not null;uniqueIndex:idx_fulfillment_pair"`
	Owner               string `gorm:"not null"`
}

func (FulfillmentEntity) TableName() string { return "fulfillments" }

func (e *FulfillmentEntity) ToDTO() Fulfillment {
	return Fulfillment{
		BaseModel:           e.BaseModel,
		EquipmentID:         e.EquipmentID,
		FulfillsEquipmentID: e.FulfillsEquipmentID,
		Owner:               e.Owner,
	}
}

func FulfillmentFromDTO(dto Fulfillment) FulfillmentEntity {
	return FulfillmentEntity{
		BaseModel:           dto.BaseModel,
		EquipmentID:         dto.EquipmentID,
		FulfillsEquipmentID: dto.FulfillsEquipmentID,
		Owner:               dto.Owner,
	}
}
