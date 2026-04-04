package handlers

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// RegisterRoutes registers all ownership group membership endpoints on the huma API.
func RegisterRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "ListOwnershipGroupMemberships",
		Method:      http.MethodGet,
		Path:        "/ownership-groups/{id}/memberships",
		Tags:        []string{"ownership-groups"},
		Summary:     "List ownership group memberships",
	}, ListOwnershipGroupMemberships)

	huma.Register(api, huma.Operation{
		OperationID:   "CreateOwnershipGroupMembership",
		Method:        http.MethodPost,
		Path:          "/ownership-groups/{id}/memberships",
		Tags:          []string{"ownership-groups"},
		Summary:       "Add member to ownership group",
		DefaultStatus: 201,
	}, CreateOwnershipGroupMembership)

	huma.Register(api, huma.Operation{
		OperationID: "UpdateOwnershipGroupMembership",
		Method:      http.MethodPut,
		Path:        "/ownership-group-memberships/{id}",
		Tags:        []string{"ownership-groups"},
		Summary:     "Update ownership group membership role",
	}, UpdateOwnershipGroupMembership)

	huma.Register(api, huma.Operation{
		OperationID: "DeleteOwnershipGroupMembership",
		Method:      http.MethodDelete,
		Path:        "/ownership-group-memberships/{id}",
		Tags:        []string{"ownership-groups"},
		Summary:     "Remove member from ownership group",
	}, DeleteOwnershipGroupMembership)
}
