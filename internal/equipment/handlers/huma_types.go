package handlers

import (
	"gesitr/internal/equipment/models"
	"gesitr/internal/humaconfig"
	"gesitr/internal/shared"
)

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

// RawBody skips huma's automatic validation — the Equipment DTO is shared
// between request and response and has server-set fields (id, createdAt, etc.)
// that aren't present in create requests.
type CreateEquipmentInput struct {
	RawBody []byte
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
	ID      uint `path:"id"`
	RawBody []byte
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
