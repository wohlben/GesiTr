package handlers

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// RegisterRoutes registers all locality and locality-availability endpoints on the huma API.
func RegisterRoutes(api huma.API) {
	// Localities
	huma.Register(api, huma.Operation{
		OperationID: "ListLocalities",
		Method:      http.MethodGet,
		Path:        "/localities",
		Tags:        []string{"localities"},
		Summary:     "List localities",
	}, ListLocalities)

	huma.Register(api, huma.Operation{
		OperationID:   "CreateLocality",
		Method:        http.MethodPost,
		Path:          "/localities",
		Tags:          []string{"localities"},
		Summary:       "Create locality",
		DefaultStatus: 201,
	}, CreateLocality)

	huma.Register(api, huma.Operation{
		OperationID: "GetLocality",
		Method:      http.MethodGet,
		Path:        "/localities/{id}",
		Tags:        []string{"localities"},
		Summary:     "Get locality",
	}, GetLocality)

	huma.Register(api, huma.Operation{
		OperationID: "UpdateLocality",
		Method:      http.MethodPut,
		Path:        "/localities/{id}",
		Tags:        []string{"localities"},
		Summary:     "Update locality",
	}, UpdateLocality)

	huma.Register(api, huma.Operation{
		OperationID: "DeleteLocality",
		Method:      http.MethodDelete,
		Path:        "/localities/{id}",
		Tags:        []string{"localities"},
		Summary:     "Delete locality",
	}, DeleteLocality)

	huma.Register(api, huma.Operation{
		OperationID: "GetLocalityPermissions",
		Method:      http.MethodGet,
		Path:        "/localities/{id}/permissions",
		Tags:        []string{"localities"},
		Summary:     "Get locality permissions",
	}, GetLocalityPermissions)

	// Locality Availabilities
	huma.Register(api, huma.Operation{
		OperationID: "ListLocalityAvailabilities",
		Method:      http.MethodGet,
		Path:        "/locality-availabilities",
		Tags:        []string{"locality-availabilities"},
		Summary:     "List locality availabilities",
	}, ListLocalityAvailabilities)

	huma.Register(api, huma.Operation{
		OperationID:   "CreateLocalityAvailability",
		Method:        http.MethodPost,
		Path:          "/locality-availabilities",
		Tags:          []string{"locality-availabilities"},
		Summary:       "Create locality availability",
		DefaultStatus: 201,
	}, CreateLocalityAvailability)

	huma.Register(api, huma.Operation{
		OperationID: "UpdateLocalityAvailability",
		Method:      http.MethodPut,
		Path:        "/locality-availabilities/{id}",
		Tags:        []string{"locality-availabilities"},
		Summary:     "Update locality availability",
	}, UpdateLocalityAvailability)

	huma.Register(api, huma.Operation{
		OperationID: "DeleteLocalityAvailability",
		Method:      http.MethodDelete,
		Path:        "/locality-availabilities/{id}",
		Tags:        []string{"locality-availabilities"},
		Summary:     "Delete locality availability",
	}, DeleteLocalityAvailability)
}
