package models

type ExerciseGroupEntity struct {
	BaseModel
	TemplateID  string  `gorm:"not null;uniqueIndex"`
	Name        string  `gorm:"not null"`
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
