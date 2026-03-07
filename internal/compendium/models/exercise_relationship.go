package models

type ExerciseRelationship struct {
	BaseModel                `tstype:",extends"`
	RelationshipType         ExerciseRelationshipType `json:"relationshipType"`
	Strength                 float64 `json:"strength"`
	Description              *string `json:"description"`
	CreatedBy                string  `json:"createdBy"`
	FromExerciseTemplateID   string  `json:"fromExerciseTemplateId"`
	ToExerciseTemplateID     string  `json:"toExerciseTemplateId"`
}
