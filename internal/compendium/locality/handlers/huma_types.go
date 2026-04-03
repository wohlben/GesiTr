package handlers

import (
	"gesitr/internal/compendium/locality/models"
	"gesitr/internal/humaconfig"
	"gesitr/internal/shared"
)

type LocalityBody struct {
	Name   string `json:"name" required:"true"`
	Public bool   `json:"public" required:"false"`
}

type ListLocalitiesInput struct {
	humaconfig.PaginationInput
	Owner  string `query:"owner" doc:"Filter by owner ('me' for current user)"`
	Public string `query:"public" doc:"'true' for public only, 'false' for private only"`
	Q      string `query:"q" doc:"Search by name"`
}

type ListLocalitiesOutput struct {
	Body humaconfig.PaginatedBody[models.Locality]
}

type CreateLocalityInput struct {
	Body LocalityBody
}

type CreateLocalityOutput struct {
	Body models.Locality
}

type GetLocalityInput struct {
	ID uint `path:"id"`
}

type GetLocalityOutput struct {
	Body models.Locality
}

type UpdateLocalityInput struct {
	ID   uint `path:"id"`
	Body LocalityBody
}

type UpdateLocalityOutput struct {
	Body models.Locality
}

type DeleteLocalityInput struct {
	ID uint `path:"id"`
}

type DeleteLocalityOutput struct{}

type GetLocalityPermissionsInput struct {
	ID uint `path:"id"`
}

type GetLocalityPermissionsOutput struct {
	Body shared.PermissionsResponse
}
