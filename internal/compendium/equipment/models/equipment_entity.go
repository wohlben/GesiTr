package models

import (
	profilemodels "gesitr/internal/profile/models"
	"gesitr/internal/shared"
)

type EquipmentEntity struct {
	shared.BaseModel
	Name             string `gorm:"not null"`
	DisplayName      string `gorm:"not null"`
	Description      string
	Category         EquipmentCategory `gorm:"not null"`
	ImageUrl         *string
	TemplateID       string                           `gorm:"not null;uniqueIndex"`
	CreatedBy        string                           `gorm:"not null"`
	CreatedByProfile *profilemodels.UserProfileEntity `gorm:"foreignKey:CreatedBy;references:ID;constraint:OnDelete:RESTRICT" json:"-"`
	Version          int                              `gorm:"not null;default:0"`
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
		TemplateID:  e.TemplateID,
		CreatedBy:   e.CreatedBy,
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
		TemplateID:  dto.TemplateID,
		CreatedBy:   dto.CreatedBy,
		Version:     dto.Version,
	}
}
