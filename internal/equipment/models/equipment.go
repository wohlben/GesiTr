package models

import "gesitr/internal/shared"

type Equipment struct {
	shared.BaseModel `tstype:",extends"`
	Name             string            `json:"name"`
	DisplayName      string            `json:"displayName"`
	Description      string            `json:"description"`
	Category         EquipmentCategory `json:"category"`
	ImageUrl         *string           `json:"imageUrl"`
	TemplateID       string            `json:"templateId"`
	Owner            string            `json:"owner"`
	Public           bool              `json:"public"`
	Version          int               `json:"version"`
}

type EquipmentCategory string

const (
	EquipmentCategoryFreeWeights EquipmentCategory = "free_weights"
	EquipmentCategoryAccessories EquipmentCategory = "accessories"
	EquipmentCategoryBenches     EquipmentCategory = "benches"
	EquipmentCategoryMachines    EquipmentCategory = "machines"
	EquipmentCategoryFunctional  EquipmentCategory = "functional"
	EquipmentCategoryOther       EquipmentCategory = "other"
)
