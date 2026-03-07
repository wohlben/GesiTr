package models

type FulfillmentEntity struct {
	BaseModel
	EquipmentTemplateID         string `gorm:"not null;uniqueIndex:idx_fulfillment_pair"`
	FulfillsEquipmentTemplateID string `gorm:"not null;uniqueIndex:idx_fulfillment_pair"`
	CreatedBy                   string `gorm:"not null"`
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
