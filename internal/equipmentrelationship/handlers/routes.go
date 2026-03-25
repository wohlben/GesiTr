package handlers

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// RegisterRoutes registers all equipment relationship endpoints on the huma API.
func RegisterRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "ListEquipmentRelationships",
		Method:      http.MethodGet,
		Path:        "/equipment-relationships",
		Tags:        []string{"equipment-relationships"},
		Summary:     "List equipment relationships",
	}, ListEquipmentRelationships)

	huma.Register(api, huma.Operation{
		OperationID:   "CreateEquipmentRelationship",
		Method:        http.MethodPost,
		Path:          "/equipment-relationships",
		Tags:          []string{"equipment-relationships"},
		Summary:       "Create equipment relationship",
		DefaultStatus: 201,
	}, CreateEquipmentRelationship)

	huma.Register(api, huma.Operation{
		OperationID: "DeleteEquipmentRelationship",
		Method:      http.MethodDelete,
		Path:        "/equipment-relationships/{id}",
		Tags:        []string{"equipment-relationships"},
		Summary:     "Delete equipment relationship",
	}, DeleteEquipmentRelationship)
}
