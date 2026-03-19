package exercise

import "gesitr/internal/shared"

// ExerciseEntity is the GORM entity for the exercises table
type ExerciseEntity struct {
	shared.BaseModel
	Name                string              `gorm:"not null"`
	Slug                string              `gorm:"not null;uniqueIndex"`
	Type                ExerciseType        `gorm:"not null"`
	TechnicalDifficulty TechnicalDifficulty `gorm:"not null"`
	BodyWeightScaling   float64
	Description         string
	AuthorName          *string
	AuthorUrl           *string
	CreatedBy           string `gorm:"not null"`
	Version             int    `gorm:"not null;default:0"`
	ParentExerciseID    *uint
	TemplateID          string `gorm:"not null;uniqueIndex"`

	Forces           []ExerciseForce               `gorm:"foreignKey:ExerciseID"`
	Muscles          []ExerciseMuscle              `gorm:"foreignKey:ExerciseID"`
	Paradigms        []ExerciseMeasurementParadigm `gorm:"foreignKey:ExerciseID"`
	Instructions     []ExerciseInstruction         `gorm:"foreignKey:ExerciseID"`
	Images           []ExerciseImage               `gorm:"foreignKey:ExerciseID"`
	AlternativeNames []ExerciseAlternativeName     `gorm:"foreignKey:ExerciseID"`
	Equipment        []ExerciseEquipment           `gorm:"foreignKey:ExerciseID"`
}

func (ExerciseEntity) TableName() string { return "exercises" }

// Child tables

type ExerciseForce struct {
	ExerciseID uint  `gorm:"primaryKey"`
	Force      Force `gorm:"primaryKey"`
}

type ExerciseMuscle struct {
	ExerciseID uint   `gorm:"primaryKey"`
	Muscle     Muscle `gorm:"primaryKey"`
	IsPrimary  bool   `gorm:"primaryKey"`
}

type ExerciseMeasurementParadigm struct {
	ExerciseID uint                `gorm:"primaryKey"`
	Paradigm   MeasurementParadigm `gorm:"primaryKey"`
}

type ExerciseInstruction struct {
	ExerciseID uint   `gorm:"primaryKey"`
	Position   int    `gorm:"primaryKey"`
	Text       string `gorm:"not null"`
}

type ExerciseImage struct {
	ExerciseID uint   `gorm:"primaryKey"`
	Position   int    `gorm:"primaryKey"`
	Path       string `gorm:"not null"`
}

type ExerciseAlternativeName struct {
	ExerciseID uint   `gorm:"primaryKey"`
	Name       string `gorm:"primaryKey"`
}

type ExerciseEquipment struct {
	ExerciseID          uint   `gorm:"primaryKey"`
	EquipmentTemplateID string `gorm:"primaryKey"`
}

// Conversion functions

func (e *ExerciseEntity) ToDTO() Exercise {
	dto := Exercise{
		BaseModel:           e.BaseModel,
		Name:                e.Name,
		Slug:                e.Slug,
		Type:                e.Type,
		TechnicalDifficulty: e.TechnicalDifficulty,
		BodyWeightScaling:   e.BodyWeightScaling,
		Description:         e.Description,
		AuthorName:          e.AuthorName,
		AuthorUrl:           e.AuthorUrl,
		CreatedBy:           e.CreatedBy,
		Version:             e.Version,
		ParentExerciseID:    e.ParentExerciseID,
		TemplateID:          e.TemplateID,
	}

	for _, f := range e.Forces {
		dto.Force = append(dto.Force, f.Force)
	}
	for _, m := range e.Muscles {
		if m.IsPrimary {
			dto.PrimaryMuscles = append(dto.PrimaryMuscles, m.Muscle)
		} else {
			dto.SecondaryMuscles = append(dto.SecondaryMuscles, m.Muscle)
		}
	}
	for _, p := range e.Paradigms {
		dto.SuggestedMeasurementParadigms = append(dto.SuggestedMeasurementParadigms, p.Paradigm)
	}
	for _, i := range e.Instructions {
		dto.Instructions = append(dto.Instructions, i.Text)
	}
	for _, img := range e.Images {
		dto.Images = append(dto.Images, img.Path)
	}
	for _, an := range e.AlternativeNames {
		dto.AlternativeNames = append(dto.AlternativeNames, an.Name)
	}
	for _, eq := range e.Equipment {
		dto.EquipmentIDs = append(dto.EquipmentIDs, eq.EquipmentTemplateID)
	}

	return dto
}

func ExerciseFromDTO(dto Exercise) ExerciseEntity {
	entity := ExerciseEntity{
		BaseModel:           dto.BaseModel,
		Name:                dto.Name,
		Slug:                dto.Slug,
		Type:                dto.Type,
		TechnicalDifficulty: dto.TechnicalDifficulty,
		BodyWeightScaling:   dto.BodyWeightScaling,
		Description:         dto.Description,
		AuthorName:          dto.AuthorName,
		AuthorUrl:           dto.AuthorUrl,
		CreatedBy:           dto.CreatedBy,
		Version:             dto.Version,
		ParentExerciseID:    dto.ParentExerciseID,
		TemplateID:          dto.TemplateID,
	}

	for _, f := range dto.Force {
		entity.Forces = append(entity.Forces, ExerciseForce{Force: f})
	}
	for _, m := range dto.PrimaryMuscles {
		entity.Muscles = append(entity.Muscles, ExerciseMuscle{Muscle: m, IsPrimary: true})
	}
	for _, m := range dto.SecondaryMuscles {
		entity.Muscles = append(entity.Muscles, ExerciseMuscle{Muscle: m, IsPrimary: false})
	}
	for _, p := range dto.SuggestedMeasurementParadigms {
		entity.Paradigms = append(entity.Paradigms, ExerciseMeasurementParadigm{Paradigm: p})
	}
	for i, text := range dto.Instructions {
		entity.Instructions = append(entity.Instructions, ExerciseInstruction{Position: i, Text: text})
	}
	for i, path := range dto.Images {
		entity.Images = append(entity.Images, ExerciseImage{Position: i, Path: path})
	}
	for _, name := range dto.AlternativeNames {
		entity.AlternativeNames = append(entity.AlternativeNames, ExerciseAlternativeName{Name: name})
	}
	for _, tid := range dto.EquipmentIDs {
		entity.Equipment = append(entity.Equipment, ExerciseEquipment{EquipmentTemplateID: tid})
	}

	return entity
}
