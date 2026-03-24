package handlers

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// RegisterRoutes registers all workout, section, and section-exercise endpoints on the huma API.
func RegisterRoutes(api huma.API) {
	// --- Workouts ---

	huma.Register(api, huma.Operation{
		OperationID: "list-workouts",
		Method:      http.MethodGet,
		Path:        "/user/workouts",
		Tags:        []string{"workouts"},
		Summary:     "List workouts",
	}, ListWorkouts)

	huma.Register(api, huma.Operation{
		OperationID:   "create-workout",
		Method:        http.MethodPost,
		Path:          "/user/workouts",
		Tags:          []string{"workouts"},
		Summary:       "Create workout",
		DefaultStatus: 201,
	}, CreateWorkout)

	huma.Register(api, huma.Operation{
		OperationID: "get-workout",
		Method:      http.MethodGet,
		Path:        "/user/workouts/{id}",
		Tags:        []string{"workouts"},
		Summary:     "Get workout",
	}, GetWorkout)

	huma.Register(api, huma.Operation{
		OperationID: "update-workout",
		Method:      http.MethodPut,
		Path:        "/user/workouts/{id}",
		Tags:        []string{"workouts"},
		Summary:     "Update workout",
	}, UpdateWorkout)

	huma.Register(api, huma.Operation{
		OperationID: "delete-workout",
		Method:      http.MethodDelete,
		Path:        "/user/workouts/{id}",
		Tags:        []string{"workouts"},
		Summary:     "Delete workout",
	}, DeleteWorkout)

	// --- Workout sections ---

	huma.Register(api, huma.Operation{
		OperationID: "list-workout-sections",
		Method:      http.MethodGet,
		Path:        "/user/workout-sections",
		Tags:        []string{"workout-sections"},
		Summary:     "List workout sections",
	}, ListWorkoutSections)

	huma.Register(api, huma.Operation{
		OperationID:   "create-workout-section",
		Method:        http.MethodPost,
		Path:          "/user/workout-sections",
		Tags:          []string{"workout-sections"},
		Summary:       "Create workout section",
		DefaultStatus: 201,
	}, CreateWorkoutSection)

	huma.Register(api, huma.Operation{
		OperationID: "get-workout-section",
		Method:      http.MethodGet,
		Path:        "/user/workout-sections/{id}",
		Tags:        []string{"workout-sections"},
		Summary:     "Get workout section",
	}, GetWorkoutSection)

	huma.Register(api, huma.Operation{
		OperationID: "delete-workout-section",
		Method:      http.MethodDelete,
		Path:        "/user/workout-sections/{id}",
		Tags:        []string{"workout-sections"},
		Summary:     "Delete workout section",
	}, DeleteWorkoutSection)

	// --- Workout section exercises ---

	huma.Register(api, huma.Operation{
		OperationID: "list-workout-section-exercises",
		Method:      http.MethodGet,
		Path:        "/user/workout-section-exercises",
		Tags:        []string{"workout-section-exercises"},
		Summary:     "List workout section exercises",
	}, ListWorkoutSectionExercises)

	huma.Register(api, huma.Operation{
		OperationID:   "create-workout-section-exercise",
		Method:        http.MethodPost,
		Path:          "/user/workout-section-exercises",
		Tags:          []string{"workout-section-exercises"},
		Summary:       "Create workout section exercise",
		DefaultStatus: 201,
	}, CreateWorkoutSectionExercise)

	huma.Register(api, huma.Operation{
		OperationID: "delete-workout-section-exercise",
		Method:      http.MethodDelete,
		Path:        "/user/workout-section-exercises/{id}",
		Tags:        []string{"workout-section-exercises"},
		Summary:     "Delete workout section exercise",
	}, DeleteWorkoutSectionExercise)
}
