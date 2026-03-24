package handlers

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// RegisterRoutes registers all exercise relationship endpoints on the huma API.
func RegisterRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "list-exercise-relationships",
		Method:      http.MethodGet,
		Path:        "/exercise-relationships",
		Tags:        []string{"exercise-relationships"},
		Summary:     "List exercise relationships",
	}, ListExerciseRelationships)

	huma.Register(api, huma.Operation{
		OperationID:   "create-exercise-relationship",
		Method:        http.MethodPost,
		Path:          "/exercise-relationships",
		Tags:          []string{"exercise-relationships"},
		Summary:       "Create exercise relationship",
		DefaultStatus: 201,
	}, CreateExerciseRelationship)

	huma.Register(api, huma.Operation{
		OperationID: "delete-exercise-relationship",
		Method:      http.MethodDelete,
		Path:        "/exercise-relationships/{id}",
		Tags:        []string{"exercise-relationships"},
		Summary:     "Delete exercise relationship",
	}, DeleteExerciseRelationship)
}
