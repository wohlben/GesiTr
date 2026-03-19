package equipmentfulfillment

import "gesitr/internal/shared"

type Fulfillment struct {
	shared.BaseModel            `tstype:",extends"`
	EquipmentTemplateID         string `json:"equipmentTemplateId"`
	FulfillsEquipmentTemplateID string `json:"fulfillsEquipmentTemplateId"`
	CreatedBy                   string `json:"createdBy"`
}
