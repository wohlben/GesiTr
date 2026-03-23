package shared

// Permission represents an allowed action on a resource.
type Permission string

const (
	PermissionRead   Permission = "READ"
	PermissionModify Permission = "MODIFY"
	PermissionDelete Permission = "DELETE"
)

// PermissionsResponse is the API response for permission queries.
type PermissionsResponse struct {
	Permissions []Permission `json:"permissions"`
}

// ResolvePermissions determines the permissions for a user on a resource.
// Returns (permissions, visible). If visible is false, the caller should return 404.
func ResolvePermissions(userID, entityOwner string, isPublic bool) ([]Permission, bool) {
	if userID == entityOwner {
		return []Permission{PermissionRead, PermissionModify, PermissionDelete}, true
	}
	if isPublic {
		return []Permission{PermissionRead}, true
	}
	return nil, false
}
