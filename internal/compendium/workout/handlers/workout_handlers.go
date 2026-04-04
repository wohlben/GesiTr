package handlers

import (
	"context"
	"time"

	"gesitr/internal/compendium/ownershipgroup"
	"gesitr/internal/compendium/workout/models"
	"gesitr/internal/compendium/workoutgroup"
	"gesitr/internal/database"
	"gesitr/internal/humaconfig"
	"gesitr/internal/shared"
	logmodels "gesitr/internal/user/workoutlog/models"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func preloadWorkout(db *gorm.DB) *gorm.DB {
	return db.Preload("Sections", func(db *gorm.DB) *gorm.DB {
		return db.Order("position")
	}).Preload("Sections.Items", func(db *gorm.DB) *gorm.DB {
		return db.Order("position")
	})
}

// canReadWorkout checks if the user can read the workout via ownership group, public visibility, or workout group membership.
func canReadWorkout(userID string, entity *models.WorkoutEntity) bool {
	access := ownershipgroup.CheckAccess(database.DB, userID, entity.OwnershipGroupID)
	if access.CanRead() {
		return true
	}
	if entity.Public {
		return true
	}
	wgAccess := workoutgroup.CheckWorkoutAccess(userID, entity.ID)
	return wgAccess.CanRead()
}

// ListWorkouts returns workouts visible to the current user.
// Default: owner's workouts + all public workouts.
// logged=me: owner's workouts + workouts with workout logs.
// GET /api/workouts
func ListWorkouts(ctx context.Context, input *ListWorkoutsInput) (*ListWorkoutsOutput, error) {
	db := database.DB.Model(&models.WorkoutEntity{})
	userID := humaconfig.GetUserID(ctx)

	visibleGroups := ownershipgroup.VisibleGroupIDs(database.DB, userID)
	groupSubquery := `workouts.id IN (
		SELECT wg.workout_id FROM workout_groups wg
		JOIN workout_group_memberships wgm ON wgm.group_id = wg.id
		WHERE wgm.user_id = ? AND wgm.deleted_at IS NULL AND wg.deleted_at IS NULL)`

	if input.Logged == "me" {
		db = db.Where(`ownership_group_id IN (?) OR workouts.id IN (
			SELECT DISTINCT workout_id FROM workout_logs
			WHERE owner = ? AND workout_id IS NOT NULL AND deleted_at IS NULL)
			OR `+groupSubquery, visibleGroups, userID, userID)
	} else if input.Owner != "" {
		if input.Owner == "me" || input.Owner == userID {
			db = db.Where("ownership_group_id IN (?)", visibleGroups)
		} else {
			db = db.Where("ownership_group_id IN (SELECT group_id FROM ownership_group_memberships WHERE user_id = ? AND role = 'owner' AND deleted_at IS NULL) AND public = ?", input.Owner, true)
		}
	} else {
		db = db.Where("ownership_group_id IN (?) OR public = ? OR "+groupSubquery, visibleGroups, true, userID)
	}

	if input.Public == "true" {
		db = db.Where("public = ?", true)
	}

	if input.Q != "" {
		pattern := "%" + input.Q + "%"
		db = db.Where("name LIKE ?", pattern)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	p := input.ToPaginationParams()
	var entities []models.WorkoutEntity
	if err := preloadWorkout(shared.ApplyPagination(db, p)).Find(&entities).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	var nonOwnedIDs []uint
	for _, e := range entities {
		access := ownershipgroup.CheckAccess(database.DB, userID, e.OwnershipGroupID)
		if !access.IsOwner {
			nonOwnedIDs = append(nonOwnedIDs, e.ID)
		}
	}
	groupInfoMap := workoutgroup.GroupInfoForWorkouts(userID, nonOwnedIDs)

	dtos := make([]models.Workout, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
		if info, ok := groupInfoMap[entities[i].ID]; ok {
			dtos[i].WorkoutGroup = &models.WorkoutGroupInfo{
				GroupName:  info.GroupName,
				Membership: info.MembershipRole,
			}
		}
	}
	return &ListWorkoutsOutput{Body: humaconfig.PaginatedBody[models.Workout]{
		Items: dtos, Total: total, Limit: p.Limit, Offset: p.Offset,
	}}, nil
}

// CreateWorkout creates a workout. If SourceWorkoutID is provided, creates
// forked+equivalent relationships (fork mechanic).
// POST /api/workouts
func CreateWorkout(ctx context.Context, input *CreateWorkoutInput) (*CreateWorkoutOutput, error) {
	entity := models.WorkoutEntity{
		Name:   input.Body.Name,
		Notes:  input.Body.Notes,
		Public: input.Body.Public,
	}
	userID := humaconfig.GetUserID(ctx)

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&entity).Error; err != nil {
			return err
		}

		groupID, err := ownershipgroup.CreateGroupForEntity(tx, userID)
		if err != nil {
			return err
		}
		entity.OwnershipGroupID = groupID
		if err := tx.Model(&entity).Update("ownership_group_id", groupID).Error; err != nil {
			return err
		}

		dto := entity.ToDTO()
		if err := tx.Create(&models.WorkoutHistoryEntity{
			WorkoutID: entity.ID,
			Version:   dto.Version,
			Snapshot:  shared.SnapshotJSON(dto),
			ChangedAt: time.Now(),
			ChangedBy: userID,
		}).Error; err != nil {
			return err
		}

		if input.Body.SourceWorkoutID != nil {
			sourceID := *input.Body.SourceWorkoutID
			forked := models.WorkoutRelationshipEntity{
				RelationshipType: models.WorkoutRelationshipTypeForked,
				Strength:         1.0,
				OwnershipGroupID: groupID,
				FromWorkoutID:    entity.ID,
				ToWorkoutID:      sourceID,
			}
			if err := tx.Create(&forked).Error; err != nil {
				return err
			}
			equivalent := models.WorkoutRelationshipEntity{
				RelationshipType: models.WorkoutRelationshipTypeEquivalent,
				Strength:         1.0,
				OwnershipGroupID: groupID,
				FromWorkoutID:    entity.ID,
				ToWorkoutID:      sourceID,
			}
			if err := tx.Create(&equivalent).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	return &CreateWorkoutOutput{Body: entity.ToDTO()}, nil
}

// GetWorkout returns a workout with its full section tree.
// Public workouts are readable by anyone; private ones require ownership or group membership.
// GET /api/workouts/{id}
func GetWorkout(ctx context.Context, input *GetWorkoutInput) (*GetWorkoutOutput, error) {
	var entity models.WorkoutEntity
	if err := preloadWorkout(database.DB).First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Workout not found")
	}
	userID := humaconfig.GetUserID(ctx)
	if !canReadWorkout(userID, &entity) {
		return nil, huma.Error403Forbidden("access denied")
	}
	dto := entity.ToDTO()
	ogAccess := ownershipgroup.CheckAccess(database.DB, userID, entity.OwnershipGroupID)
	if !ogAccess.IsOwner {
		wgAccess := workoutgroup.CheckWorkoutAccess(userID, entity.ID)
		if wgAccess.GroupName != "" {
			dto.WorkoutGroup = &models.WorkoutGroupInfo{
				GroupName:  wgAccess.GroupName,
				Membership: wgAccess.MembershipRole,
			}
		}
	}
	return &GetWorkoutOutput{Body: dto}, nil
}

// UpdateWorkout updates workout metadata. Increments version and creates a history snapshot.
// PUT /api/workouts/{id}
func UpdateWorkout(ctx context.Context, input *UpdateWorkoutInput) (*UpdateWorkoutOutput, error) {
	var existing models.WorkoutEntity
	if err := preloadWorkout(database.DB).First(&existing, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Workout not found")
	}
	userID := humaconfig.GetUserID(ctx)
	ogAccess := ownershipgroup.CheckAccess(database.DB, userID, existing.OwnershipGroupID)
	wgAccess := workoutgroup.CheckWorkoutAccess(userID, existing.ID)
	if !ogAccess.CanModify() && !wgAccess.CanModify() {
		return nil, huma.Error403Forbidden("access denied")
	}

	oldDTO := existing.ToDTO()

	existing.Name = input.Body.Name
	existing.Notes = input.Body.Notes
	if ogAccess.CanModify() {
		existing.Public = input.Body.Public
	}

	newDTO := existing.ToDTO()
	if !models.WorkoutChanged(oldDTO, newDTO) {
		return &UpdateWorkoutOutput{Body: newDTO}, nil
	}

	existing.Version = oldDTO.Version + 1

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&existing).Error; err != nil {
			return err
		}

		updatedDTO := existing.ToDTO()
		if err := tx.Create(&models.WorkoutHistoryEntity{
			WorkoutID: existing.ID,
			Version:   updatedDTO.Version,
			Snapshot:  shared.SnapshotJSON(updatedDTO),
			ChangedAt: time.Now(),
			ChangedBy: userID,
		}).Error; err != nil {
			return err
		}
		// Invalidate stale planning logs so /start creates fresh ones.
		return tx.Where("workout_id = ? AND status = ?", existing.ID, logmodels.WorkoutLogStatusPlanning).
			Delete(&logmodels.WorkoutLogEntity{}).Error
	})
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	if err := preloadWorkout(database.DB).First(&existing, existing.ID).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	dto := existing.ToDTO()
	if !ogAccess.IsOwner && wgAccess.GroupName != "" {
		dto.WorkoutGroup = &models.WorkoutGroupInfo{
			GroupName:  wgAccess.GroupName,
			Membership: wgAccess.MembershipRole,
		}
	}
	return &UpdateWorkoutOutput{Body: dto}, nil
}

// DeleteWorkout deletes a workout. Owner only.
// DELETE /api/workouts/{id}
func DeleteWorkout(ctx context.Context, input *DeleteWorkoutInput) (*DeleteWorkoutOutput, error) {
	var entity models.WorkoutEntity
	if err := database.DB.First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Workout not found")
	}
	access := ownershipgroup.CheckAccess(database.DB, humaconfig.GetUserID(ctx), entity.OwnershipGroupID)
	if !access.CanDelete() {
		return nil, huma.Error403Forbidden("access denied")
	}
	if err := database.DB.Delete(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return nil, nil
}

// GetWorkoutPermissions returns the current user's permissions on a workout.
// GET /api/workouts/{id}/permissions
func GetWorkoutPermissions(ctx context.Context, input *GetWorkoutPermissionsInput) (*GetWorkoutPermissionsOutput, error) {
	var entity models.WorkoutEntity
	if err := database.DB.First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Workout not found")
	}
	userID := humaconfig.GetUserID(ctx)

	// Check ownership group permissions first.
	access := ownershipgroup.CheckAccess(database.DB, userID, entity.OwnershipGroupID)
	perms, _ := shared.ResolvePermissionsFromAccess(access, entity.Public)

	// Fall back to workout-group-based access.
	if len(perms) == 0 {
		wgAccess := workoutgroup.CheckWorkoutAccess(userID, entity.ID)
		if wgAccess.CanDelete() {
			perms = []shared.Permission{shared.PermissionRead, shared.PermissionModify, shared.PermissionDelete}
		} else if wgAccess.CanModify() {
			perms = []shared.Permission{shared.PermissionRead, shared.PermissionModify}
		} else if wgAccess.CanRead() {
			perms = []shared.Permission{shared.PermissionRead}
		}
	}

	if len(perms) == 0 {
		return nil, huma.Error404NotFound("Workout not found")
	}

	return &GetWorkoutPermissionsOutput{Body: shared.PermissionsResponse{Permissions: perms}}, nil
}

// ListWorkoutVersions returns all version snapshots for a workout.
// GET /api/workouts/{id}/versions
func ListWorkoutVersions(ctx context.Context, input *ListWorkoutVersionsInput) (*ListWorkoutVersionsOutput, error) {
	var entity models.WorkoutEntity
	if err := database.DB.First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Workout not found")
	}
	userID := humaconfig.GetUserID(ctx)
	if !canReadWorkout(userID, &entity) {
		return nil, huma.Error403Forbidden("access denied")
	}

	var history []models.WorkoutHistoryEntity
	if err := database.DB.Where("workout_id = ?", input.ID).Order("version ASC").Find(&history).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	entries := make([]shared.VersionEntry, len(history))
	for i := range history {
		entries[i] = history[i].ToVersionEntry()
	}
	return &ListWorkoutVersionsOutput{Body: entries}, nil
}

// GetWorkoutVersion returns a specific version snapshot.
// GET /api/workouts/{id}/versions/{version}
func GetWorkoutVersion(ctx context.Context, input *GetWorkoutVersionInput) (*GetWorkoutVersionOutput, error) {
	var entity models.WorkoutEntity
	if err := database.DB.First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Workout not found")
	}
	userID := humaconfig.GetUserID(ctx)
	if !canReadWorkout(userID, &entity) {
		return nil, huma.Error403Forbidden("access denied")
	}

	var h models.WorkoutHistoryEntity
	if err := database.DB.Where("workout_id = ? AND version = ?", input.ID, input.Version).First(&h).Error; err != nil {
		return nil, huma.Error404NotFound("Version not found")
	}
	return &GetWorkoutVersionOutput{Body: h.ToVersionEntry()}, nil
}
