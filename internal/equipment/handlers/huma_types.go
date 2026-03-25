package handlers

import (
	"gesitr/internal/equipment/models"
	"gesitr/internal/humaconfig"
	"gesitr/internal/shared"
)

type EquipmentBody struct {
	Name        string                   `json:"name" required:"true"`
	DisplayName string                   `json:"displayName" required:"true"`
	Description string                   `json:"description" required:"false"`
	Category    models.EquipmentCategory `json:"category" required:"true"`
	ImageUrl    *string                  `json:"imageUrl,omitempty"`
	TemplateID  string                   `json:"templateId" required:"false"`
	Public      bool                     `json:"public" required:"false"`
}

type ListEquipmentInput struct {
	humaconfig.PaginationInput
	Owner    string `query:"owner" doc:"Filter by owner ('me' for current user)"`
	Public   string `query:"public" doc:"'true' to show only public equipment"`
	Q        string `query:"q" doc:"Search by name or display name"`
	Category string `query:"category" doc:"Filter by equipment category"`
}

type ListEquipmentOutput struct {
	Body humaconfig.PaginatedBody[models.Equipment]
}

type CreateEquipmentInput struct {
	Body EquipmentBody
}

type CreateEquipmentOutput struct {
	Body models.Equipment
}

type GetEquipmentInput struct {
	ID uint `path:"id"`
}

type GetEquipmentOutput struct {
	Body models.Equipment
}

type UpdateEquipmentInput struct {
	ID   uint `path:"id"`
	Body EquipmentBody
}

type UpdateEquipmentOutput struct {
	Body models.Equipment
}

type DeleteEquipmentInput struct {
	ID uint `path:"id"`
}

type DeleteEquipmentOutput struct{}

type GetEquipmentPermissionsInput struct {
	ID uint `path:"id"`
}

type GetEquipmentPermissionsOutput struct {
	Body shared.PermissionsResponse
}

type ListEquipmentVersionsInput struct {
	ID uint `path:"id"`
}

type ListEquipmentVersionsOutput struct {
	Body []shared.VersionEntry
}

type GetEquipmentVersionInput struct {
	TemplateID string `path:"templateId"`
	Version    int    `path:"version"`
}

type GetEquipmentVersionOutput struct {
	Body shared.VersionEntry
}
