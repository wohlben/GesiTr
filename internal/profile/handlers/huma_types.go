package handlers

import (
	"gesitr/internal/profile/models"
)

// --- GetMyProfile ---

type GetMyProfileInput struct{}

type GetMyProfileOutput struct {
	Body models.UserProfile
}

// --- UpdateMyProfile ---

type UpdateMyProfileBody struct {
	Name string `json:"name" required:"true"`
}

type UpdateMyProfileInput struct {
	Body UpdateMyProfileBody
}

type UpdateMyProfileOutput struct {
	Body models.UserProfile
}

// --- GetProfile ---

type GetProfileInput struct {
	ID string `path:"id"`
}

type GetProfileOutput struct {
	Body models.UserProfile
}
