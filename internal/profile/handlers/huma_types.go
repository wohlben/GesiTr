package handlers

import "gesitr/internal/profile/models"

type GetProfileInput struct{}

type GetProfileOutput struct {
	Body models.Profile
}

type CreateProfileInput struct{}

type CreateProfileOutput struct {
	Body models.Profile
}
