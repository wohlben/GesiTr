package workoutschedule

import (
	"fmt"
	"time"

	workoutmodels "gesitr/internal/user/workout/models"
	workoutlogmodels "gesitr/internal/user/workoutlog/models"
	"gesitr/internal/user/workoutschedule/models"

	"gorm.io/gorm"
)

// GenerateForUser runs the two-phase generation for all active schedules
// belonging to the given user. It is idempotent.
//
// Phase 1 (Clone): if the last period has ended and no next period exists,
// clone it forward with the same commitment pattern.
//
// Phase 2 (Activate): for periods where periodStart ≤ now, create WorkoutLogs
// for each unlinked commitment.
func GenerateForUser(db *gorm.DB, userID string, now time.Time) error {
	var schedules []models.WorkoutScheduleEntity
	if err := db.Where("owner = ?", userID).Find(&schedules).Error; err != nil {
		return fmt.Errorf("query schedules: %w", err)
	}

	for i := range schedules {
		if !schedules[i].IsActive(now) {
			continue
		}
		if err := generateForSchedule(db, &schedules[i], now); err != nil {
			return fmt.Errorf("generate for schedule %d: %w", schedules[i].ID, err)
		}
	}
	return nil
}

func generateForSchedule(db *gorm.DB, schedule *models.WorkoutScheduleEntity, now time.Time) error {
	// Phase 1: Clone — ensure there is always a planned period ahead.
	// Loop because if the app was offline, multiple periods may need to be
	// cloned forward before one lands in the future.
	for {
		cloned, err := cloneIfNeeded(db, schedule, now)
		if err != nil {
			return err
		}
		if !cloned {
			break
		}
	}

	// Phase 2: Activate — create WorkoutLogs for unlinked commitments in active periods
	if err := activateCommitments(db, schedule, now); err != nil {
		return err
	}

	return nil
}

// cloneIfNeeded clones the last period forward if no planned period exists.
// Returns true if a period was cloned (caller should loop until false).
func cloneIfNeeded(db *gorm.DB, schedule *models.WorkoutScheduleEntity, now time.Time) (bool, error) {
	// If there's already a planned period (starts in the future), nothing to do
	var plannedCount int64
	db.Model(&models.SchedulePeriodEntity{}).
		Where("schedule_id = ? AND period_start > ?", schedule.ID, now).
		Count(&plannedCount)
	if plannedCount > 0 {
		return false, nil
	}

	// Find the last period for this schedule (template for cloning)
	var lastPeriod models.SchedulePeriodEntity
	err := db.Where("schedule_id = ?", schedule.ID).
		Order("period_end DESC").First(&lastPeriod).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil // no periods yet — user hasn't configured the first one
	}
	if err != nil {
		return false, fmt.Errorf("find last period: %w", err)
	}

	// Compute next period
	nextStart := startOfDay(lastPeriod.PeriodEnd.Add(24 * time.Hour))
	var nextEnd time.Time
	if lastPeriod.Mode == models.PeriodModeMonthly {
		// Monthly: advance one calendar month from the original start's day offset
		nextEnd = startOfDay(nextStart.AddDate(0, 1, -1))
	} else {
		// Normal: same duration in days
		duration := lastPeriod.PeriodEnd.Sub(lastPeriod.PeriodStart)
		nextEnd = nextStart.Add(duration)
	}

	// Check schedule end date
	if schedule.EndDate != nil && nextStart.After(*schedule.EndDate) {
		return false, nil
	}

	return true, db.Transaction(func(tx *gorm.DB) error {
		newPeriod := models.SchedulePeriodEntity{
			ScheduleID:  schedule.ID,
			PeriodStart: nextStart,
			PeriodEnd:   nextEnd,
			Type:        lastPeriod.Type,
			Mode:        lastPeriod.Mode,
		}
		if err := tx.Create(&newPeriod).Error; err != nil {
			return fmt.Errorf("create period: %w", err)
		}

		// Clone commitments from the template (last) period
		var templateCommitments []models.ScheduleCommitmentEntity
		if err := tx.Where("period_id = ?", lastPeriod.ID).Find(&templateCommitments).Error; err != nil {
			return fmt.Errorf("find template commitments: %w", err)
		}

		for _, tc := range templateCommitments {
			newCommitment := models.ScheduleCommitmentEntity{
				PeriodID: newPeriod.ID,
			}
			// For fixed_date: preserve the day offset within the period
			if tc.Date != nil {
				offset := tc.Date.Sub(lastPeriod.PeriodStart)
				newDate := nextStart.Add(offset)
				newCommitment.Date = &newDate
			}
			if err := tx.Create(&newCommitment).Error; err != nil {
				return fmt.Errorf("clone commitment: %w", err)
			}
		}

		return nil
	})
}

func activateCommitments(db *gorm.DB, schedule *models.WorkoutScheduleEntity, now time.Time) error {
	// Find periods that are active (periodStart ≤ now)
	var periods []models.SchedulePeriodEntity
	if err := db.Where("schedule_id = ? AND period_start <= ?", schedule.ID, now).
		Find(&periods).Error; err != nil {
		return fmt.Errorf("find active periods: %w", err)
	}

	// Fetch the workout name
	var workout workoutmodels.WorkoutEntity
	if err := db.First(&workout, schedule.WorkoutID).Error; err != nil {
		return fmt.Errorf("workout %d not found: %w", schedule.WorkoutID, err)
	}

	status := workoutlogmodels.WorkoutLogStatus(schedule.InitialStatus)

	for _, period := range periods {
		// Find unlinked commitments for this period
		var commitments []models.ScheduleCommitmentEntity
		if err := db.Where("period_id = ? AND workout_log_id IS NULL", period.ID).
			Find(&commitments).Error; err != nil {
			return fmt.Errorf("find unlinked commitments: %w", err)
		}

		for _, commitment := range commitments {
			if err := db.Transaction(func(tx *gorm.DB) error {
				log := workoutlogmodels.WorkoutLogEntity{
					Owner:      schedule.Owner,
					WorkoutID:  &schedule.WorkoutID,
					Name:       workout.Name,
					Status:     status,
					ScheduleID: &schedule.ID,
					PeriodID:   &period.ID,
					DueStart:   &period.PeriodStart,
					DueEnd:     &period.PeriodEnd,
					Date:       commitment.Date,
				}
				if err := tx.Create(&log).Error; err != nil {
					return fmt.Errorf("create log: %w", err)
				}

				// Link the commitment to the newly created log
				if err := tx.Model(&commitment).Update("workout_log_id", log.ID).Error; err != nil {
					return fmt.Errorf("link commitment: %w", err)
				}

				return nil
			}); err != nil {
				return err
			}
		}
	}
	return nil
}

func startOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
