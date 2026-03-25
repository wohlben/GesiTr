package handlers

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// RegisterRoutes registers all equipment endpoints on the huma API.
func RegisterRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "ListEquipment",
		Method:      http.MethodGet,
		Path:        "/equipment",
		Tags:        []string{"equipment"},
		Summary:     "List equipment",
	}, ListEquipment)

	huma.Register(api, huma.Operation{
		OperationID:   "CreateEquipment",
		Method:        http.MethodPost,
		Path:          "/equipment",
		Tags:          []string{"equipment"},
		Summary:       "Create equipment",
		DefaultStatus: 201,
	}, CreateEquipment)

	huma.Register(api, huma.Operation{
		OperationID: "GetEquipment",
		Method:      http.MethodGet,
		Path:        "/equipment/{id}",
		Tags:        []string{"equipment"},
		Summary:     "Get equipment",
	}, GetEquipment)

	huma.Register(api, huma.Operation{
		OperationID: "UpdateEquipment",
		Method:      http.MethodPut,
		Path:        "/equipment/{id}",
		Tags:        []string{"equipment"},
		Summary:     "Update equipment",
	}, UpdateEquipment)

	huma.Register(api, huma.Operation{
		OperationID: "DeleteEquipment",
		Method:      http.MethodDelete,
		Path:        "/equipment/{id}",
		Tags:        []string{"equipment"},
		Summary:     "Delete equipment",
	}, DeleteEquipment)

	huma.Register(api, huma.Operation{
		OperationID: "GetEquipmentPermissions",
		Method:      http.MethodGet,
		Path:        "/equipment/{id}/permissions",
		Tags:        []string{"equipment"},
		Summary:     "Get equipment permissions",
	}, GetEquipmentPermissions)

	huma.Register(api, huma.Operation{
		OperationID: "ListEquipmentVersions",
		Method:      http.MethodGet,
		Path:        "/equipment/{id}/versions",
		Tags:        []string{"equipment"},
		Summary:     "List equipment versions",
	}, ListEquipmentVersions)

	huma.Register(api, huma.Operation{
		OperationID: "GetEquipmentVersion",
		Method:      http.MethodGet,
		Path:        "/equipment/{id}/versions/{version}",
		Tags:        []string{"equipment"},
		Summary:     "Get equipment version",
	}, GetEquipmentVersion)
}
