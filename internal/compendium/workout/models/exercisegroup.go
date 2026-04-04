package models

import "gesitr/internal/shared"

type ExerciseGroup struct {
	shared.BaseModel `tstype:",extends"`
	Name             *string `json:"name"`
	OwnershipGroupID uint    `json:"ownershipGroupId"`
}

type ExerciseGroupMember struct {
	shared.BaseModel `tstype:",extends"`
	GroupID          uint `json:"groupId"`
	ExerciseID       uint `json:"exerciseId"`
	OwnershipGroupID uint `json:"ownershipGroupId"`
}
