package handlers

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// RegisterRoutes registers exercise name preference endpoints on the huma API.
func RegisterRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "ListExerciseNamePreferences",
		Method:      http.MethodGet,
		Path:        "/user/exercise-name-preferences",
		Tags:        []string{"exercise-name-preferences"},
		Summary:     "List exercise name preferences",
	}, ListExerciseNamePreferences)

	huma.Register(api, huma.Operation{
		OperationID: "SetExerciseNamePreference",
		Method:      http.MethodPut,
		Path:        "/user/exercise-name-preferences/{exerciseId}",
		Tags:        []string{"exercise-name-preferences"},
		Summary:     "Set preferred exercise name",
	}, SetExerciseNamePreference)

	huma.Register(api, huma.Operation{
		OperationID: "DeleteExerciseNamePreference",
		Method:      http.MethodDelete,
		Path:        "/user/exercise-name-preferences/{exerciseId}",
		Tags:        []string{"exercise-name-preferences"},
		Summary:     "Delete exercise name preference",
	}, DeleteExerciseNamePreference)
}
