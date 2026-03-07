package models

type Fulfillment struct {
	BaseModel                   `tstype:",extends"`
	EquipmentTemplateID         string `json:"equipmentTemplateId"`
	FulfillsEquipmentTemplateID string `json:"fulfillsEquipmentTemplateId"`
	CreatedBy                   string `json:"createdBy"`
}
