package models

import (
	"gesitr/internal/shared"
)

type LocalityEntity struct {
	shared.BaseModel
	Name             string `gorm:"not null"`
	OwnershipGroupID uint   `gorm:"index"`
	Public           bool   `gorm:"not null;default:false;index"`
}

func (LocalityEntity) TableName() string { return "localities" }

func (e *LocalityEntity) ToDTO() Locality {
	return Locality{
		BaseModel:        e.BaseModel,
		Name:             e.Name,
		OwnershipGroupID: e.OwnershipGroupID,
		Public:           e.Public,
	}
}

func LocalityFromDTO(dto Locality) LocalityEntity {
	return LocalityEntity{
		BaseModel:        dto.BaseModel,
		Name:             dto.Name,
		OwnershipGroupID: dto.OwnershipGroupID,
		Public:           dto.Public,
	}
}
