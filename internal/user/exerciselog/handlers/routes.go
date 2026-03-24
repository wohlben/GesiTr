package handlers

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// RegisterRoutes registers all exercise log endpoints on the huma API.
func RegisterRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "list-exercise-logs",
		Method:      http.MethodGet,
		Path:        "/user/exercise-logs",
		Tags:        []string{"exercise-logs"},
		Summary:     "List exercise logs",
	}, ListExerciseLogs)

	huma.Register(api, huma.Operation{
		OperationID:   "create-exercise-log",
		Method:        http.MethodPost,
		Path:          "/user/exercise-logs",
		Tags:          []string{"exercise-logs"},
		Summary:       "Create exercise log",
		DefaultStatus: 201,
	}, CreateExerciseLog)

	huma.Register(api, huma.Operation{
		OperationID: "get-exercise-log",
		Method:      http.MethodGet,
		Path:        "/user/exercise-logs/{id}",
		Tags:        []string{"exercise-logs"},
		Summary:     "Get exercise log",
	}, GetExerciseLog)

	huma.Register(api, huma.Operation{
		OperationID: "update-exercise-log",
		Method:      http.MethodPatch,
		Path:        "/user/exercise-logs/{id}",
		Tags:        []string{"exercise-logs"},
		Summary:     "Update exercise log",
	}, UpdateExerciseLog)

	huma.Register(api, huma.Operation{
		OperationID: "delete-exercise-log",
		Method:      http.MethodDelete,
		Path:        "/user/exercise-logs/{id}",
		Tags:        []string{"exercise-logs"},
		Summary:     "Delete exercise log",
	}, DeleteExerciseLog)
}
