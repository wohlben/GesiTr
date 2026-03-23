package models

import "gesitr/internal/shared"

type ExerciseGroup struct {
	shared.BaseModel `tstype:",extends"`
	TemplateID       string  `json:"templateId"`
	Name             string  `json:"name"`
	Description      *string `json:"description"`
	Owner            string  `json:"owner"`
}

type ExerciseGroupMember struct {
	shared.BaseModel `tstype:",extends"`
	GroupID          uint   `json:"groupId"`
	ExerciseID       uint   `json:"exerciseId"`
	Owner            string `json:"owner"`
}
