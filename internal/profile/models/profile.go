package models

import "time"

type UserProfile struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UpdateProfileRequest struct {
	Name string `json:"name" binding:"required"`
}
