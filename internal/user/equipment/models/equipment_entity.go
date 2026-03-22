package models

import (
	profilemodels "gesitr/internal/profile/models"
	"gesitr/internal/shared"
)

type UserEquipmentEntity struct {
	shared.BaseModel
	Owner                 string                           `gorm:"not null;index;uniqueIndex:idx_owner_compendium_equipment"`
	OwnerProfile          *profilemodels.UserProfileEntity `gorm:"foreignKey:Owner;references:ID;constraint:OnDelete:RESTRICT" json:"-"`
	CompendiumEquipmentID string                           `gorm:"not null;uniqueIndex:idx_owner_compendium_equipment"`
	CompendiumVersion     int                              `gorm:"not null"`
}

func (UserEquipmentEntity) TableName() string { return "user_equipment" }

func (e *UserEquipmentEntity) ToDTO() UserEquipment {
	return UserEquipment{
		BaseModel:             e.BaseModel,
		Owner:                 e.Owner,
		CompendiumEquipmentID: e.CompendiumEquipmentID,
		CompendiumVersion:     e.CompendiumVersion,
	}
}

func UserEquipmentFromDTO(dto UserEquipment) UserEquipmentEntity {
	return UserEquipmentEntity{
		BaseModel:             dto.BaseModel,
		Owner:                 dto.Owner,
		CompendiumEquipmentID: dto.CompendiumEquipmentID,
		CompendiumVersion:     dto.CompendiumVersion,
	}
}
