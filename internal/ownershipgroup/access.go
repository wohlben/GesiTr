package ownershipgroup

import (
	"gesitr/internal/ownershipgroup/models"

	"gorm.io/gorm"
)

// EntityAccess describes a user's access level for an entity via its ownership group.
type EntityAccess struct {
	IsOwner  bool
	IsMember bool
	Role     models.OwnershipGroupRole
}

// CheckAccess determines a user's access level for a given ownership group.
func CheckAccess(db *gorm.DB, userID string, ownershipGroupID uint) EntityAccess {
	if ownershipGroupID == 0 {
		return EntityAccess{}
	}

	var membership models.OwnershipGroupMembershipEntity
	err := db.
		Where("group_id = ? AND user_id = ?", ownershipGroupID, userID).
		First(&membership).Error
	if err != nil {
		return EntityAccess{}
	}

	return EntityAccess{
		IsOwner:  membership.Role == models.RoleOwner,
		IsMember: true,
		Role:     membership.Role,
	}
}

// CheckAccessBatch determines access levels for multiple ownership groups at once.
func CheckAccessBatch(db *gorm.DB, userID string, groupIDs []uint) map[uint]EntityAccess {
	if len(groupIDs) == 0 {
		return nil
	}

	var memberships []models.OwnershipGroupMembershipEntity
	db.Where("group_id IN ? AND user_id = ?", groupIDs, userID).Find(&memberships)

	m := make(map[uint]EntityAccess, len(memberships))
	for _, mem := range memberships {
		m[mem.GroupID] = EntityAccess{
			IsOwner:  mem.Role == models.RoleOwner,
			IsMember: true,
			Role:     mem.Role,
		}
	}
	return m
}

// VisibleGroupIDs returns a GORM subquery selecting group IDs the user is a member of.
// Use with: db.Where("ownership_group_id IN (?)", VisibleGroupIDs(db, userID))
func VisibleGroupIDs(db *gorm.DB, userID string) *gorm.DB {
	return db.Table("ownership_group_memberships").
		Select("group_id").
		Where("user_id = ? AND deleted_at IS NULL", userID)
}

func (a EntityAccess) CanRead() bool {
	return a.IsOwner || a.IsMember
}

func (a EntityAccess) CanModify() bool {
	return a.IsOwner
}

func (a EntityAccess) CanDelete() bool {
	return a.IsOwner
}
