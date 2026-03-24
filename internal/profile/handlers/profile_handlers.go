package handlers

import (
	"context"
	"encoding/json"

	"gesitr/internal/database"
	"gesitr/internal/humaconfig"
	"gesitr/internal/profile/models"

	"github.com/danielgtaylor/huma/v2"
)

// GetMyProfile returns the current user's profile.
// GET /api/user/profile
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
func UpdateMyProfile(ctx context.Context, input *UpdateMyProfileInput) (*UpdateMyProfileOutput, error) {
	userID := humaconfig.GetUserID(ctx)

	var req models.UpdateProfileRequest
	if err := json.Unmarshal(input.RawBody, &req); err != nil {
		return nil, huma.Error400BadRequest(err.Error())
	}
	if req.Name == "" {
		return nil, huma.Error400BadRequest("name is required")
	}

	var entity models.UserProfileEntity
	if err := database.DB.First(&entity, "id = ?", userID).Error; err != nil {
		return nil, huma.Error404NotFound("profile not found")
	}

	entity.Name = req.Name
	if err := database.DB.Save(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError("failed to update profile")
	}

	return &UpdateMyProfileOutput{Body: entity.ToDTO()}, nil
}

// GetProfile returns a user's profile by ID. Public endpoint.
// GET /api/profiles/{id}
func GetProfile(ctx context.Context, input *GetProfileInput) (*GetProfileOutput, error) {
	var entity models.UserProfileEntity
	if err := database.DB.First(&entity, "id = ?", input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("profile not found")
	}

	return &GetProfileOutput{Body: entity.ToDTO()}, nil
}
