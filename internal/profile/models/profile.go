package models

import "time"

// ProfileEntity stores the mapping between a Keycloak user ID and a display username.
// The UserID (keycloak sub) is the primary key — no auto-increment ID.
type ProfileEntity struct {
	UserID    string    `gorm:"primaryKey;column:user_id" json:"userId"`
	Username  string    `gorm:"not null" json:"username"`
	CreatedAt time.Time `json:"createdAt"`
}

func (ProfileEntity) TableName() string { return "profiles" }

type Profile struct {
	UserID    string    `json:"userId"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"createdAt"`
}

func (e *ProfileEntity) ToDTO() Profile {
	return Profile{
		UserID:    e.UserID,
		Username:  e.Username,
		CreatedAt: e.CreatedAt,
	}
}
