package models

import (
	"gesitr/internal/shared"
)

type ExerciseGroupEntity struct {
	shared.BaseModel
	Name  *string
	Owner string `gorm:"not null"`
}

func (ExerciseGroupEntity) TableName() string { return "exercise_groups" }

func (e *ExerciseGroupEntity) ToDTO() ExerciseGroup {
	return ExerciseGroup{
		BaseModel: e.BaseModel,
		Name:      e.Name,
		Owner:     e.Owner,
	}
}

func ExerciseGroupFromDTO(dto ExerciseGroup) ExerciseGroupEntity {
	return ExerciseGroupEntity{
		BaseModel: dto.BaseModel,
		Name:      dto.Name,
		Owner:     dto.Owner,
	}
}

type ExerciseGroupMemberEntity struct {
	shared.BaseModel
	GroupID    uint   `gorm:"not null;uniqueIndex:idx_group_member"`
	ExerciseID uint   `gorm:"not null;uniqueIndex:idx_group_member"`
	Owner      string `gorm:"not null"`
}

func (ExerciseGroupMemberEntity) TableName() string { return "exercise_group_members" }

func (e *ExerciseGroupMemberEntity) ToDTO() ExerciseGroupMember {
	return ExerciseGroupMember{
		BaseModel:  e.BaseModel,
		GroupID:    e.GroupID,
		ExerciseID: e.ExerciseID,
		Owner:      e.Owner,
	}
}

func ExerciseGroupMemberFromDTO(dto ExerciseGroupMember) ExerciseGroupMemberEntity {
	return ExerciseGroupMemberEntity{
		BaseModel:  dto.BaseModel,
		GroupID:    dto.GroupID,
		ExerciseID: dto.ExerciseID,
		Owner:      dto.Owner,
	}
}
