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

// RawBody skips huma's automatic validation — the EquipmentRelationship DTO is shared
// between request and response and has server-set fields (id, createdAt, etc.)
// that aren't present in create requests.
type CreateEquipmentRelationshipInput struct {
	RawBody []byte
}

type CreateEquipmentRelationshipOutput struct {
	Body models.EquipmentRelationship
}

type DeleteEquipmentRelationshipInput struct {
	ID uint `path:"id"`
}

type DeleteEquipmentRelationshipOutput struct{}
