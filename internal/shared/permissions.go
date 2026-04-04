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

// AccessChecker is implemented by ownershipgroup.EntityAccess to avoid an import cycle.
type AccessChecker interface {
	CanRead() bool
	CanModify() bool
	CanDelete() bool
}

// ResolvePermissionsFromAccess determines permissions using an ownership group access check.
// Returns (permissions, visible). If visible is false, the caller should return 404.
func ResolvePermissionsFromAccess(access AccessChecker, isPublic bool) ([]Permission, bool) {
	if access.CanDelete() {
		return []Permission{PermissionRead, PermissionModify, PermissionDelete}, true
	}
	if access.CanModify() {
		return []Permission{PermissionRead, PermissionModify}, true
	}
	if access.CanRead() {
		return []Permission{PermissionRead}, true
	}
	if isPublic {
		return []Permission{PermissionRead}, true
	}
	return nil, false
}
