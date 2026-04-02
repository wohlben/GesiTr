package handlers

import (
	"gesitr/internal/user/mastery/models"
)

type ListEquipmentMasteryInput struct{}

type ListEquipmentMasteryOutput struct {
	Body []models.EquipmentMastery
}

type GetEquipmentMasteryInput struct {
	EquipmentID uint `path:"equipmentId"`
}

type GetEquipmentMasteryOutput struct {
	Body models.EquipmentMastery
}
