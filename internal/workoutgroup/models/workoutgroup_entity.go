package models

import (
	profilemodels "gesitr/internal/profile/models"
	"gesitr/internal/shared"
	workoutmodels "gesitr/internal/workout/models"
)

type WorkoutGroupEntity struct {
	shared.BaseModel
	Name         string                           `gorm:"not null"`
	WorkoutID    uint                             `gorm:"not null;index"`
	Workout      *workoutmodels.WorkoutEntity     `gorm:"foreignKey:WorkoutID;constraint:OnDelete:CASCADE" json:"-"`
	Owner        string                           `gorm:"not null;index"`
	OwnerProfile *profilemodels.UserProfileEntity `gorm:"foreignKey:Owner;references:ID;constraint:OnDelete:RESTRICT" json:"-"`
	Memberships  []WorkoutGroupMembershipEntity   `gorm:"foreignKey:GroupID"`
}

func (WorkoutGroupEntity) TableName() string { return "workout_groups" }

func (e *WorkoutGroupEntity) ToDTO() WorkoutGroup {
	return WorkoutGroup{
		BaseModel: e.BaseModel,
		Name:      e.Name,
		WorkoutID: e.WorkoutID,
		Owner:     e.Owner,
	}
}

func WorkoutGroupFromDTO(dto WorkoutGroup) WorkoutGroupEntity {
	return WorkoutGroupEntity{
		BaseModel: dto.BaseModel,
		Name:      dto.Name,
		WorkoutID: dto.WorkoutID,
		Owner:     dto.Owner,
	}
}

type WorkoutGroupMembershipEntity struct {
	shared.BaseModel
	GroupID     uint                             `gorm:"not null;uniqueIndex:idx_group_user"`
	Group       *WorkoutGroupEntity              `gorm:"foreignKey:GroupID;constraint:OnDelete:CASCADE" json:"-"`
	UserID      string                           `gorm:"not null;uniqueIndex:idx_group_user"`
	UserProfile *profilemodels.UserProfileEntity `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:RESTRICT" json:"-"`
	Role        WorkoutGroupRole                 `gorm:"not null;default:'member'"`
}

func (WorkoutGroupMembershipEntity) TableName() string { return "workout_group_memberships" }

func (e *WorkoutGroupMembershipEntity) ToDTO() WorkoutGroupMembership {
	return WorkoutGroupMembership{
		BaseModel: e.BaseModel,
		GroupID:   e.GroupID,
		UserID:    e.UserID,
		Role:      e.Role,
	}
}

func WorkoutGroupMembershipFromDTO(dto WorkoutGroupMembership) WorkoutGroupMembershipEntity {
	return WorkoutGroupMembershipEntity{
		BaseModel: dto.BaseModel,
		GroupID:   dto.GroupID,
		UserID:    dto.UserID,
		Role:      dto.Role,
	}
}
