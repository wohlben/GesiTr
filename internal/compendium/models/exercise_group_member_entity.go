package models

type ExerciseGroupMemberEntity struct {
	BaseModel
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
