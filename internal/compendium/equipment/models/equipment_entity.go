package models

import (
	"gesitr/internal/shared"
)

type EquipmentEntity struct {
	shared.BaseModel
	Name        string `gorm:"not null"`
	DisplayName string `gorm:"not null"`
	Description string
	Category    EquipmentCategory `gorm:"not null"`
	ImageUrl    *string
	Owner       string `gorm:"not null;index"`
	Public      bool   `gorm:"not null;default:false;index"`
	Version     int    `gorm:"not null;default:0"`
}

func (EquipmentEntity) TableName() string { return "equipment" }

func (e *EquipmentEntity) ToDTO() Equipment {
	return Equipment{
		BaseModel:   e.BaseModel,
		Name:        e.Name,
		DisplayName: e.DisplayName,
		Description: e.Description,
		Category:    e.Category,
		ImageUrl:    e.ImageUrl,
		Owner:       e.Owner,
		Public:      e.Public,
		Version:     e.Version,
	}
}

func EquipmentFromDTO(dto Equipment) EquipmentEntity {
	return EquipmentEntity{
		BaseModel:   dto.BaseModel,
		Name:        dto.Name,
		DisplayName: dto.DisplayName,
		Description: dto.Description,
		Category:    dto.Category,
		ImageUrl:    dto.ImageUrl,
		Owner:       dto.Owner,
		Public:      dto.Public,
		Version:     dto.Version,
	}
}
