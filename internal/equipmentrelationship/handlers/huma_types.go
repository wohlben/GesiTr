package handlers

import (
	"gesitr/internal/equipmentrelationship/models"
)

type ListEquipmentRelationshipsInput struct {
	Owner            string `query:"owner" doc:"Filter by owner"`
	FromEquipmentID  string `query:"fromEquipmentId" doc:"Filter by source equipment ID"`
	ToEquipmentID    string `query:"toEquipmentId" doc:"Filter by target equipment ID"`
	RelationshipType string `query:"relationshipType" doc:"Filter by relationship type"`
}

type ListEquipmentRelationshipsOutput struct {
	Body []models.EquipmentRelationship
}

type EquipmentRelationshipBody struct {
	RelationshipType models.EquipmentRelationshipType `json:"relationshipType" required:"true"`
	Strength         float64                          `json:"strength" required:"false"`
	FromEquipmentID  uint                             `json:"fromEquipmentId" required:"true"`
	ToEquipmentID    uint                             `json:"toEquipmentId" required:"true"`
}

type CreateEquipmentRelationshipInput struct {
	Body EquipmentRelationshipBody
}

type CreateEquipmentRelationshipOutput struct {
	Body models.EquipmentRelationship
}

type DeleteEquipmentRelationshipInput struct {
	ID uint `path:"id"`
}

type DeleteEquipmentRelationshipOutput struct{}
