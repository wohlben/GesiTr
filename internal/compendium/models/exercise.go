package models

// Exercise is the API DTO (tygo generates TypeScript from this)
type Exercise struct {
	BaseModel                     `tstype:",extends"`
	Name                          string                `json:"name"`
	Slug                          string                `json:"slug"`
	Type                          ExerciseType          `json:"type"`
	Force                         []Force               `json:"force" gorm:"-"`
	PrimaryMuscles                []Muscle              `json:"primaryMuscles" gorm:"-"`
	SecondaryMuscles              []Muscle              `json:"secondaryMuscles" gorm:"-"`
	TechnicalDifficulty           TechnicalDifficulty   `json:"technicalDifficulty"`
	BodyWeightScaling             float64               `json:"bodyWeightScaling"`
	SuggestedMeasurementParadigms []MeasurementParadigm `json:"suggestedMeasurementParadigms" gorm:"-"`
	Description                   string                `json:"description"`
	Instructions                  []string              `json:"instructions" gorm:"-"`
	Images                        []string              `json:"images" gorm:"-"`
	AlternativeNames              []string              `json:"alternativeNames" gorm:"-"`
	AuthorName                    *string               `json:"authorName"`
	AuthorUrl                     *string               `json:"authorUrl"`
	CreatedBy                     string                `json:"createdBy"`
	Version                       int                   `json:"version"`
	ParentExerciseID              *uint                 `json:"parentExerciseId"`
	TemplateID                    *string               `json:"templateId"`
	EquipmentIDs                  []string              `json:"equipmentIds" gorm:"-"`
}
