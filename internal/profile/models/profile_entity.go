package models

import (
	"time"

	"gorm.io/gorm"
)

type UserProfileEntity struct {
	ID        string `gorm:"primaryKey;size:255"`
	Name      string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (UserProfileEntity) TableName() string { return "user_profiles" }

func (e *UserProfileEntity) ToDTO() UserProfile {
	return UserProfile{
		ID:        e.ID,
		Name:      e.Name,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

func UserProfileFromDTO(dto UserProfile) UserProfileEntity {
	return UserProfileEntity{
		ID:   dto.ID,
		Name: dto.Name,
	}
}
