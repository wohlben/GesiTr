package models

type UserEquipmentEntity struct {
	BaseModel
	Owner               string `gorm:"not null;index;uniqueIndex:idx_owner_equipment_template"`
	EquipmentTemplateID string `gorm:"not null;uniqueIndex:idx_owner_equipment_template"`
	CompendiumVersion   int    `gorm:"not null"`
}

func (UserEquipmentEntity) TableName() string { return "user_equipment" }

func (e *UserEquipmentEntity) ToDTO() UserEquipment {
	return UserEquipment{
		BaseModel:           e.BaseModel,
		Owner:               e.Owner,
		EquipmentTemplateID: e.EquipmentTemplateID,
		CompendiumVersion:   e.CompendiumVersion,
	}
}

func UserEquipmentFromDTO(dto UserEquipment) UserEquipmentEntity {
	return UserEquipmentEntity{
		BaseModel:           dto.BaseModel,
		Owner:               dto.Owner,
		EquipmentTemplateID: dto.EquipmentTemplateID,
		CompendiumVersion:   dto.CompendiumVersion,
	}
}
