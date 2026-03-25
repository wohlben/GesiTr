package models

import "gesitr/internal/shared"

type ExerciseGroup struct {
	shared.BaseModel `tstype:",extends"`
	Name             *string `json:"name"`
	Owner            string  `json:"owner"`
}

type ExerciseGroupMember struct {
	shared.BaseModel `tstype:",extends"`
	GroupID          uint   `json:"groupId"`
	ExerciseID       uint   `json:"exerciseId"`
	Owner            string `json:"owner"`
}
