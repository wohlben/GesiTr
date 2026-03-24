package handlers

import (
	"gesitr/internal/equipmentfulfillment/models"
)

type ListFulfillmentsInput struct {
	EquipmentID         string `query:"equipmentId" doc:"Filter by equipment ID"`
	FulfillsEquipmentID string `query:"fulfillsEquipmentId" doc:"Filter by fulfilled equipment ID"`
}

type ListFulfillmentsOutput struct {
	Body []models.Fulfillment
}

// RawBody skips huma's automatic validation — the Fulfillment DTO is shared
// between request and response and has server-set fields (id, createdAt, etc.)
// that aren't present in create requests.
type CreateFulfillmentInput struct {
	RawBody []byte
}

type CreateFulfillmentOutput struct {
	Body models.Fulfillment
}

type DeleteFulfillmentInput struct {
	ID uint `path:"id"`
}

type DeleteFulfillmentOutput struct{}
