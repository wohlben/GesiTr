package models

type UserExercise struct {
	BaseModel          `tstype:",extends"`
	Owner              string `json:"owner"`
	ExerciseTemplateID string `json:"exerciseTemplateId"`
	CompendiumVersion  int    `json:"compendiumVersion"`
}
