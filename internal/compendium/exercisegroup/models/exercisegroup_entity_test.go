package models

import (
	"testing"

	"gesitr/internal/shared"
)

func TestExerciseGroupEntityTableName(t *testing.T) {
	if got := (ExerciseGroupEntity{}).TableName(); got != "exercise_groups" {
		t.Errorf("TableName() = %q, want %q", got, "exercise_groups")
	}
}

func TestExerciseGroupEntityToDTO(t *testing.T) {
	desc := "group desc"
	e := &ExerciseGroupEntity{
		BaseModel:   shared.BaseModel{ID: 10},
		TemplateID:  "g1",
		Name:        "Group One",
		Description: &desc,
		CreatedBy:   "system",
	}
	dto := e.ToDTO()
	if dto.ID != 10 || dto.TemplateID != "g1" || dto.Name != "Group One" || *dto.Description != "group desc" || dto.CreatedBy != "system" {
		t.Error("ToDTO field mismatch")
	}
}

func TestExerciseGroupFromDTO(t *testing.T) {
	desc := "d"
	dto := ExerciseGroup{
		BaseModel:   shared.BaseModel{ID: 7},
		TemplateID:  "g2",
		Name:        "Group Two",
		Description: &desc,
		CreatedBy:   "user",
	}
	e := ExerciseGroupFromDTO(dto)
	if e.ID != 7 || e.TemplateID != "g2" || e.Name != "Group Two" || *e.Description != "d" || e.CreatedBy != "user" {
		t.Error("FromDTO field mismatch")
	}
}

func TestExerciseGroupMemberEntityTableName(t *testing.T) {
	if got := (ExerciseGroupMemberEntity{}).TableName(); got != "exercise_group_members" {
		t.Errorf("TableName() = %q, want %q", got, "exercise_group_members")
	}
}

func TestExerciseGroupMemberEntityToDTO(t *testing.T) {
	e := &ExerciseGroupMemberEntity{
		BaseModel:          shared.BaseModel{ID: 20},
		GroupTemplateID:    "g1",
		ExerciseTemplateID: "ex1",
		AddedBy:            "user",
	}
	dto := e.ToDTO()
	if dto.ID != 20 || dto.GroupTemplateID != "g1" || dto.ExerciseTemplateID != "ex1" || dto.AddedBy != "user" {
		t.Error("ToDTO field mismatch")
	}
}

func TestExerciseGroupMemberFromDTO(t *testing.T) {
	dto := ExerciseGroupMember{
		BaseModel:          shared.BaseModel{ID: 15},
		GroupTemplateID:    "g3",
		ExerciseTemplateID: "ex5",
		AddedBy:            "admin",
	}
	e := ExerciseGroupMemberFromDTO(dto)
	if e.ID != 15 || e.GroupTemplateID != "g3" || e.ExerciseTemplateID != "ex5" || e.AddedBy != "admin" {
		t.Error("FromDTO field mismatch")
	}
}
