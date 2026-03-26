package handlers

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// RegisterRoutes registers all workout group and membership endpoints on the huma API.
func RegisterRoutes(api huma.API) {
	// --- Workout groups ---

	huma.Register(api, huma.Operation{
		OperationID: "ListWorkoutGroups",
		Method:      http.MethodGet,
		Path:        "/user/workout-groups",
		Tags:        []string{"workout-groups"},
		Summary:     "List workout groups",
	}, ListWorkoutGroups)

	huma.Register(api, huma.Operation{
		OperationID:   "CreateWorkoutGroup",
		Method:        http.MethodPost,
		Path:          "/user/workout-groups",
		Tags:          []string{"workout-groups"},
		Summary:       "Create workout group",
		DefaultStatus: 201,
	}, CreateWorkoutGroup)

	huma.Register(api, huma.Operation{
		OperationID: "GetWorkoutGroup",
		Method:      http.MethodGet,
		Path:        "/user/workout-groups/{id}",
		Tags:        []string{"workout-groups"},
		Summary:     "Get workout group",
	}, GetWorkoutGroup)

	huma.Register(api, huma.Operation{
		OperationID: "UpdateWorkoutGroup",
		Method:      http.MethodPut,
		Path:        "/user/workout-groups/{id}",
		Tags:        []string{"workout-groups"},
		Summary:     "Update workout group",
	}, UpdateWorkoutGroup)

	huma.Register(api, huma.Operation{
		OperationID: "DeleteWorkoutGroup",
		Method:      http.MethodDelete,
		Path:        "/user/workout-groups/{id}",
		Tags:        []string{"workout-groups"},
		Summary:     "Delete workout group",
	}, DeleteWorkoutGroup)

	// --- Workout group memberships ---

	huma.Register(api, huma.Operation{
		OperationID: "ListWorkoutGroupMemberships",
		Method:      http.MethodGet,
		Path:        "/user/workout-group-memberships",
		Tags:        []string{"workout-group-memberships"},
		Summary:     "List workout group memberships",
	}, ListWorkoutGroupMemberships)

	huma.Register(api, huma.Operation{
		OperationID:   "CreateWorkoutGroupMembership",
		Method:        http.MethodPost,
		Path:          "/user/workout-group-memberships",
		Tags:          []string{"workout-group-memberships"},
		Summary:       "Create workout group membership",
		DefaultStatus: 201,
	}, CreateWorkoutGroupMembership)

	huma.Register(api, huma.Operation{
		OperationID: "UpdateWorkoutGroupMembership",
		Method:      http.MethodPut,
		Path:        "/user/workout-group-memberships/{id}",
		Tags:        []string{"workout-group-memberships"},
		Summary:     "Update workout group membership role",
	}, UpdateWorkoutGroupMembership)

	huma.Register(api, huma.Operation{
		OperationID: "DeleteWorkoutGroupMembership",
		Method:      http.MethodDelete,
		Path:        "/user/workout-group-memberships/{id}",
		Tags:        []string{"workout-group-memberships"},
		Summary:     "Delete workout group membership",
	}, DeleteWorkoutGroupMembership)
}
