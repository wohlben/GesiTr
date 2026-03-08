package models

type Equipment struct {
	BaseModel   `tstype:",extends"`
	Name        string            `json:"name"`
	DisplayName string            `json:"displayName"`
	Description string            `json:"description"`
	Category    EquipmentCategory `json:"category"`
	ImageUrl    *string           `json:"imageUrl"`
	TemplateID  string            `json:"templateId"`
	CreatedBy   string            `json:"createdBy"`
	Version     int               `json:"version"`
}
