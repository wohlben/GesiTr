package models

import "gesitr/internal/shared"

type EquipmentRelationship struct {
	shared.BaseModel `tstype:",extends"`
	RelationshipType EquipmentRelationshipType `json:"relationshipType"`
	Strength         float64                   `json:"strength"`
	Owner            string                    `json:"owner"`
	FromEquipmentID  uint                      `json:"fromEquipmentId"`
	ToEquipmentID    uint                      `json:"toEquipmentId"`
}

type EquipmentRelationshipType string

const (
	EquipmentRelationshipTypeEquivalent EquipmentRelationshipType = "equivalent"
)
