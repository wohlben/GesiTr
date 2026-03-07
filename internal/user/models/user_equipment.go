package models

type UserEquipment struct {
	BaseModel           `tstype:",extends"`
	Owner               string `json:"owner"`
	EquipmentTemplateID string `json:"equipmentTemplateId"`
	CompendiumVersion   int    `json:"compendiumVersion"`
}
