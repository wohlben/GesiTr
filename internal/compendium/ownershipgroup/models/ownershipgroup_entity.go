package models

import (
	"gesitr/internal/shared"
)

type OwnershipGroupEntity struct {
	shared.BaseModel
	Memberships []OwnershipGroupMembershipEntity `gorm:"foreignKey:GroupID"`
}

func (OwnershipGroupEntity) TableName() string { return "ownership_groups" }

func (e *OwnershipGroupEntity) ToDTO() OwnershipGroup {
	return OwnershipGroup{
		BaseModel: e.BaseModel,
	}
}

type OwnershipGroupMembershipEntity struct {
	shared.BaseModel
	GroupID uint                  `gorm:"not null;uniqueIndex:idx_ownership_group_user"`
	Group   *OwnershipGroupEntity `gorm:"foreignKey:GroupID;constraint:OnDelete:CASCADE" json:"-"`
	UserID  string                `gorm:"not null;uniqueIndex:idx_ownership_group_user"`
	Role    OwnershipGroupRole    `gorm:"not null;default:'member'"`
}

func (OwnershipGroupMembershipEntity) TableName() string { return "ownership_group_memberships" }

func (e *OwnershipGroupMembershipEntity) ToDTO() OwnershipGroupMembership {
	return OwnershipGroupMembership{
		BaseModel: e.BaseModel,
		GroupID:   e.GroupID,
		UserID:    e.UserID,
		Role:      e.Role,
	}
}

func OwnershipGroupMembershipFromDTO(dto OwnershipGroupMembership) OwnershipGroupMembershipEntity {
	return OwnershipGroupMembershipEntity{
		BaseModel: dto.BaseModel,
		GroupID:   dto.GroupID,
		UserID:    dto.UserID,
		Role:      dto.Role,
	}
}
