package models

import (
	profilemodels "gesitr/internal/profile/models"
	"gesitr/internal/shared"
)

type FulfillmentEntity struct {
	shared.BaseModel
	EquipmentTemplateID         string                           `gorm:"not null;uniqueIndex:idx_fulfillment_pair"`
	FulfillsEquipmentTemplateID string                           `gorm:"not null;uniqueIndex:idx_fulfillment_pair"`
	CreatedBy                   string                           `gorm:"not null"`
	CreatedByProfile            *profilemodels.UserProfileEntity `gorm:"foreignKey:CreatedBy;references:ID;constraint:OnDelete:RESTRICT" json:"-"`
}

func (FulfillmentEntity) TableName() string { return "fulfillments" }

func (e *FulfillmentEntity) ToDTO() Fulfillment {
	return Fulfillment{
		BaseModel:                   e.BaseModel,
		EquipmentTemplateID:         e.EquipmentTemplateID,
		FulfillsEquipmentTemplateID: e.FulfillsEquipmentTemplateID,
		CreatedBy:                   e.CreatedBy,
	}
}

func FulfillmentFromDTO(dto Fulfillment) FulfillmentEntity {
	return FulfillmentEntity{
		BaseModel:                   dto.BaseModel,
		EquipmentTemplateID:         dto.EquipmentTemplateID,
		FulfillsEquipmentTemplateID: dto.FulfillsEquipmentTemplateID,
		CreatedBy:                   dto.CreatedBy,
	}
}
