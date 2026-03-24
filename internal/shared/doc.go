// Package shared provides cross-cutting types and utilities used by all
// domain packages in the GesiTr API.
//
// # Permissions
//
// The [ResolvePermissions] function implements the core authorization logic.
// Every resource in GesiTr has an owner and a public/private visibility flag.
// Permissions are resolved as follows:
//
//   - Owner: receives [PermissionRead], [PermissionModify], and [PermissionDelete]
//   - Non-owner viewing a public resource: receives [PermissionRead] only
//   - Non-owner viewing a private resource: receives no permissions (empty list)
//
// Handler packages call ResolvePermissions to populate [PermissionsResponse] on
// dedicated /permissions endpoints (exercises, equipment, exercise groups).
// User-scoped resources like workouts enforce ownership directly with 403 responses
// rather than exposing a permissions endpoint.
//
// # Pagination
//
// [ParsePagination] and [ApplyPagination] standardize list endpoint pagination.
//
// # History
//
// [SnapshotJSON] and [VersionEntry] support versioned history for exercises
// and equipment.
//
// # Base Model
//
// [BaseModel] provides the standard GORM fields (ID, CreatedAt, UpdatedAt,
// DeletedAt) used by all entity types.
package shared
