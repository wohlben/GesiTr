package models

import "gesitr/internal/shared"

type ExerciseGroupEntity struct {
	shared.BaseModel
	TemplateID  string `gorm:"not null;uniqueIndex"`
	Name        string `gorm:"not null"`
	Description *string
	CreatedBy   string `gorm:"not null"`
}

func (ExerciseGroupEntity) TableName() string { return "exercise_groups" }

func (e *ExerciseGroupEntity) ToDTO() ExerciseGroup {
	return ExerciseGroup{
		BaseModel:   e.BaseModel,
		TemplateID:  e.TemplateID,
		Name:        e.Name,
		Description: e.Description,
		CreatedBy:   e.CreatedBy,
	}
}

func ExerciseGroupFromDTO(dto ExerciseGroup) ExerciseGroupEntity {
	return ExerciseGroupEntity{
		BaseModel:   dto.BaseModel,
		TemplateID:  dto.TemplateID,
		Name:        dto.Name,
		Description: dto.Description,
		CreatedBy:   dto.CreatedBy,
	}
}

type ExerciseGroupMemberEntity struct {
	shared.BaseModel
	GroupTemplateID    string `gorm:"not null;uniqueIndex:idx_group_member"`
	ExerciseTemplateID string `gorm:"not null;uniqueIndex:idx_group_member"`
	AddedBy            string `gorm:"not null"`
}

func (ExerciseGroupMemberEntity) TableName() string { return "exercise_group_members" }

func (e *ExerciseGroupMemberEntity) ToDTO() ExerciseGroupMember {
	return ExerciseGroupMember{
		BaseModel:          e.BaseModel,
		GroupTemplateID:    e.GroupTemplateID,
		ExerciseTemplateID: e.ExerciseTemplateID,
		AddedBy:            e.AddedBy,
	}
}

func ExerciseGroupMemberFromDTO(dto ExerciseGroupMember) ExerciseGroupMemberEntity {
	return ExerciseGroupMemberEntity{
		BaseModel:          dto.BaseModel,
		GroupTemplateID:    dto.GroupTemplateID,
		ExerciseTemplateID: dto.ExerciseTemplateID,
		AddedBy:            dto.AddedBy,
	}
}
