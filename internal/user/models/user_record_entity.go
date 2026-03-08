package models

type UserRecordEntity struct {
	BaseModel
	UserExerciseID          uint    `gorm:"not null;uniqueIndex:idx_user_exercise_measurement"`
	MeasurementType         string  `gorm:"not null;uniqueIndex:idx_user_exercise_measurement"`
	RecordValue             float64 `gorm:"not null"`
	ActualReps              *int
	ActualWeight            *float64
	ActualDuration          *int
	ActualDistance          *float64
	ActualTime              *int
	WorkoutLogExerciseSetID uint `gorm:"not null"`
}

func (UserRecordEntity) TableName() string { return "user_records" }

func (e *UserRecordEntity) ToDTO() UserRecord {
	return UserRecord{
		BaseModel:               e.BaseModel,
		UserExerciseID:          e.UserExerciseID,
		MeasurementType:         e.MeasurementType,
		RecordValue:             e.RecordValue,
		ActualReps:              e.ActualReps,
		ActualWeight:            e.ActualWeight,
		ActualDuration:          e.ActualDuration,
		ActualDistance:          e.ActualDistance,
		ActualTime:              e.ActualTime,
		WorkoutLogExerciseSetID: e.WorkoutLogExerciseSetID,
	}
}

func UserRecordFromDTO(dto UserRecord) UserRecordEntity {
	return UserRecordEntity{
		BaseModel:               dto.BaseModel,
		UserExerciseID:          dto.UserExerciseID,
		MeasurementType:         dto.MeasurementType,
		RecordValue:             dto.RecordValue,
		ActualReps:              dto.ActualReps,
		ActualWeight:            dto.ActualWeight,
		ActualDuration:          dto.ActualDuration,
		ActualDistance:          dto.ActualDistance,
		ActualTime:              dto.ActualTime,
		WorkoutLogExerciseSetID: dto.WorkoutLogExerciseSetID,
	}
}
