package handlers

import "gesitr/internal/compendium/ownershipgroup/models"

// --- Ownership group membership types ---

type ListOwnershipGroupMembershipsInput struct {
	GroupID uint `path:"id"`
}

type ListOwnershipGroupMembershipsOutput struct {
	Body []models.OwnershipGroupMembership
}

type CreateOwnershipGroupMembershipInput struct {
	GroupID uint `path:"id"`
	Body    struct {
		UserID string `json:"userId" required:"true"`
	}
}

type CreateOwnershipGroupMembershipOutput struct {
	Body models.OwnershipGroupMembership
}

type UpdateOwnershipGroupMembershipInput struct {
	ID   uint `path:"id"`
	Body struct {
		Role models.OwnershipGroupRole `json:"role" required:"true"`
	}
}

type UpdateOwnershipGroupMembershipOutput struct {
	Body models.OwnershipGroupMembership
}

type DeleteOwnershipGroupMembershipInput struct {
	ID uint `path:"id"`
}

type DeleteOwnershipGroupMembershipOutput struct{}
