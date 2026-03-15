package models

import "time"

type WorkoutLogExerciseSetEntity struct {
	BaseModel
	WorkoutLogExerciseID uint             `gorm:"not null;index"`
	WorkoutLogID         uint             `gorm:"not null;index"`
	SetNumber            int              `gorm:"not null"`
	Status               WorkoutLogStatus `gorm:"not null;default:'planning'"`
	StatusChangedAt      *time.Time
	BreakAfterSeconds    *int
	TargetReps           *int
	TargetWeight         *float64
	TargetDuration       *int
	TargetDistance       *float64
	TargetTime           *int
	ActualReps           *int
	ActualWeight         *float64
	ActualDuration       *int
	ActualDistance       *float64
	ActualTime           *int
}

func (WorkoutLogExerciseSetEntity) TableName() string { return "workout_log_exercise_sets" }

func (e *WorkoutLogExerciseSetEntity) ToDTO() WorkoutLogExerciseSet {
	return WorkoutLogExerciseSet{
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
		ActualReps:           e.ActualReps,
		ActualWeight:         e.ActualWeight,
		ActualDuration:       e.ActualDuration,
		ActualDistance:       e.ActualDistance,
		ActualTime:           e.ActualTime,
	}
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
		ActualReps:           dto.ActualReps,
		ActualWeight:         dto.ActualWeight,
		ActualDuration:       dto.ActualDuration,
		ActualDistance:       dto.ActualDistance,
		ActualTime:           dto.ActualTime,
	}
}
