package handlers

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// RegisterRoutes registers all equipment relationship endpoints on the huma API.
func RegisterRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "list-equipment-relationships",
		Method:      http.MethodGet,
		Path:        "/equipment-relationships",
		Tags:        []string{"equipment-relationships"},
		Summary:     "List equipment relationships",
	}, ListEquipmentRelationships)

	huma.Register(api, huma.Operation{
		OperationID:   "create-equipment-relationship",
		Method:        http.MethodPost,
		Path:          "/equipment-relationships",
		Tags:          []string{"equipment-relationships"},
		Summary:       "Create equipment relationship",
		DefaultStatus: 201,
	}, CreateEquipmentRelationship)

	huma.Register(api, huma.Operation{
		OperationID: "delete-equipment-relationship",
		Method:      http.MethodDelete,
		Path:        "/equipment-relationships/{id}",
		Tags:        []string{"equipment-relationships"},
		Summary:     "Delete equipment relationship",
	}, DeleteEquipmentRelationship)
}
