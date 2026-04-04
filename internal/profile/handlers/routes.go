package handlers

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "GetProfile",
		Method:      http.MethodGet,
		Path:        "/profile",
		Tags:        []string{"profile"},
		Summary:     "Get current user's profile",
	}, GetProfile)

	huma.Register(api, huma.Operation{
		OperationID:   "CreateProfile",
		Method:        http.MethodPost,
		Path:          "/profile",
		Tags:          []string{"profile"},
		Summary:       "Create profile for current user",
		DefaultStatus: 201,
	}, CreateProfile)
}
