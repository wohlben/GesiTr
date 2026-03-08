package models

type UserEquipment struct {
	BaseModel             `tstype:",extends"`
	Owner                 string `json:"owner"`
	CompendiumEquipmentID string `json:"compendiumEquipmentId"`
	CompendiumVersion     int    `json:"compendiumVersion"`
}
