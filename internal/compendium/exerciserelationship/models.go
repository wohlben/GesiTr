package exerciserelationship

import "gesitr/internal/shared"

type ExerciseRelationship struct {
	shared.BaseModel       `tstype:",extends"`
	RelationshipType       ExerciseRelationshipType `json:"relationshipType"`
	Strength               float64                  `json:"strength"`
	Description            *string                  `json:"description"`
	CreatedBy              string                   `json:"createdBy"`
	FromExerciseTemplateID string                   `json:"fromExerciseTemplateId"`
	ToExerciseTemplateID   string                   `json:"toExerciseTemplateId"`
}

type ExerciseRelationshipType string

const (
	ExerciseRelationshipTypeAccessory           ExerciseRelationshipType = "accessory"
	ExerciseRelationshipTypeAlternative         ExerciseRelationshipType = "alternative"
	ExerciseRelationshipTypeAntagonist          ExerciseRelationshipType = "antagonist"
	ExerciseRelationshipTypeBilateralUnilateral ExerciseRelationshipType = "bilateral_unilateral"
	ExerciseRelationshipTypeComplementary       ExerciseRelationshipType = "complementary"
	ExerciseRelationshipTypeEasierAlternative   ExerciseRelationshipType = "easier_alternative"
	ExerciseRelationshipTypeEquipmentVariation  ExerciseRelationshipType = "equipment_variation"
	ExerciseRelationshipTypeHarderAlternative   ExerciseRelationshipType = "harder_alternative"
	ExerciseRelationshipTypePreparation         ExerciseRelationshipType = "preparation"
	ExerciseRelationshipTypePrerequisite        ExerciseRelationshipType = "prerequisite"
	ExerciseRelationshipTypeProgressesTo        ExerciseRelationshipType = "progresses_to"
	ExerciseRelationshipTypeProgression         ExerciseRelationshipType = "progression"
	ExerciseRelationshipTypeRegressesTo         ExerciseRelationshipType = "regresses_to"
	ExerciseRelationshipTypeRegression          ExerciseRelationshipType = "regression"
	ExerciseRelationshipTypeRelated             ExerciseRelationshipType = "related"
	ExerciseRelationshipTypeSimilar             ExerciseRelationshipType = "similar"
	ExerciseRelationshipTypeSupersetWith        ExerciseRelationshipType = "superset_with"
	ExerciseRelationshipTypeSupports            ExerciseRelationshipType = "supports"
	ExerciseRelationshipTypeVariant             ExerciseRelationshipType = "variant"
	ExerciseRelationshipTypeVariation           ExerciseRelationshipType = "variation"
)
