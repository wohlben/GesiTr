package handlers

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// RegisterRoutes registers all workout log endpoints on the huma API.
func RegisterRoutes(api huma.API) {
	// --- Workout logs ---

	huma.Register(api, huma.Operation{
		OperationID: "list-workout-logs",
		Method:      http.MethodGet,
		Path:        "/user/workout-logs",
		Tags:        []string{"workout-logs"},
		Summary:     "List workout logs",
	}, ListWorkoutLogs)

	huma.Register(api, huma.Operation{
		OperationID:   "create-workout-log",
		Method:        http.MethodPost,
		Path:          "/user/workout-logs",
		Tags:          []string{"workout-logs"},
		Summary:       "Create workout log",
		DefaultStatus: 201,
	}, CreateWorkoutLog)

	huma.Register(api, huma.Operation{
		OperationID: "get-workout-log",
		Method:      http.MethodGet,
		Path:        "/user/workout-logs/{id}",
		Tags:        []string{"workout-logs"},
		Summary:     "Get workout log",
	}, GetWorkoutLog)

	huma.Register(api, huma.Operation{
		OperationID: "update-workout-log",
		Method:      http.MethodPatch,
		Path:        "/user/workout-logs/{id}",
		Tags:        []string{"workout-logs"},
		Summary:     "Update workout log",
	}, UpdateWorkoutLog)

	huma.Register(api, huma.Operation{
		OperationID: "delete-workout-log",
		Method:      http.MethodDelete,
		Path:        "/user/workout-logs/{id}",
		Tags:        []string{"workout-logs"},
		Summary:     "Delete workout log",
	}, DeleteWorkoutLog)

	huma.Register(api, huma.Operation{
		OperationID: "start-workout-log",
		Method:      http.MethodPost,
		Path:        "/user/workout-logs/{id}/start",
		Tags:        []string{"workout-logs"},
		Summary:     "Start workout log",
	}, StartWorkoutLog)

	huma.Register(api, huma.Operation{
		OperationID:   "start-adhoc-workout-log",
		Method:        http.MethodPost,
		Path:          "/user/workout-logs/adhoc",
		Tags:          []string{"workout-logs"},
		Summary:       "Start adhoc workout log",
		DefaultStatus: 201,
	}, StartAdhocWorkoutLog)

	huma.Register(api, huma.Operation{
		OperationID: "finish-workout-log",
		Method:      http.MethodPost,
		Path:        "/user/workout-logs/{id}/finish",
		Tags:        []string{"workout-logs"},
		Summary:     "Finish workout log",
	}, FinishWorkoutLog)

	huma.Register(api, huma.Operation{
		OperationID: "abandon-workout-log",
		Method:      http.MethodPost,
		Path:        "/user/workout-logs/{id}/abandon",
		Tags:        []string{"workout-logs"},
		Summary:     "Abandon workout log",
	}, AbandonWorkoutLog)

	// --- Workout log sections ---

	huma.Register(api, huma.Operation{
		OperationID: "list-workout-log-sections",
		Method:      http.MethodGet,
		Path:        "/user/workout-log-sections",
		Tags:        []string{"workout-log-sections"},
		Summary:     "List workout log sections",
	}, ListWorkoutLogSections)

	huma.Register(api, huma.Operation{
		OperationID:   "create-workout-log-section",
		Method:        http.MethodPost,
		Path:          "/user/workout-log-sections",
		Tags:          []string{"workout-log-sections"},
		Summary:       "Create workout log section",
		DefaultStatus: 201,
	}, CreateWorkoutLogSection)

	huma.Register(api, huma.Operation{
		OperationID: "get-workout-log-section",
		Method:      http.MethodGet,
		Path:        "/user/workout-log-sections/{id}",
		Tags:        []string{"workout-log-sections"},
		Summary:     "Get workout log section",
	}, GetWorkoutLogSection)

	huma.Register(api, huma.Operation{
		OperationID: "update-workout-log-section",
		Method:      http.MethodPatch,
		Path:        "/user/workout-log-sections/{id}",
		Tags:        []string{"workout-log-sections"},
		Summary:     "Update workout log section",
	}, UpdateWorkoutLogSection)

	huma.Register(api, huma.Operation{
		OperationID: "delete-workout-log-section",
		Method:      http.MethodDelete,
		Path:        "/user/workout-log-sections/{id}",
		Tags:        []string{"workout-log-sections"},
		Summary:     "Delete workout log section",
	}, DeleteWorkoutLogSection)

	// --- Workout log exercises ---

	huma.Register(api, huma.Operation{
		OperationID: "list-workout-log-exercises",
		Method:      http.MethodGet,
		Path:        "/user/workout-log-exercises",
		Tags:        []string{"workout-log-exercises"},
		Summary:     "List workout log exercises",
	}, ListWorkoutLogExercises)

	huma.Register(api, huma.Operation{
		OperationID:   "create-workout-log-exercise",
		Method:        http.MethodPost,
		Path:          "/user/workout-log-exercises",
		Tags:          []string{"workout-log-exercises"},
		Summary:       "Create workout log exercise",
		DefaultStatus: 201,
	}, CreateWorkoutLogExercise)

	huma.Register(api, huma.Operation{
		OperationID: "update-workout-log-exercise",
		Method:      http.MethodPatch,
		Path:        "/user/workout-log-exercises/{id}",
		Tags:        []string{"workout-log-exercises"},
		Summary:     "Update workout log exercise",
	}, UpdateWorkoutLogExercise)

	huma.Register(api, huma.Operation{
		OperationID: "delete-workout-log-exercise",
		Method:      http.MethodDelete,
		Path:        "/user/workout-log-exercises/{id}",
		Tags:        []string{"workout-log-exercises"},
		Summary:     "Delete workout log exercise",
	}, DeleteWorkoutLogExercise)

	// --- Workout log exercise sets ---

	huma.Register(api, huma.Operation{
		OperationID: "list-workout-log-exercise-sets",
		Method:      http.MethodGet,
		Path:        "/user/workout-log-exercise-sets",
		Tags:        []string{"workout-log-exercise-sets"},
		Summary:     "List workout log exercise sets",
	}, ListWorkoutLogExerciseSets)

	huma.Register(api, huma.Operation{
		OperationID:   "create-workout-log-exercise-set",
		Method:        http.MethodPost,
		Path:          "/user/workout-log-exercise-sets",
		Tags:          []string{"workout-log-exercise-sets"},
		Summary:       "Create workout log exercise set",
		DefaultStatus: 201,
	}, CreateWorkoutLogExerciseSet)

	huma.Register(api, huma.Operation{
		OperationID: "update-workout-log-exercise-set",
		Method:      http.MethodPatch,
		Path:        "/user/workout-log-exercise-sets/{id}",
		Tags:        []string{"workout-log-exercise-sets"},
		Summary:     "Update workout log exercise set",
	}, UpdateWorkoutLogExerciseSet)

	huma.Register(api, huma.Operation{
		OperationID: "delete-workout-log-exercise-set",
		Method:      http.MethodDelete,
		Path:        "/user/workout-log-exercise-sets/{id}",
		Tags:        []string{"workout-log-exercise-sets"},
		Summary:     "Delete workout log exercise set",
	}, DeleteWorkoutLogExerciseSet)
}
