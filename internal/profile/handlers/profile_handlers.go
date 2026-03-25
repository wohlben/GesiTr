package handlers

import (
	"context"

	"gesitr/internal/database"
	"gesitr/internal/humaconfig"
	"gesitr/internal/profile/models"

	"github.com/danielgtaylor/huma/v2"
)

// GetMyProfile returns the current user's profile.
// GET /api/user/profile
//
// OpenAPI: /api/docs#/operations/get-my-profile
func GetMyProfile(ctx context.Context, input *GetMyProfileInput) (*GetMyProfileOutput, error) {
	userID := humaconfig.GetUserID(ctx)

	var entity models.UserProfileEntity
	if err := database.DB.First(&entity, "id = ?", userID).Error; err != nil {
		return nil, huma.Error404NotFound("profile not found")
	}

	return &GetMyProfileOutput{Body: entity.ToDTO()}, nil
}

// UpdateMyProfile updates the current user's profile.
// PATCH /api/user/profile
//
// OpenAPI: /api/docs#/operations/update-my-profile
func UpdateMyProfile(ctx context.Context, input *UpdateMyProfileInput) (*UpdateMyProfileOutput, error) {
	userID := humaconfig.GetUserID(ctx)

	var entity models.UserProfileEntity
	if err := database.DB.First(&entity, "id = ?", userID).Error; err != nil {
		return nil, huma.Error404NotFound("profile not found")
	}

	entity.Name = input.Body.Name
	if err := database.DB.Save(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError("failed to update profile")
	}

	return &UpdateMyProfileOutput{Body: entity.ToDTO()}, nil
}

// GetProfile returns a user's profile by ID. Public endpoint.
// GET /api/profiles/{id}
//
// OpenAPI: /api/docs#/operations/get-profile
func GetProfile(ctx context.Context, input *GetProfileInput) (*GetProfileOutput, error) {
	var entity models.UserProfileEntity
	if err := database.DB.First(&entity, "id = ?", input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("profile not found")
	}

	return &GetProfileOutput{Body: entity.ToDTO()}, nil
}
