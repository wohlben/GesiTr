package handlers

import (
	"context"

	"gesitr/internal/database"
	"gesitr/internal/humaconfig"
	"gesitr/internal/profile/models"

	"github.com/danielgtaylor/huma/v2"
)

// GetProfile returns the current user's profile.
// GET /api/profile
func GetProfile(ctx context.Context, _ *GetProfileInput) (*GetProfileOutput, error) {
	userID := humaconfig.GetUserID(ctx)

	var entity models.ProfileEntity
	if err := database.DB.Where("user_id = ?", userID).First(&entity).Error; err != nil {
		return nil, huma.Error404NotFound("profile not found")
	}
	dto := entity.ToDTO()
	return &GetProfileOutput{Body: dto}, nil
}

// CreateProfile creates a profile for the current user using the auth headers.
// POST /api/profile
func CreateProfile(ctx context.Context, _ *CreateProfileInput) (*CreateProfileOutput, error) {
	userID := humaconfig.GetUserID(ctx)
	username := humaconfig.GetUsername(ctx)
	if username == "" {
		username = userID // fallback to user ID if no username header
	}

	entity := models.ProfileEntity{
		UserID:   userID,
		Username: username,
	}
	if err := database.DB.Create(&entity).Error; err != nil {
		// Likely duplicate — profile already exists.
		return nil, huma.Error409Conflict("profile already exists")
	}
	dto := entity.ToDTO()
	return &CreateProfileOutput{Body: dto}, nil
}
