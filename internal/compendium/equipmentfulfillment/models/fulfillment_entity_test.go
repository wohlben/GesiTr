package models

import (
	"testing"

	"gesitr/internal/shared"
)

func TestFulfillmentEntityTableName(t *testing.T) {
	if got := (FulfillmentEntity{}).TableName(); got != "fulfillments" {
		t.Errorf("TableName() = %q, want %q", got, "fulfillments")
	}
}

func TestFulfillmentEntityToDTO(t *testing.T) {
	e := &FulfillmentEntity{
		BaseModel:                   shared.BaseModel{ID: 1},
		EquipmentTemplateID:         "eq1",
		FulfillsEquipmentTemplateID: "eq2",
		CreatedBy:                   "system",
	}
	dto := e.ToDTO()
	if dto.ID != 1 || dto.EquipmentTemplateID != "eq1" || dto.FulfillsEquipmentTemplateID != "eq2" || dto.CreatedBy != "system" {
		t.Error("ToDTO field mismatch")
	}
}

func TestFulfillmentFromDTO(t *testing.T) {
	dto := Fulfillment{
		BaseModel:                   shared.BaseModel{ID: 2},
		EquipmentTemplateID:         "a",
		FulfillsEquipmentTemplateID: "b",
		CreatedBy:                   "user",
	}
	e := FulfillmentFromDTO(dto)
	if e.ID != 2 || e.EquipmentTemplateID != "a" || e.FulfillsEquipmentTemplateID != "b" || e.CreatedBy != "user" {
		t.Error("FulfillmentFromDTO field mismatch")
	}
}
