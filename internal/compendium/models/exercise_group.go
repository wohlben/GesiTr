package models

type ExerciseGroup struct {
	BaseModel   `tstype:",extends"`
	TemplateID  string  `json:"templateId"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	CreatedBy   string  `json:"createdBy"`
}
