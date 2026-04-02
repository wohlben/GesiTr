package handlers

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// RegisterRoutes registers all equipment, fulfillment, and relationship endpoints on the huma API.
func RegisterRoutes(api huma.API) {
	// Equipment
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

	// Fulfillments
	huma.Register(api, huma.Operation{
		OperationID: "ListFulfillments",
		Method:      http.MethodGet,
		Path:        "/fulfillments",
		Tags:        []string{"fulfillments"},
		Summary:     "List fulfillments",
	}, ListFulfillments)

	huma.Register(api, huma.Operation{
		OperationID:   "CreateFulfillment",
		Method:        http.MethodPost,
		Path:          "/fulfillments",
		Tags:          []string{"fulfillments"},
		Summary:       "Create fulfillment",
		DefaultStatus: 201,
	}, CreateFulfillment)

	huma.Register(api, huma.Operation{
		OperationID: "DeleteFulfillment",
		Method:      http.MethodDelete,
		Path:        "/fulfillments/{id}",
		Tags:        []string{"fulfillments"},
		Summary:     "Delete fulfillment",
	}, DeleteFulfillment)

	// Equipment relationships
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
