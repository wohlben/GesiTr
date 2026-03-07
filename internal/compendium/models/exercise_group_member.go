package models

type ExerciseGroupMember struct {
	BaseModel              `tstype:",extends"`
	GroupTemplateID        string `json:"groupTemplateId"`
	ExerciseTemplateID     string `json:"exerciseTemplateId"`
	AddedBy                string `json:"addedBy"`
}
