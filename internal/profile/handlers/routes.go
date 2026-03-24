package handlers

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// RegisterRoutes registers all profile endpoints on the huma API.
func RegisterRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "get-my-profile",
		Method:      http.MethodGet,
		Path:        "/user/profile",
		Tags:        []string{"profiles"},
		Summary:     "Get current user's profile",
	}, GetMyProfile)

	huma.Register(api, huma.Operation{
		OperationID: "update-my-profile",
		Method:      http.MethodPatch,
		Path:        "/user/profile",
		Tags:        []string{"profiles"},
		Summary:     "Update current user's profile",
	}, UpdateMyProfile)

	huma.Register(api, huma.Operation{
		OperationID: "get-profile",
		Method:      http.MethodGet,
		Path:        "/profiles/{id}",
		Tags:        []string{"profiles"},
		Summary:     "Get a user's profile",
	}, GetProfile)
}
