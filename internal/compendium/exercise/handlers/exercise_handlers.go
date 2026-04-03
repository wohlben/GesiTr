package handlers

import (
	"context"
	"encoding/json"
	"time"

	"gesitr/internal/compendium/exercise/models"
	"gesitr/internal/database"
	"gesitr/internal/humaconfig"
	"gesitr/internal/shared"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func exerciseDTOFromBody(b ExerciseBody) models.Exercise {
	nameDTOs := make([]models.ExerciseNameDTO, len(b.Names))
	for i, n := range b.Names {
		nameDTOs[i] = models.ExerciseNameDTO{Name: n}
	}
	return models.Exercise{
		Names:                         nameDTOs,
		Type:                          b.Type,
		Force:                         b.Force,
		PrimaryMuscles:                b.PrimaryMuscles,
		SecondaryMuscles:              b.SecondaryMuscles,
		TechnicalDifficulty:           b.TechnicalDifficulty,
		BodyWeightScaling:             b.BodyWeightScaling,
		SuggestedMeasurementParadigms: b.SuggestedMeasurementParadigms,
		Description:                   b.Description,
		Instructions:                  b.Instructions,
		Images:                        b.Images,
		AuthorName:                    b.AuthorName,
		AuthorUrl:                     b.AuthorUrl,
		Public:                        b.Public,
		ParentExerciseID:              b.ParentExerciseID,
		EquipmentIDs:                  b.EquipmentIDs,
	}
}

var exercisePreloads = []string{
	"Forces", "Muscles", "Paradigms", "Instructions", "Images", "Equipment",
}

func preloadExercise(db *gorm.DB) *gorm.DB {
	for _, p := range exercisePreloads {
		db = db.Preload(p)
	}
	// Sort names by popularity (most users' preferred name first), then position as tiebreaker.
	db = db.Preload("Names", func(tx *gorm.DB) *gorm.DB {
		return tx.Order("(SELECT COUNT(*) FROM exercise_name_preferences WHERE exercise_name_id = exercise_names.id) DESC, exercise_names.position ASC")
	})
	return db
}

// ListExercises returns exercises visible to the current user: their own
// exercises plus all public exercises. Filter by owner or public query params.
// GET /api/exercises
//
// OpenAPI: /api/docs#/operations/ListExercises
func ListExercises(ctx context.Context, input *ListExercisesInput) (*ListExercisesOutput, error) {
	db := database.DB.Model(&models.ExerciseEntity{})

	userID := humaconfig.GetUserID(ctx)
	if input.Mastery == "me" {
		db = db.Where("owner = ? OR exercises.id IN (SELECT exercise_id FROM mastery_experience WHERE owner = ?)", userID, userID)
	} else if input.Owner != "" {
		if input.Owner == "me" || input.Owner == userID {
			db = db.Where("owner = ?", userID)
		} else {
			db = db.Where("owner = ? AND public = ?", input.Owner, true)
		}
	} else {
		db = db.Where("owner = ? OR public = ?", userID, true)
	}
	if input.Public == "true" {
		db = db.Where("public = ?", true)
	}

	if input.Q != "" {
		pattern := "%" + input.Q + "%"
		db = db.Where(
			"exercises.id IN (SELECT exercise_id FROM exercise_names WHERE name LIKE ?)",
			pattern,
		)
	}
	if input.Type != "" {
		db = db.Where("exercises.type = ?", input.Type)
	}
	if input.Difficulty != "" {
		db = db.Where("exercises.technical_difficulty = ?", input.Difficulty)
	}
	if input.Force != "" {
		db = db.Where("exercises.id IN (SELECT exercise_id FROM exercise_forces WHERE force = ?)", input.Force)
	}
	if input.Muscle != "" {
		db = db.Where("exercises.id IN (SELECT exercise_id FROM exercise_muscles WHERE muscle = ?)", input.Muscle)
	}
	if input.PrimaryMuscle != "" {
		db = db.Where("exercises.id IN (SELECT exercise_id FROM exercise_muscles WHERE muscle = ? AND is_primary = true)", input.PrimaryMuscle)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	p := input.ToPaginationParams()
	var entities []models.ExerciseEntity
	if err := preloadExercise(shared.ApplyPagination(db, p)).Find(&entities).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	dtos := make([]models.Exercise, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
	}
	return &ListExercisesOutput{Body: humaconfig.PaginatedBody[models.Exercise]{
		Items: dtos, Total: total, Limit: p.Limit, Offset: p.Offset,
	}}, nil
}

// CreateExercise creates an exercise owned by the current user. The exercise
// can reference equipment via equipmentIds — equipment must already exist
// (see [gesitr/internal/compendium/equipment/handlers.CreateEquipment]). To use this
// exercise in a workout, create an exercise scheme via [CreateExerciseScheme].
// POST /api/exercises
//
// OpenAPI: /api/docs#/operations/CreateExercise
func CreateExercise(ctx context.Context, input *CreateExerciseInput) (*CreateExerciseOutput, error) {
	if len(input.Body.Names) == 0 {
		return nil, huma.Error422UnprocessableEntity("at least one name is required")
	}

	dto := exerciseDTOFromBody(input.Body)

	entity := models.ExerciseFromDTO(dto)
	entity.Owner = humaconfig.GetUserID(ctx)
	entity.Public = dto.Public

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&entity).Error; err != nil {
			return err
		}

		if err := preloadExercise(tx).First(&entity, entity.ID).Error; err != nil {
			return err
		}

		resultDTO := entity.ToDTO()
		if err := tx.Create(&models.ExerciseHistoryEntity{
			ExerciseID: entity.ID,
			Version:    resultDTO.Version,
			Snapshot:   shared.SnapshotJSON(resultDTO),
			ChangedAt:  time.Now(),
			ChangedBy:  resultDTO.Owner,
		}).Error; err != nil {
			return err
		}

		if input.Body.SourceExerciseID != nil {
			sourceID := *input.Body.SourceExerciseID
			forked := models.ExerciseRelationshipEntity{
				RelationshipType: models.ExerciseRelationshipTypeForked,
				Strength:         1.0,
				Owner:            entity.Owner,
				FromExerciseID:   entity.ID,
				ToExerciseID:     sourceID,
			}
			if err := tx.Create(&forked).Error; err != nil {
				return err
			}
			equivalent := models.ExerciseRelationshipEntity{
				RelationshipType: models.ExerciseRelationshipTypeEquivalent,
				Strength:         1.0,
				Owner:            entity.Owner,
				FromExerciseID:   entity.ID,
				ToExerciseID:     sourceID,
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

	return &CreateExerciseOutput{Body: entity.ToDTO()}, nil
}

// GetExercisePermissions returns the current user's permissions on an exercise.
// See [gesitr/internal/shared.ResolvePermissions] for the permission model.
// GET /api/exercises/:id/permissions
//
// OpenAPI: /api/docs#/operations/GetExercisePermissions
func GetExercisePermissions(ctx context.Context, input *GetExercisePermissionsInput) (*GetExercisePermissionsOutput, error) {
	var entity models.ExerciseEntity
	if err := database.DB.First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Exercise not found")
	}
	userID := humaconfig.GetUserID(ctx)
	perms, _ := shared.ResolvePermissions(userID, entity.Owner, entity.Public)
	if perms == nil {
		perms = []shared.Permission{}
	}
	return &GetExercisePermissionsOutput{Body: shared.PermissionsResponse{Permissions: perms}}, nil
}

// GetExercise returns a single exercise. Public exercises are visible to all
// users; private exercises are visible only to their owner.
// GET /api/exercises/:id
//
// OpenAPI: /api/docs#/operations/GetExercise
func GetExercise(ctx context.Context, input *GetExerciseInput) (*GetExerciseOutput, error) {
	var entity models.ExerciseEntity
	if err := preloadExercise(database.DB).First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Exercise not found")
	}
	userID := humaconfig.GetUserID(ctx)
	perms, _ := shared.ResolvePermissions(userID, entity.Owner, entity.Public)
	if len(perms) == 0 {
		return nil, huma.Error403Forbidden("access denied")
	}
	return &GetExerciseOutput{Body: entity.ToDTO()}, nil
}

// UpdateExercise updates an exercise. Creates a version history entry.
// Owner only — returns 403 for non-owners. PUT /api/exercises/:id
//
// OpenAPI: /api/docs#/operations/UpdateExercise
func UpdateExercise(ctx context.Context, input *UpdateExerciseInput) (*UpdateExerciseOutput, error) {
	if len(input.Body.Names) == 0 {
		return nil, huma.Error422UnprocessableEntity("at least one name is required")
	}

	var existing models.ExerciseEntity
	if err := preloadExercise(database.DB).First(&existing, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Exercise not found")
	}

	if existing.Owner != humaconfig.GetUserID(ctx) {
		return nil, huma.Error403Forbidden("access denied")
	}

	dto := exerciseDTOFromBody(input.Body)
	dto.Owner = existing.Owner
	oldDTO := existing.ToDTO()

	if !models.ExerciseChanged(oldDTO, dto) {
		return &UpdateExerciseOutput{Body: oldDTO}, nil
	}

	entity := models.ExerciseFromDTO(dto)
	entity.ID = existing.ID
	entity.Owner = existing.Owner
	entity.Version = existing.Version + 1

	forces := entity.Forces
	muscles := entity.Muscles
	paradigms := entity.Paradigms
	instructions := entity.Instructions
	images := entity.Images
	names := entity.Names
	equipment := entity.Equipment
	entity.Forces = nil
	entity.Muscles = nil
	entity.Paradigms = nil
	entity.Instructions = nil
	entity.Images = nil
	entity.Names = nil
	entity.Equipment = nil

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("exercise_id = ?", entity.ID).Delete(&models.ExerciseForce{}).Error; err != nil {
			return err
		}
		if err := tx.Where("exercise_id = ?", entity.ID).Delete(&models.ExerciseMuscle{}).Error; err != nil {
			return err
		}
		if err := tx.Where("exercise_id = ?", entity.ID).Delete(&models.ExerciseMeasurementParadigm{}).Error; err != nil {
			return err
		}
		if err := tx.Where("exercise_id = ?", entity.ID).Delete(&models.ExerciseInstruction{}).Error; err != nil {
			return err
		}
		if err := tx.Where("exercise_id = ?", entity.ID).Delete(&models.ExerciseImage{}).Error; err != nil {
			return err
		}
		if err := tx.Where("exercise_id = ?", entity.ID).Delete(&models.ExerciseName{}).Error; err != nil {
			return err
		}
		if err := tx.Where("exercise_id = ?", entity.ID).Delete(&models.ExerciseEquipment{}).Error; err != nil {
			return err
		}

		if err := tx.Save(&entity).Error; err != nil {
			return err
		}

		for i := range forces {
			forces[i].ExerciseID = entity.ID
			if err := tx.Create(&forces[i]).Error; err != nil {
				return err
			}
		}
		for i := range muscles {
			muscles[i].ExerciseID = entity.ID
			if err := tx.Create(&muscles[i]).Error; err != nil {
				return err
			}
		}
		for i := range paradigms {
			paradigms[i].ExerciseID = entity.ID
			if err := tx.Create(&paradigms[i]).Error; err != nil {
				return err
			}
		}
		for i := range instructions {
			instructions[i].ExerciseID = entity.ID
			if err := tx.Create(&instructions[i]).Error; err != nil {
				return err
			}
		}
		for i := range images {
			images[i].ExerciseID = entity.ID
			if err := tx.Create(&images[i]).Error; err != nil {
				return err
			}
		}
		for i := range names {
			names[i].ExerciseID = entity.ID
			if err := tx.Create(&names[i]).Error; err != nil {
				return err
			}
		}
		for i := range equipment {
			equipment[i].ExerciseID = entity.ID
			if err := tx.Create(&equipment[i]).Error; err != nil {
				return err
			}
		}

		entity.Forces = forces
		entity.Muscles = muscles
		entity.Paradigms = paradigms
		entity.Instructions = instructions
		entity.Images = images
		entity.Names = names
		entity.Equipment = equipment

		resultDTO := entity.ToDTO()
		if err := tx.Create(&models.ExerciseHistoryEntity{
			ExerciseID: entity.ID,
			Version:    resultDTO.Version,
			Snapshot:   shared.SnapshotJSON(resultDTO),
			ChangedAt:  time.Now(),
			ChangedBy:  resultDTO.Owner,
		}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	if err := preloadExercise(database.DB).First(&entity, entity.ID).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &UpdateExerciseOutput{Body: entity.ToDTO()}, nil
}

// ListExerciseVersions returns the version history for an exercise. Each
// update via [UpdateExercise] creates a new version entry.
// GET /api/exercises/:id/versions
//
// OpenAPI: /api/docs#/operations/ListExerciseVersions
func ListExerciseVersions(ctx context.Context, input *ListExerciseVersionsInput) (*ListExerciseVersionsOutput, error) {
	var entity models.ExerciseEntity
	if err := database.DB.First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Exercise not found")
	}

	var history []models.ExerciseHistoryEntity
	if err := database.DB.Where("exercise_id = ?", entity.ID).Order("version ASC").Find(&history).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	entries := make([]shared.VersionEntry, len(history))
	for i := range history {
		entries[i] = history[i].ToVersionEntry()
	}
	return &ListExerciseVersionsOutput{Body: entries}, nil
}

// GetExerciseVersion returns a specific historical version of an exercise
// by templateId and version number.
// GET /api/exercises/templates/:templateId/versions/:version
//
// OpenAPI: /api/docs#/operations/GetExerciseVersion
func GetExerciseVersion(ctx context.Context, input *GetExerciseVersionInput) (*GetExerciseVersionOutput, error) {
	var history models.ExerciseHistoryEntity
	if err := database.DB.Where("exercise_id = ? AND version = ?", input.ID, input.Version).First(&history).Error; err != nil {
		return nil, huma.Error404NotFound("Version not found")
	}

	var snap models.Exercise
	json.Unmarshal([]byte(history.Snapshot), &snap)
	userID := humaconfig.GetUserID(ctx)
	perms, _ := shared.ResolvePermissions(userID, snap.Owner, snap.Public)
	if len(perms) == 0 {
		return nil, huma.Error403Forbidden("access denied")
	}

	return &GetExerciseVersionOutput{Body: history.ToVersionEntry()}, nil
}

// DeleteExercise deletes an exercise. Owner only.
// DELETE /api/exercises/:id
//
// OpenAPI: /api/docs#/operations/DeleteExercise
func DeleteExercise(ctx context.Context, input *DeleteExerciseInput) (*DeleteExerciseOutput, error) {
	var entity models.ExerciseEntity
	if err := database.DB.First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Exercise not found")
	}
	if entity.Owner != humaconfig.GetUserID(ctx) {
		return nil, huma.Error403Forbidden("access denied")
	}
	if err := database.DB.Unscoped().Select(clause.Associations).Delete(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return nil, nil
}

// DeleteExerciseVersion deletes a specific version from exercise history. Owner only.
// DELETE /api/exercises/:id/versions/:version
//
// OpenAPI: /api/docs#/operations/DeleteExerciseVersion
func DeleteExerciseVersion(ctx context.Context, input *DeleteExerciseVersionInput) (*DeleteExerciseVersionOutput, error) {
	var entity models.ExerciseEntity
	if err := database.DB.First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Exercise not found")
	}
	if entity.Owner != humaconfig.GetUserID(ctx) {
		return nil, huma.Error403Forbidden("access denied")
	}
	result := database.DB.Where("exercise_id = ? AND version = ?", input.ID, input.Version).Delete(&models.ExerciseHistoryEntity{})
	if result.Error != nil {
		return nil, huma.Error500InternalServerError(result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return nil, huma.Error404NotFound("Version not found")
	}
	return nil, nil
}

// DeleteAllExerciseVersions deletes all version history for an exercise. Owner only.
// DELETE /api/exercises/:id/versions
//
// OpenAPI: /api/docs#/operations/DeleteAllExerciseVersions
func DeleteAllExerciseVersions(ctx context.Context, input *DeleteAllExerciseVersionsInput) (*DeleteAllExerciseVersionsOutput, error) {
	var entity models.ExerciseEntity
	if err := database.DB.First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Exercise not found")
	}
	if entity.Owner != humaconfig.GetUserID(ctx) {
		return nil, huma.Error403Forbidden("access denied")
	}
	if err := database.DB.Where("exercise_id = ?", input.ID).Delete(&models.ExerciseHistoryEntity{}).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return nil, nil
}
