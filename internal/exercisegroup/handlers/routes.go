package handlers

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// RegisterRoutes registers all exercise group and exercise group member endpoints on the huma API.
func RegisterRoutes(api huma.API) {
	// --- Exercise groups ---

	huma.Register(api, huma.Operation{
		OperationID: "list-exercise-groups",
		Method:      http.MethodGet,
		Path:        "/exercise-groups",
		Tags:        []string{"exercise-groups"},
		Summary:     "List exercise groups",
	}, ListExerciseGroups)

	huma.Register(api, huma.Operation{
		OperationID:   "create-exercise-group",
		Method:        http.MethodPost,
		Path:          "/exercise-groups",
		Tags:          []string{"exercise-groups"},
		Summary:       "Create exercise group",
		DefaultStatus: 201,
	}, CreateExerciseGroup)

	huma.Register(api, huma.Operation{
		OperationID: "get-exercise-group",
		Method:      http.MethodGet,
		Path:        "/exercise-groups/{id}",
		Tags:        []string{"exercise-groups"},
		Summary:     "Get exercise group",
	}, GetExerciseGroup)

	huma.Register(api, huma.Operation{
		OperationID: "update-exercise-group",
		Method:      http.MethodPut,
		Path:        "/exercise-groups/{id}",
		Tags:        []string{"exercise-groups"},
		Summary:     "Update exercise group",
	}, UpdateExerciseGroup)

	huma.Register(api, huma.Operation{
		OperationID: "delete-exercise-group",
		Method:      http.MethodDelete,
		Path:        "/exercise-groups/{id}",
		Tags:        []string{"exercise-groups"},
		Summary:     "Delete exercise group",
	}, DeleteExerciseGroup)

	huma.Register(api, huma.Operation{
		OperationID: "get-exercise-group-permissions",
		Method:      http.MethodGet,
		Path:        "/exercise-groups/{id}/permissions",
		Tags:        []string{"exercise-groups"},
		Summary:     "Get exercise group permissions",
	}, GetExerciseGroupPermissions)

	// --- Exercise group members ---

	huma.Register(api, huma.Operation{
		OperationID: "list-exercise-group-members",
		Method:      http.MethodGet,
		Path:        "/exercise-group-members",
		Tags:        []string{"exercise-group-members"},
		Summary:     "List exercise group members",
	}, ListExerciseGroupMembers)

	huma.Register(api, huma.Operation{
		OperationID:   "create-exercise-group-member",
		Method:        http.MethodPost,
		Path:          "/exercise-group-members",
		Tags:          []string{"exercise-group-members"},
		Summary:       "Create exercise group member",
		DefaultStatus: 201,
	}, CreateExerciseGroupMember)

	huma.Register(api, huma.Operation{
		OperationID: "delete-exercise-group-member",
		Method:      http.MethodDelete,
		Path:        "/exercise-group-members/{id}",
		Tags:        []string{"exercise-group-members"},
		Summary:     "Delete exercise group member",
	}, DeleteExerciseGroupMember)
}
