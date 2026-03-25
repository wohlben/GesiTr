package models

import (
	profilemodels "gesitr/internal/profile/models"
	"gesitr/internal/shared"
)

type ExerciseGroupEntity struct {
	shared.BaseModel
	Name         string `gorm:"not null"`
	Description  *string
	Owner        string                           `gorm:"not null"`
	OwnerProfile *profilemodels.UserProfileEntity `gorm:"foreignKey:Owner;references:ID;constraint:OnDelete:RESTRICT" json:"-"`
}

func (ExerciseGroupEntity) TableName() string { return "exercise_groups" }

func (e *ExerciseGroupEntity) ToDTO() ExerciseGroup {
	return ExerciseGroup{
		BaseModel:   e.BaseModel,
		Name:        e.Name,
		Description: e.Description,
		Owner:       e.Owner,
	}
}

func ExerciseGroupFromDTO(dto ExerciseGroup) ExerciseGroupEntity {
	return ExerciseGroupEntity{
		BaseModel:   dto.BaseModel,
		Name:        dto.Name,
		Description: dto.Description,
		Owner:       dto.Owner,
	}
}

type ExerciseGroupMemberEntity struct {
	shared.BaseModel
	GroupID      uint                             `gorm:"not null;uniqueIndex:idx_group_member"`
	ExerciseID   uint                             `gorm:"not null;uniqueIndex:idx_group_member"`
	Owner        string                           `gorm:"not null"`
	OwnerProfile *profilemodels.UserProfileEntity `gorm:"foreignKey:Owner;references:ID;constraint:OnDelete:RESTRICT" json:"-"`
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
