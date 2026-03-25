package handlers

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// RegisterRoutes registers all equipment fulfillment endpoints on the huma API.
func RegisterRoutes(api huma.API) {
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
}
