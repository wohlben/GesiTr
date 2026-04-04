package models

import "gesitr/internal/shared"

type WorkoutGroupInfo struct {
	GroupName  string `json:"groupName"`
	Membership string `json:"membership"`
}

type Workout struct {
	shared.BaseModel `tstype:",extends"`
	OwnershipGroupID uint              `json:"ownershipGroupId"`
	Name             string            `json:"name"`
	Notes            *string           `json:"notes"`
	Public           bool              `json:"public"`
	Version          int               `json:"version"`
	Sections         []WorkoutSection  `json:"sections" gorm:"-"`
	WorkoutGroup     *WorkoutGroupInfo `json:"workoutGroup,omitempty" gorm:"-"`
}

type WorkoutSectionType string

const (
	WorkoutSectionTypeMain          WorkoutSectionType = "main"
	WorkoutSectionTypeSupplementary WorkoutSectionType = "supplementary"
)

type WorkoutSection struct {
	shared.BaseModel     `tstype:",extends"`
	WorkoutID            uint                 `json:"workoutId"`
	Type                 WorkoutSectionType   `json:"type"`
	Label                *string              `json:"label"`
	Position             int                  `json:"position"`
	RestBetweenExercises *int                 `json:"restBetweenExercises"`
	Items                []WorkoutSectionItem `json:"items" gorm:"-"`
}

type WorkoutSectionItemType string

const (
	WorkoutSectionItemTypeExercise      WorkoutSectionItemType = "exercise"
	WorkoutSectionItemTypeExerciseGroup WorkoutSectionItemType = "exercise_group"
)

type WorkoutSectionItem struct {
	shared.BaseModel `tstype:",extends"`
	WorkoutSectionID uint                   `json:"workoutSectionId"`
	Type             WorkoutSectionItemType `json:"type"`
	ExerciseID       *uint                  `json:"exerciseId"`
	ExerciseGroupID  *uint                  `json:"exerciseGroupId"`
	Data             *string                `json:"data"`
	Position         int                    `json:"position"`
}
