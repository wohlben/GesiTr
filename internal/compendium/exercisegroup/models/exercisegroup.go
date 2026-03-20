package models

import "gesitr/internal/shared"

type ExerciseGroup struct {
	shared.BaseModel `tstype:",extends"`
	TemplateID       string  `json:"templateId"`
	Name             string  `json:"name"`
	Description      *string `json:"description"`
	CreatedBy        string  `json:"createdBy"`
}

type ExerciseGroupMember struct {
	shared.BaseModel   `tstype:",extends"`
	GroupTemplateID    string `json:"groupTemplateId"`
	ExerciseTemplateID string `json:"exerciseTemplateId"`
	AddedBy            string `json:"addedBy"`
}
