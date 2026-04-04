package ownershipgroup

import (
	"gesitr/internal/compendium/ownershipgroup/models"

	"gorm.io/gorm"
)

// CreateGroupForEntity creates an ownership group with the given user as its owner.
// Call this inside a transaction when creating a new compendium entity.
// Returns the new group's ID.
func CreateGroupForEntity(tx *gorm.DB, ownerUserID string) (uint, error) {
	group := models.OwnershipGroupEntity{}
	if err := tx.Create(&group).Error; err != nil {
		return 0, err
	}

	membership := models.OwnershipGroupMembershipEntity{
		GroupID: group.ID,
		UserID:  ownerUserID,
		Role:    models.RoleOwner,
	}
	if err := tx.Create(&membership).Error; err != nil {
		return 0, err
	}

	return group.ID, nil
}
