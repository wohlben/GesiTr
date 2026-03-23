package models

import "gesitr/internal/shared"

// Exercise is the API DTO (tygo generates TypeScript from this)
type Exercise struct {
	shared.BaseModel              `tstype:",extends"`
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
	Owner                         string                `json:"owner"`
	Public                        bool                  `json:"public"`
	Version                       int                   `json:"version"`
	ParentExerciseID              *uint                 `json:"parentExerciseId"`
	TemplateID                    string                `json:"templateId"`
	EquipmentIDs                  []uint                `json:"equipmentIds" gorm:"-"`
}

type ExerciseType string

const (
	ExerciseTypeStrength   ExerciseType = "STRENGTH"
	ExerciseTypeCardio     ExerciseType = "CARDIO"
	ExerciseTypeStretching ExerciseType = "STRETCHING"
	ExerciseTypeStrongman  ExerciseType = "STRONGMAN"
)

type Force string

const (
	ForcePull     Force = "PULL"
	ForcePush     Force = "PUSH"
	ForceStatic   Force = "STATIC"
	ForceDynamic  Force = "DYNAMIC"
	ForceHinge    Force = "HINGE"
	ForceRotation Force = "ROTATION"
)

type Muscle string

const (
	MuscleAbs        Muscle = "ABS"
	MuscleAdductors  Muscle = "ADDUCTORS"
	MuscleBiceps     Muscle = "BICEPS"
	MuscleCalves     Muscle = "CALVES"
	MuscleChest      Muscle = "CHEST"
	MuscleForearms   Muscle = "FOREARMS"
	MuscleGlutes     Muscle = "GLUTES"
	MuscleHamstrings Muscle = "HAMSTRINGS"
	MuscleHipFlexors Muscle = "HIP_FLEXORS"
	MuscleLats       Muscle = "LATS"
	MuscleLowerBack  Muscle = "LOWER_BACK"
	MuscleNeck       Muscle = "NECK"
	MuscleObliques   Muscle = "OBLIQUES"
	MuscleQuads      Muscle = "QUADS"
	MuscleTraps      Muscle = "TRAPS"
	MuscleTriceps    Muscle = "TRICEPS"
	MuscleFrontDelts Muscle = "FRONT_DELTS"
	MuscleRearDelts  Muscle = "REAR_DELTS"
	MuscleRhomboids  Muscle = "RHOMBOIDS"
	MuscleSideDelts  Muscle = "SIDE_DELTS"
)

type TechnicalDifficulty string

const (
	DifficultyBeginner     TechnicalDifficulty = "beginner"
	DifficultyIntermediate TechnicalDifficulty = "intermediate"
	DifficultyAdvanced     TechnicalDifficulty = "advanced"
)

type MeasurementParadigm string

const (
	MeasurementRepBased      MeasurementParadigm = "REP_BASED"
	MeasurementAMRAP         MeasurementParadigm = "AMRAP"
	MeasurementTimeBased     MeasurementParadigm = "TIME_BASED"
	MeasurementDistanceBased MeasurementParadigm = "DISTANCE_BASED"
	MeasurementEMOM          MeasurementParadigm = "EMOM"
	MeasurementRoundsForTime MeasurementParadigm = "ROUNDS_FOR_TIME"
	MeasurementTime          MeasurementParadigm = "TIME"
	MeasurementDistance      MeasurementParadigm = "DISTANCE"
)
