package handlers

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// RegisterRoutes registers all equipment endpoints on the huma API.
func RegisterRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "list-equipment",
		Method:      http.MethodGet,
		Path:        "/equipment",
		Tags:        []string{"equipment"},
		Summary:     "List equipment",
	}, ListEquipment)

	huma.Register(api, huma.Operation{
		OperationID:   "create-equipment",
		Method:        http.MethodPost,
		Path:          "/equipment",
		Tags:          []string{"equipment"},
		Summary:       "Create equipment",
		DefaultStatus: 201,
	}, CreateEquipment)

	huma.Register(api, huma.Operation{
		OperationID: "get-equipment",
		Method:      http.MethodGet,
		Path:        "/equipment/{id}",
		Tags:        []string{"equipment"},
		Summary:     "Get equipment",
	}, GetEquipment)

	huma.Register(api, huma.Operation{
		OperationID: "update-equipment",
		Method:      http.MethodPut,
		Path:        "/equipment/{id}",
		Tags:        []string{"equipment"},
		Summary:     "Update equipment",
	}, UpdateEquipment)

	huma.Register(api, huma.Operation{
		OperationID: "delete-equipment",
		Method:      http.MethodDelete,
		Path:        "/equipment/{id}",
		Tags:        []string{"equipment"},
		Summary:     "Delete equipment",
	}, DeleteEquipment)

	huma.Register(api, huma.Operation{
		OperationID: "get-equipment-permissions",
		Method:      http.MethodGet,
		Path:        "/equipment/{id}/permissions",
		Tags:        []string{"equipment"},
		Summary:     "Get equipment permissions",
	}, GetEquipmentPermissions)

	huma.Register(api, huma.Operation{
		OperationID: "list-equipment-versions",
		Method:      http.MethodGet,
		Path:        "/equipment/{id}/versions",
		Tags:        []string{"equipment"},
		Summary:     "List equipment versions",
	}, ListEquipmentVersions)

	huma.Register(api, huma.Operation{
		OperationID: "get-equipment-version",
		Method:      http.MethodGet,
		Path:        "/equipment/templates/{templateId}/versions/{version}",
		Tags:        []string{"equipment"},
		Summary:     "Get equipment version",
	}, GetEquipmentVersion)
}
