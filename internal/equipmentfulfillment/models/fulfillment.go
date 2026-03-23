package models

import "gesitr/internal/shared"

type Fulfillment struct {
	shared.BaseModel    `tstype:",extends"`
	EquipmentID         uint   `json:"equipmentId"`
	FulfillsEquipmentID uint   `json:"fulfillsEquipmentId"`
	Owner               string `json:"owner"`
}
