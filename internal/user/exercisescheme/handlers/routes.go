package handlers

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// RegisterRoutes registers all exercise-scheme endpoints on the huma API.
func RegisterRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "ListExerciseSchemes",
		Method:      http.MethodGet,
		Path:        "/user/exercise-schemes",
		Tags:        []string{"exercise-schemes"},
		Summary:     "List exercise schemes",
	}, ListExerciseSchemes)

	huma.Register(api, huma.Operation{
		OperationID:   "CreateExerciseScheme",
		Method:        http.MethodPost,
		Path:          "/user/exercise-schemes",
		Tags:          []string{"exercise-schemes"},
		Summary:       "Create exercise scheme",
		DefaultStatus: 201,
	}, CreateExerciseScheme)

	huma.Register(api, huma.Operation{
		OperationID: "GetExerciseScheme",
		Method:      http.MethodGet,
		Path:        "/user/exercise-schemes/{id}",
		Tags:        []string{"exercise-schemes"},
		Summary:     "Get exercise scheme",
	}, GetExerciseScheme)

	huma.Register(api, huma.Operation{
		OperationID: "UpdateExerciseScheme",
		Method:      http.MethodPut,
		Path:        "/user/exercise-schemes/{id}",
		Tags:        []string{"exercise-schemes"},
		Summary:     "Update exercise scheme",
	}, UpdateExerciseScheme)

	huma.Register(api, huma.Operation{
		OperationID: "DeleteExerciseScheme",
		Method:      http.MethodDelete,
		Path:        "/user/exercise-schemes/{id}",
		Tags:        []string{"exercise-schemes"},
		Summary:     "Delete exercise scheme",
	}, DeleteExerciseScheme)
}
