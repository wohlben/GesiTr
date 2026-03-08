package models

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

type EquipmentCategory string

const (
	EquipmentCategoryFreeWeights EquipmentCategory = "free_weights"
	EquipmentCategoryAccessories EquipmentCategory = "accessories"
	EquipmentCategoryBenches     EquipmentCategory = "benches"
	EquipmentCategoryMachines    EquipmentCategory = "machines"
	EquipmentCategoryFunctional  EquipmentCategory = "functional"
	EquipmentCategoryOther       EquipmentCategory = "other"
)

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
