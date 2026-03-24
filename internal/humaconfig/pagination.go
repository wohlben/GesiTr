package humaconfig

import "gesitr/internal/shared"

// PaginationInput can be embedded in huma input structs to add
// limit/offset query parameters with validation.
type PaginationInput struct {
	Limit  int `query:"limit" default:"50" minimum:"1" maximum:"200" doc:"Max items to return"`
	Offset int `query:"offset" default:"0" minimum:"0" doc:"Number of items to skip"`
}

// ToPaginationParams converts to the shared type used by ApplyPagination.
func (p PaginationInput) ToPaginationParams() shared.PaginationParams {
	return shared.PaginationParams{Limit: p.Limit, Offset: p.Offset}
}
