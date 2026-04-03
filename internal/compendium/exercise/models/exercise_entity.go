package models

import (
	"gesitr/internal/shared"
)

// ExerciseEntity is the GORM entity for the exercises table
type ExerciseEntity struct {
	shared.BaseModel
	Type                ExerciseType        `gorm:"not null"`
	TechnicalDifficulty TechnicalDifficulty `gorm:"not null"`
	BodyWeightScaling   float64
	Description         string
	AuthorName          *string
	AuthorUrl           *string
	Owner               string `gorm:"not null;index"`
	Public              bool   `gorm:"not null;default:false;index"`
	Version             int    `gorm:"not null;default:0"`
	ParentExerciseID    *uint

	Forces       []ExerciseForce               `gorm:"foreignKey:ExerciseID"`
	Muscles      []ExerciseMuscle              `gorm:"foreignKey:ExerciseID"`
	Paradigms    []ExerciseMeasurementParadigm `gorm:"foreignKey:ExerciseID"`
	Instructions []ExerciseInstruction         `gorm:"foreignKey:ExerciseID"`
	Images       []ExerciseImage               `gorm:"foreignKey:ExerciseID"`
	Names        []ExerciseName                `gorm:"foreignKey:ExerciseID"`
	Equipment    []ExerciseEquipment           `gorm:"foreignKey:ExerciseID"`
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

type ExerciseName struct {
	ID         uint   `gorm:"primaryKey;autoIncrement"`
	ExerciseID uint   `gorm:"not null;index"`
	Position   int    `gorm:"not null"`
	Name       string `gorm:"not null"`
}

type ExerciseEquipment struct {
	ExerciseID  uint `gorm:"primaryKey"`
	EquipmentID uint `gorm:"primaryKey"`
}

// Conversion functions

func (e *ExerciseEntity) ToDTO() Exercise {
	dto := Exercise{
		BaseModel:           e.BaseModel,
		Type:                e.Type,
		TechnicalDifficulty: e.TechnicalDifficulty,
		BodyWeightScaling:   e.BodyWeightScaling,
		Description:         e.Description,
		AuthorName:          e.AuthorName,
		AuthorUrl:           e.AuthorUrl,
		Owner:               e.Owner,
		Public:              e.Public,
		Version:             e.Version,
		ParentExerciseID:    e.ParentExerciseID,
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
	for _, n := range e.Names {
		dto.Names = append(dto.Names, ExerciseNameDTO{ID: n.ID, Name: n.Name})
	}
	for _, eq := range e.Equipment {
		dto.EquipmentIDs = append(dto.EquipmentIDs, eq.EquipmentID)
	}

	return dto
}

func ExerciseFromDTO(dto Exercise) ExerciseEntity {
	entity := ExerciseEntity{
		BaseModel:           dto.BaseModel,
		Type:                dto.Type,
		TechnicalDifficulty: dto.TechnicalDifficulty,
		BodyWeightScaling:   dto.BodyWeightScaling,
		Description:         dto.Description,
		AuthorName:          dto.AuthorName,
		AuthorUrl:           dto.AuthorUrl,
		Owner:               dto.Owner,
		Public:              dto.Public,
		Version:             dto.Version,
		ParentExerciseID:    dto.ParentExerciseID,
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
	for i, n := range dto.Names {
		entity.Names = append(entity.Names, ExerciseName{Position: i, Name: n.Name})
	}
	for _, eid := range dto.EquipmentIDs {
		entity.Equipment = append(entity.Equipment, ExerciseEquipment{EquipmentID: eid})
	}

	return entity
}
