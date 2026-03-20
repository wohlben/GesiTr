package models

import "gesitr/internal/shared"

type UserEquipment struct {
	shared.BaseModel      `tstype:",extends"`
	Owner                 string `json:"owner"`
	CompendiumEquipmentID string `json:"compendiumEquipmentId"`
	CompendiumVersion     int    `json:"compendiumVersion"`
}
