package models

import "gesitr/internal/shared"

type LocalityAvailability struct {
	shared.BaseModel `tstype:",extends"`
	LocalityID       uint   `json:"localityId"`
	EquipmentID      uint   `json:"equipmentId"`
	Available        bool   `json:"available"`
	Owner            string `json:"owner"`
	EquipmentName    string `json:"equipmentName"`
}
