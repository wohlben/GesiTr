package handlers

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// RegisterRoutes registers all exercise and exercise-scheme endpoints on the huma API.
func RegisterRoutes(api huma.API) {
	// --- Exercises ---

	huma.Register(api, huma.Operation{
		OperationID: "ListExercises",
		Method:      http.MethodGet,
		Path:        "/exercises",
		Tags:        []string{"exercises"},
		Summary:     "List exercises",
	}, ListExercises)

	huma.Register(api, huma.Operation{
		OperationID:   "CreateExercise",
		Method:        http.MethodPost,
		Path:          "/exercises",
		Tags:          []string{"exercises"},
		Summary:       "Create exercise",
		DefaultStatus: 201,
	}, CreateExercise)

	huma.Register(api, huma.Operation{
		OperationID: "GetExercise",
		Method:      http.MethodGet,
		Path:        "/exercises/{id}",
		Tags:        []string{"exercises"},
		Summary:     "Get exercise",
	}, GetExercise)

	huma.Register(api, huma.Operation{
		OperationID: "UpdateExercise",
		Method:      http.MethodPut,
		Path:        "/exercises/{id}",
		Tags:        []string{"exercises"},
		Summary:     "Update exercise",
	}, UpdateExercise)

	huma.Register(api, huma.Operation{
		OperationID: "DeleteExercise",
		Method:      http.MethodDelete,
		Path:        "/exercises/{id}",
		Tags:        []string{"exercises"},
		Summary:     "Delete exercise",
	}, DeleteExercise)

	huma.Register(api, huma.Operation{
		OperationID: "GetExercisePermissions",
		Method:      http.MethodGet,
		Path:        "/exercises/{id}/permissions",
		Tags:        []string{"exercises"},
		Summary:     "Get exercise permissions",
	}, GetExercisePermissions)

	huma.Register(api, huma.Operation{
		OperationID: "ListExerciseVersions",
		Method:      http.MethodGet,
		Path:        "/exercises/{id}/versions",
		Tags:        []string{"exercises"},
		Summary:     "List exercise versions",
	}, ListExerciseVersions)

	huma.Register(api, huma.Operation{
		OperationID: "GetExerciseVersion",
		Method:      http.MethodGet,
		Path:        "/exercises/{id}/versions/{version}",
		Tags:        []string{"exercises"},
		Summary:     "Get exercise version",
	}, GetExerciseVersion)

	huma.Register(api, huma.Operation{
		OperationID: "DeleteExerciseVersion",
		Method:      http.MethodDelete,
		Path:        "/exercises/{id}/versions/{version}",
		Tags:        []string{"exercises"},
		Summary:     "Delete exercise version",
	}, DeleteExerciseVersion)

	huma.Register(api, huma.Operation{
		OperationID: "DeleteAllExerciseVersions",
		Method:      http.MethodDelete,
		Path:        "/exercises/{id}/versions",
		Tags:        []string{"exercises"},
		Summary:     "Delete all exercise versions",
	}, DeleteAllExerciseVersions)

	// --- Exercise schemes ---

	huma.Register(api, huma.Operation{
		OperationID: "ListExerciseSchemes",
		Method:      http.MethodGet,
		Path:        "/exercise-schemes",
		Tags:        []string{"exercise-schemes"},
		Summary:     "List exercise schemes",
	}, ListExerciseSchemes)

	huma.Register(api, huma.Operation{
		OperationID:   "CreateExerciseScheme",
		Method:        http.MethodPost,
		Path:          "/exercise-schemes",
		Tags:          []string{"exercise-schemes"},
		Summary:       "Create exercise scheme",
		DefaultStatus: 201,
	}, CreateExerciseScheme)

	huma.Register(api, huma.Operation{
		OperationID: "GetExerciseScheme",
		Method:      http.MethodGet,
		Path:        "/exercise-schemes/{id}",
		Tags:        []string{"exercise-schemes"},
		Summary:     "Get exercise scheme",
	}, GetExerciseScheme)

	huma.Register(api, huma.Operation{
		OperationID: "UpdateExerciseScheme",
		Method:      http.MethodPut,
		Path:        "/exercise-schemes/{id}",
		Tags:        []string{"exercise-schemes"},
		Summary:     "Update exercise scheme",
	}, UpdateExerciseScheme)

	huma.Register(api, huma.Operation{
		OperationID: "DeleteExerciseScheme",
		Method:      http.MethodDelete,
		Path:        "/exercise-schemes/{id}",
		Tags:        []string{"exercise-schemes"},
		Summary:     "Delete exercise scheme",
	}, DeleteExerciseScheme)

	// --- Exercise relationships ---

	huma.Register(api, huma.Operation{
		OperationID: "ListExerciseRelationships",
		Method:      http.MethodGet,
		Path:        "/exercise-relationships",
		Tags:        []string{"exercise-relationships"},
		Summary:     "List exercise relationships",
	}, ListExerciseRelationships)

	huma.Register(api, huma.Operation{
		OperationID:   "CreateExerciseRelationship",
		Method:        http.MethodPost,
		Path:          "/exercise-relationships",
		Tags:          []string{"exercise-relationships"},
		Summary:       "Create exercise relationship",
		DefaultStatus: 201,
	}, CreateExerciseRelationship)

	huma.Register(api, huma.Operation{
		OperationID: "DeleteExerciseRelationship",
		Method:      http.MethodDelete,
		Path:        "/exercise-relationships/{id}",
		Tags:        []string{"exercise-relationships"},
		Summary:     "Delete exercise relationship",
	}, DeleteExerciseRelationship)
}
