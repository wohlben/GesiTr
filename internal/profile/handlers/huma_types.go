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

type UpdateMyProfileInput struct {
	RawBody []byte
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
