package handlers

import (
	"context"
	"strconv"

	"gesitr/internal/database"
	"gesitr/internal/humaconfig"
	"gesitr/internal/user/workoutlog/models"

	"github.com/danielgtaylor/huma/v2"
)

// parseUint converts a string to uint, returning 0 on failure.
func parseUint(s string) uint {
	v, _ := strconv.ParseUint(s, 10, 64)
	return uint(v)
}

// requireOwner checks that the authenticated user matches the given owner.
func requireOwner(ctx context.Context, owner string) error {
	if owner != humaconfig.GetUserID(ctx) {
		return huma.Error403Forbidden("access denied")
	}
	return nil
}

// requireLogOwner fetches the log by ID and checks that the authenticated user owns it.
func requireLogOwner(ctx context.Context, logID uint) error {
	var log models.WorkoutLogEntity
	if err := database.DB.Select("owner").First(&log, logID).Error; err != nil {
		return huma.Error404NotFound("Workout log not found")
	}
	return requireOwner(ctx, log.Owner)
}

// requireLogStatus fetches the log by ID, checks ownership, and checks its status is in the allowed set.
// Returns the log entity or an error.
func requireLogStatus(ctx context.Context, logID uint, allowed ...models.WorkoutLogStatus) (models.WorkoutLogEntity, error) {
	var log models.WorkoutLogEntity
	if err := database.DB.First(&log, logID).Error; err != nil {
		return log, huma.Error404NotFound("Workout log not found")
	}
	if err := requireOwner(ctx, log.Owner); err != nil {
		return log, err
	}
	for _, s := range allowed {
		if log.Status == s {
			return log, nil
		}
	}
	return log, huma.Error409Conflict("Workout log status is " + string(log.Status) + ", operation not allowed")
}
