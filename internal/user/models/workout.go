package models

type Workout struct {
	BaseModel `tstype:",extends"`
	Owner     string           `json:"owner"`
	Name      string           `json:"name"`
	Notes     *string          `json:"notes"`
	Sections  []WorkoutSection `json:"sections" gorm:"-"`
}
