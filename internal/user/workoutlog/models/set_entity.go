package models

import (
	"time"

	exerciselogmodels "gesitr/internal/user/exerciselog/models"

	"gesitr/internal/shared"
)

type WorkoutLogExerciseSetEntity struct {
	shared.BaseModel
	WorkoutLogExerciseID uint                 `gorm:"not null;index"`
	WorkoutLogID         uint                 `gorm:"not null;index"`
	SetNumber            int                  `gorm:"not null"`
	Status               WorkoutLogItemStatus `gorm:"not null;default:'planning'"`
	StatusChangedAt      *time.Time
	BreakAfterSeconds    *int
	TargetReps           *int
	TargetWeight         *float64
	TargetDuration       *int
	TargetDistance       *float64
	TargetTime           *int
	ExerciseLog          *exerciselogmodels.ExerciseLogEntity `gorm:"foreignKey:WorkoutLogExerciseSetID"`
}

func (WorkoutLogExerciseSetEntity) TableName() string { return "workout_log_exercise_sets" }

func (e *WorkoutLogExerciseSetEntity) ToDTO() WorkoutLogExerciseSet {
	dto := WorkoutLogExerciseSet{
		BaseModel:            e.BaseModel,
		WorkoutLogExerciseID: e.WorkoutLogExerciseID,
		WorkoutLogID:         e.WorkoutLogID,
		SetNumber:            e.SetNumber,
		Status:               e.Status,
		StatusChangedAt:      e.StatusChangedAt,
		BreakAfterSeconds:    e.BreakAfterSeconds,
		TargetReps:           e.TargetReps,
		TargetWeight:         e.TargetWeight,
		TargetDuration:       e.TargetDuration,
		TargetDistance:       e.TargetDistance,
		TargetTime:           e.TargetTime,
	}
	if e.ExerciseLog != nil {
		el := e.ExerciseLog.ToDTO()
		dto.ExerciseLog = &el
	}
	return dto
}

func WorkoutLogExerciseSetFromDTO(dto WorkoutLogExerciseSet) WorkoutLogExerciseSetEntity {
	return WorkoutLogExerciseSetEntity{
		BaseModel:            dto.BaseModel,
		WorkoutLogExerciseID: dto.WorkoutLogExerciseID,
		WorkoutLogID:         dto.WorkoutLogID,
		SetNumber:            dto.SetNumber,
		Status:               dto.Status,
		StatusChangedAt:      dto.StatusChangedAt,
		BreakAfterSeconds:    dto.BreakAfterSeconds,
		TargetReps:           dto.TargetReps,
		TargetWeight:         dto.TargetWeight,
		TargetDuration:       dto.TargetDuration,
		TargetDistance:       dto.TargetDistance,
		TargetTime:           dto.TargetTime,
	}
}
