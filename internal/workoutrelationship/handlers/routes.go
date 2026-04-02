package handlers

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "ListWorkoutRelationships",
		Method:      http.MethodGet,
		Path:        "/workout-relationships",
		Tags:        []string{"workout-relationships"},
		Summary:     "List workout relationships",
	}, ListWorkoutRelationships)

	huma.Register(api, huma.Operation{
		OperationID:   "CreateWorkoutRelationship",
		Method:        http.MethodPost,
		Path:          "/workout-relationships",
		Tags:          []string{"workout-relationships"},
		Summary:       "Create workout relationship",
		DefaultStatus: 201,
	}, CreateWorkoutRelationship)

	huma.Register(api, huma.Operation{
		OperationID: "DeleteWorkoutRelationship",
		Method:      http.MethodDelete,
		Path:        "/workout-relationships/{id}",
		Tags:        []string{"workout-relationships"},
		Summary:     "Delete workout relationship",
	}, DeleteWorkoutRelationship)
}
