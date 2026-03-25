package handlers

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// RegisterRoutes registers all exercise group and exercise group member endpoints on the huma API.
func RegisterRoutes(api huma.API) {
	// --- Exercise groups ---

	huma.Register(api, huma.Operation{
		OperationID: "ListExerciseGroups",
		Method:      http.MethodGet,
		Path:        "/exercise-groups",
		Tags:        []string{"exercise-groups"},
		Summary:     "List exercise groups",
	}, ListExerciseGroups)

	huma.Register(api, huma.Operation{
		OperationID:   "CreateExerciseGroup",
		Method:        http.MethodPost,
		Path:          "/exercise-groups",
		Tags:          []string{"exercise-groups"},
		Summary:       "Create exercise group",
		DefaultStatus: 201,
	}, CreateExerciseGroup)

	huma.Register(api, huma.Operation{
		OperationID: "GetExerciseGroup",
		Method:      http.MethodGet,
		Path:        "/exercise-groups/{id}",
		Tags:        []string{"exercise-groups"},
		Summary:     "Get exercise group",
	}, GetExerciseGroup)

	huma.Register(api, huma.Operation{
		OperationID: "UpdateExerciseGroup",
		Method:      http.MethodPut,
		Path:        "/exercise-groups/{id}",
		Tags:        []string{"exercise-groups"},
		Summary:     "Update exercise group",
	}, UpdateExerciseGroup)

	huma.Register(api, huma.Operation{
		OperationID: "DeleteExerciseGroup",
		Method:      http.MethodDelete,
		Path:        "/exercise-groups/{id}",
		Tags:        []string{"exercise-groups"},
		Summary:     "Delete exercise group",
	}, DeleteExerciseGroup)

	huma.Register(api, huma.Operation{
		OperationID: "GetExerciseGroupPermissions",
		Method:      http.MethodGet,
		Path:        "/exercise-groups/{id}/permissions",
		Tags:        []string{"exercise-groups"},
		Summary:     "Get exercise group permissions",
	}, GetExerciseGroupPermissions)

	// --- Exercise group members ---

	huma.Register(api, huma.Operation{
		OperationID: "ListExerciseGroupMembers",
		Method:      http.MethodGet,
		Path:        "/exercise-group-members",
		Tags:        []string{"exercise-group-members"},
		Summary:     "List exercise group members",
	}, ListExerciseGroupMembers)

	huma.Register(api, huma.Operation{
		OperationID:   "CreateExerciseGroupMember",
		Method:        http.MethodPost,
		Path:          "/exercise-group-members",
		Tags:          []string{"exercise-group-members"},
		Summary:       "Create exercise group member",
		DefaultStatus: 201,
	}, CreateExerciseGroupMember)

	huma.Register(api, huma.Operation{
		OperationID: "DeleteExerciseGroupMember",
		Method:      http.MethodDelete,
		Path:        "/exercise-group-members/{id}",
		Tags:        []string{"exercise-group-members"},
		Summary:     "Delete exercise group member",
	}, DeleteExerciseGroupMember)
}
