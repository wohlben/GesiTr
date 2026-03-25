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

type FulfillmentBody struct {
	EquipmentID         uint `json:"equipmentId" required:"true"`
	FulfillsEquipmentID uint `json:"fulfillsEquipmentId" required:"true"`
}

type CreateFulfillmentInput struct {
	Body FulfillmentBody
}

type CreateFulfillmentOutput struct {
	Body models.Fulfillment
}

type DeleteFulfillmentInput struct {
	ID uint `path:"id"`
}

type DeleteFulfillmentOutput struct{}
