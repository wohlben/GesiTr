package models

import "gesitr/internal/shared"

// OwnershipGroupRole defines the roles within an ownership group.
type OwnershipGroupRole string

const (
	RoleOwner  OwnershipGroupRole = "owner"
	RoleAdmin  OwnershipGroupRole = "admin"  // reserved — no extra privileges yet
	RoleMember OwnershipGroupRole = "member" // read access to private entities
)

// OwnershipGroup is a lightweight entity that holds memberships for a compendium entity.
type OwnershipGroup struct {
	shared.BaseModel `tstype:",extends"`
}

// OwnershipGroupMembership links a user to an ownership group with a role.
type OwnershipGroupMembership struct {
	shared.BaseModel `tstype:",extends"`
	GroupID          uint               `json:"groupId"`
	UserID           string             `json:"userId"`
	Role             OwnershipGroupRole `json:"role"`
}
