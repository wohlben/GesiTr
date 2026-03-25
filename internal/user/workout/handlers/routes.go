package handlers

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// RegisterRoutes registers all workout, section, and section-exercise endpoints on the huma API.
func RegisterRoutes(api huma.API) {
	// --- Workouts ---

	huma.Register(api, huma.Operation{
		OperationID: "ListWorkouts",
		Method:      http.MethodGet,
		Path:        "/user/workouts",
		Tags:        []string{"workouts"},
		Summary:     "List workouts",
	}, ListWorkouts)

	huma.Register(api, huma.Operation{
		OperationID:   "CreateWorkout",
		Method:        http.MethodPost,
		Path:          "/user/workouts",
		Tags:          []string{"workouts"},
		Summary:       "Create workout",
		DefaultStatus: 201,
	}, CreateWorkout)

	huma.Register(api, huma.Operation{
		OperationID: "GetWorkout",
		Method:      http.MethodGet,
		Path:        "/user/workouts/{id}",
		Tags:        []string{"workouts"},
		Summary:     "Get workout",
	}, GetWorkout)

	huma.Register(api, huma.Operation{
		OperationID: "UpdateWorkout",
		Method:      http.MethodPut,
		Path:        "/user/workouts/{id}",
		Tags:        []string{"workouts"},
		Summary:     "Update workout",
	}, UpdateWorkout)

	huma.Register(api, huma.Operation{
		OperationID: "DeleteWorkout",
		Method:      http.MethodDelete,
		Path:        "/user/workouts/{id}",
		Tags:        []string{"workouts"},
		Summary:     "Delete workout",
	}, DeleteWorkout)

	// --- Workout sections ---

	huma.Register(api, huma.Operation{
		OperationID: "ListWorkoutSections",
		Method:      http.MethodGet,
		Path:        "/user/workout-sections",
		Tags:        []string{"workout-sections"},
		Summary:     "List workout sections",
	}, ListWorkoutSections)

	huma.Register(api, huma.Operation{
		OperationID:   "CreateWorkoutSection",
		Method:        http.MethodPost,
		Path:          "/user/workout-sections",
		Tags:          []string{"workout-sections"},
		Summary:       "Create workout section",
		DefaultStatus: 201,
	}, CreateWorkoutSection)

	huma.Register(api, huma.Operation{
		OperationID: "GetWorkoutSection",
		Method:      http.MethodGet,
		Path:        "/user/workout-sections/{id}",
		Tags:        []string{"workout-sections"},
		Summary:     "Get workout section",
	}, GetWorkoutSection)

	huma.Register(api, huma.Operation{
		OperationID: "DeleteWorkoutSection",
		Method:      http.MethodDelete,
		Path:        "/user/workout-sections/{id}",
		Tags:        []string{"workout-sections"},
		Summary:     "Delete workout section",
	}, DeleteWorkoutSection)

	// --- Workout section items ---

	huma.Register(api, huma.Operation{
		OperationID: "ListWorkoutSectionItems",
		Method:      http.MethodGet,
		Path:        "/user/workout-section-items",
		Tags:        []string{"workout-section-items"},
		Summary:     "List workout section items",
	}, ListWorkoutSectionItems)

	huma.Register(api, huma.Operation{
		OperationID:   "CreateWorkoutSectionItem",
		Method:        http.MethodPost,
		Path:          "/user/workout-section-items",
		Tags:          []string{"workout-section-items"},
		Summary:       "Create workout section item",
		DefaultStatus: 201,
	}, CreateWorkoutSectionItem)

	huma.Register(api, huma.Operation{
		OperationID: "DeleteWorkoutSectionItem",
		Method:      http.MethodDelete,
		Path:        "/user/workout-section-items/{id}",
		Tags:        []string{"workout-section-items"},
		Summary:     "Delete workout section item",
	}, DeleteWorkoutSectionItem)
}
