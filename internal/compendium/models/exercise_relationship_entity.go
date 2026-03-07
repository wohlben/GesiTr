package models

type ExerciseRelationshipEntity struct {
	BaseModel
	RelationshipType       ExerciseRelationshipType `gorm:"not null;uniqueIndex:idx_exercise_relationship"`
	Strength               float64 `gorm:"not null"`
	Description            *string
	CreatedBy              string `gorm:"not null"`
	FromExerciseTemplateID string `gorm:"not null;uniqueIndex:idx_exercise_relationship"`
	ToExerciseTemplateID   string `gorm:"not null;uniqueIndex:idx_exercise_relationship"`
}

func (ExerciseRelationshipEntity) TableName() string { return "exercise_relationships" }

func (e *ExerciseRelationshipEntity) ToDTO() ExerciseRelationship {
	return ExerciseRelationship{
		BaseModel:              e.BaseModel,
		RelationshipType:       e.RelationshipType,
		Strength:               e.Strength,
		Description:            e.Description,
		CreatedBy:              e.CreatedBy,
		FromExerciseTemplateID: e.FromExerciseTemplateID,
		ToExerciseTemplateID:   e.ToExerciseTemplateID,
	}
}

func ExerciseRelationshipFromDTO(dto ExerciseRelationship) ExerciseRelationshipEntity {
	return ExerciseRelationshipEntity{
		BaseModel:              dto.BaseModel,
		RelationshipType:       dto.RelationshipType,
		Strength:               dto.Strength,
		Description:            dto.Description,
		CreatedBy:              dto.CreatedBy,
		FromExerciseTemplateID: dto.FromExerciseTemplateID,
		ToExerciseTemplateID:   dto.ToExerciseTemplateID,
	}
}
