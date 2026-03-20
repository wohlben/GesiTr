package models

import (
	"testing"
	"time"

	"gesitr/internal/shared"
)

func TestEquipmentEntityTableName(t *testing.T) {
	if got := (EquipmentEntity{}).TableName(); got != "equipment" {
		t.Errorf("TableName() = %q, want %q", got, "equipment")
	}
}

func TestEquipmentEntityToDTO(t *testing.T) {
	now := time.Now()
	imgUrl := "http://example.com/img.png"
	e := &EquipmentEntity{
		BaseModel:   shared.BaseModel{ID: 1, CreatedAt: now, UpdatedAt: now},
		Name:        "barbell",
		DisplayName: "Barbell",
		Description: "A long bar",
		Category:    EquipmentCategoryFreeWeights,
		ImageUrl:    &imgUrl,
		TemplateID:  "barbell",
		CreatedBy:   "system",
		Version:     2,
	}
	dto := e.ToDTO()
	if dto.ID != 1 {
		t.Errorf("ID = %d, want 1", dto.ID)
	}
	if dto.Name != "barbell" {
		t.Errorf("Name = %q", dto.Name)
	}
	if dto.DisplayName != "Barbell" {
		t.Errorf("DisplayName = %q", dto.DisplayName)
	}
	if dto.Description != "A long bar" {
		t.Errorf("Description = %q", dto.Description)
	}
	if dto.Category != EquipmentCategoryFreeWeights {
		t.Errorf("Category = %q", dto.Category)
	}
	if *dto.ImageUrl != imgUrl {
		t.Errorf("ImageUrl = %q", *dto.ImageUrl)
	}
	if dto.TemplateID != "barbell" {
		t.Errorf("TemplateID = %q", dto.TemplateID)
	}
	if dto.CreatedBy != "system" {
		t.Errorf("CreatedBy = %q", dto.CreatedBy)
	}
	if dto.Version != 2 {
		t.Errorf("Version = %d", dto.Version)
	}
}

func TestEquipmentFromDTO(t *testing.T) {
	imgUrl := "http://example.com/img.png"
	dto := Equipment{
		BaseModel:   shared.BaseModel{ID: 3},
		Name:        "bench",
		DisplayName: "Bench",
		Description: "A flat bench",
		Category:    EquipmentCategoryBenches,
		ImageUrl:    &imgUrl,
		TemplateID:  "bench",
		CreatedBy:   "user",
		Version:     1,
	}
	e := EquipmentFromDTO(dto)
	if e.ID != 3 || e.Name != "bench" || e.Category != EquipmentCategoryBenches || e.Version != 1 {
		t.Error("EquipmentFromDTO field mismatch")
	}
	if e.DisplayName != "Bench" || e.Description != "A flat bench" || e.TemplateID != "bench" || e.CreatedBy != "user" {
		t.Error("EquipmentFromDTO field mismatch")
	}
	if *e.ImageUrl != imgUrl {
		t.Error("EquipmentFromDTO ImageUrl mismatch")
	}
}
