package workoutschedule

import (
	"testing"
	"time"

	profilemodels "gesitr/internal/profile/models"
	workoutmodels "gesitr/internal/user/workout/models"
	workoutlogmodels "gesitr/internal/user/workoutlog/models"
	"gesitr/internal/user/workoutschedule/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?_foreign_keys=on"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatal(err)
	}
	db.AutoMigrate(
		&profilemodels.UserProfileEntity{},
		&workoutmodels.WorkoutEntity{},
		&workoutmodels.WorkoutSectionEntity{},
		&workoutmodels.WorkoutSectionItemEntity{},
		&workoutlogmodels.WorkoutLogEntity{},
		&workoutlogmodels.WorkoutLogSectionEntity{},
		&workoutlogmodels.WorkoutLogExerciseEntity{},
		&workoutlogmodels.WorkoutLogExerciseSetEntity{},
		&models.WorkoutScheduleEntity{},
		&models.SchedulePeriodEntity{},
		&models.ScheduleCommitmentEntity{},
	)
	db.Create(&profilemodels.UserProfileEntity{ID: "alice", Name: "alice"})
	return db
}

func createWorkout(t *testing.T, db *gorm.DB) workoutmodels.WorkoutEntity {
	t.Helper()
	w := workoutmodels.WorkoutEntity{Owner: "alice", Name: "Test Workout"}
	if err := db.Create(&w).Error; err != nil {
		t.Fatal(err)
	}
	return w
}

func createScheduleWithPeriodAndCommitments(t *testing.T, db *gorm.DB, workout workoutmodels.WorkoutEntity, periodStart, periodEnd time.Time, commitmentDates []*time.Time) (models.WorkoutScheduleEntity, models.SchedulePeriodEntity) {
	t.Helper()

	schedule := models.WorkoutScheduleEntity{
		Owner:         "alice",
		WorkoutID:     workout.ID,
		StartDate:     periodStart,
		InitialStatus: "committed",
	}
	db.Create(&schedule)

	period := models.SchedulePeriodEntity{
		ScheduleID:  schedule.ID,
		PeriodStart: periodStart,
		PeriodEnd:   periodEnd,
		Type:        models.ScheduleTypeFixedDate,
	}
	db.Create(&period)

	for _, d := range commitmentDates {
		db.Create(&models.ScheduleCommitmentEntity{
			PeriodID: period.ID,
			Date:     d,
		})
	}

	return schedule, period
}

func TestActivation_CreatesWorkoutLogs(t *testing.T) {
	db := setupTestDB(t)
	workout := createWorkout(t, db)

	// Create a period that started yesterday with 2 commitments
	yesterday := startOfDay(time.Now().AddDate(0, 0, -1))
	nextWeek := startOfDay(time.Now().AddDate(0, 0, 6))
	day1 := yesterday
	day3 := yesterday.AddDate(0, 0, 2)

	schedule, period := createScheduleWithPeriodAndCommitments(t, db, workout,
		yesterday, nextWeek, []*time.Time{&day1, &day3})

	// Run generation — should activate (period has started)
	if err := GenerateForUser(db, "alice", time.Now()); err != nil {
		t.Fatal(err)
	}

	// Verify workout logs were created
	var logs []workoutlogmodels.WorkoutLogEntity
	db.Where("schedule_id = ?", schedule.ID).Find(&logs)
	if len(logs) != 2 {
		t.Fatalf("expected 2 logs, got %d", len(logs))
	}

	// Verify commitments are now linked
	var commitments []models.ScheduleCommitmentEntity
	db.Where("period_id = ?", period.ID).Find(&commitments)
	for _, c := range commitments {
		if c.WorkoutLogID == nil {
			t.Error("commitment should be linked to a workout log after activation")
		}
	}

	// Verify log properties
	for _, log := range logs {
		if log.Status != workoutlogmodels.WorkoutLogStatusCommitted {
			t.Errorf("expected committed status, got %s", log.Status)
		}
		if log.Name != "Test Workout" {
			t.Errorf("expected workout name, got %s", log.Name)
		}
	}
}

func TestActivation_Idempotent(t *testing.T) {
	db := setupTestDB(t)
	workout := createWorkout(t, db)

	yesterday := startOfDay(time.Now().AddDate(0, 0, -1))
	nextWeek := startOfDay(time.Now().AddDate(0, 0, 6))

	schedule, _ := createScheduleWithPeriodAndCommitments(t, db, workout,
		yesterday, nextWeek, []*time.Time{&yesterday})

	// Run twice
	GenerateForUser(db, "alice", time.Now())
	GenerateForUser(db, "alice", time.Now())

	var logCount int64
	db.Model(&workoutlogmodels.WorkoutLogEntity{}).Where("schedule_id = ?", schedule.ID).Count(&logCount)
	if logCount != 1 {
		t.Errorf("expected 1 log (idempotent), got %d", logCount)
	}
}

func TestActivation_FuturePeriodNotActivated(t *testing.T) {
	db := setupTestDB(t)
	workout := createWorkout(t, db)

	// Period starts tomorrow — should NOT activate
	tomorrow := startOfDay(time.Now().AddDate(0, 0, 1))
	nextWeek := startOfDay(time.Now().AddDate(0, 0, 8))

	schedule, _ := createScheduleWithPeriodAndCommitments(t, db, workout,
		tomorrow, nextWeek, []*time.Time{&tomorrow})

	GenerateForUser(db, "alice", time.Now())

	var logCount int64
	db.Model(&workoutlogmodels.WorkoutLogEntity{}).Where("schedule_id = ?", schedule.ID).Count(&logCount)
	if logCount != 0 {
		t.Errorf("future period should not activate, got %d logs", logCount)
	}

	// Commitments should still be unlinked
	var commitments []models.ScheduleCommitmentEntity
	db.Where("workout_log_id IS NOT NULL").Find(&commitments)
	if len(commitments) != 0 {
		t.Error("no commitments should be linked for future periods")
	}
}

func TestClone_ClonesLastPeriod(t *testing.T) {
	db := setupTestDB(t)
	workout := createWorkout(t, db)

	// Create a period that ended yesterday
	lastWeek := startOfDay(time.Now().AddDate(0, 0, -8))
	yesterday := startOfDay(time.Now().AddDate(0, 0, -1))
	day1 := lastWeek
	day4 := lastWeek.AddDate(0, 0, 3)

	schedule, _ := createScheduleWithPeriodAndCommitments(t, db, workout,
		lastWeek, yesterday, []*time.Time{&day1, &day4})

	// Run generation — should clone the period forward
	if err := GenerateForUser(db, "alice", time.Now()); err != nil {
		t.Fatal(err)
	}

	// Should now have 2 periods
	var periods []models.SchedulePeriodEntity
	db.Where("schedule_id = ?", schedule.ID).Order("period_start").Find(&periods)
	if len(periods) != 2 {
		t.Fatalf("expected 2 periods (original + clone), got %d", len(periods))
	}

	// New period should have same duration
	origDuration := yesterday.Sub(lastWeek)
	newDuration := periods[1].PeriodEnd.Sub(periods[1].PeriodStart)
	if origDuration != newDuration {
		t.Errorf("cloned period duration mismatch: %v vs %v", origDuration, newDuration)
	}

	// New period should start the day after the old one ended
	expectedStart := startOfDay(yesterday.AddDate(0, 0, 1))
	if !periods[1].PeriodStart.Equal(expectedStart) {
		t.Errorf("expected clone to start at %v, got %v", expectedStart, periods[1].PeriodStart)
	}

	// New period should have 2 commitments (same as template)
	var newCommitments []models.ScheduleCommitmentEntity
	db.Where("period_id = ?", periods[1].ID).Find(&newCommitments)
	if len(newCommitments) != 2 {
		t.Fatalf("expected 2 cloned commitments, got %d", len(newCommitments))
	}

	// Verify day offsets are preserved
	for i, c := range newCommitments {
		if c.Date == nil {
			t.Errorf("commitment %d should have a date (fixed_date clone)", i)
		}
	}
}

func TestClone_Idempotent(t *testing.T) {
	db := setupTestDB(t)
	workout := createWorkout(t, db)

	lastWeek := startOfDay(time.Now().AddDate(0, 0, -8))
	yesterday := startOfDay(time.Now().AddDate(0, 0, -1))

	createScheduleWithPeriodAndCommitments(t, db, workout,
		lastWeek, yesterday, []*time.Time{&lastWeek})

	// Run twice
	GenerateForUser(db, "alice", time.Now())
	GenerateForUser(db, "alice", time.Now())

	var periodCount int64
	db.Model(&models.SchedulePeriodEntity{}).Count(&periodCount)
	if periodCount != 2 {
		t.Errorf("expected exactly 2 periods (original + 1 clone), got %d", periodCount)
	}
}

func TestScheduleIsActive(t *testing.T) {
	now := time.Date(2026, 3, 30, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		start    time.Time
		end      *time.Time
		expected bool
	}{
		{"active (started, no end)", now.AddDate(0, 0, -1), nil, true},
		{"active (started, future end)", now.AddDate(0, 0, -1), ptr(now.AddDate(0, 0, 7)), true},
		{"inactive (not started)", now.AddDate(0, 0, 1), nil, false},
		{"inactive (ended)", now.AddDate(0, 0, -7), ptr(now.AddDate(0, 0, -1)), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := models.WorkoutScheduleEntity{StartDate: tt.start, EndDate: tt.end}
			if s.IsActive(now) != tt.expected {
				t.Errorf("IsActive() = %v, want %v", !tt.expected, tt.expected)
			}
		})
	}
}

func TestActivation_InitialStatusProposed(t *testing.T) {
	db := setupTestDB(t)
	workout := createWorkout(t, db)

	yesterday := startOfDay(time.Now().AddDate(0, 0, -1))
	nextWeek := startOfDay(time.Now().AddDate(0, 0, 6))

	schedule := models.WorkoutScheduleEntity{
		Owner:         "alice",
		WorkoutID:     workout.ID,
		StartDate:     yesterday,
		InitialStatus: "proposed",
	}
	db.Create(&schedule)

	period := models.SchedulePeriodEntity{
		ScheduleID:  schedule.ID,
		PeriodStart: yesterday,
		PeriodEnd:   nextWeek,
		Type:        models.ScheduleTypeFixedDate,
	}
	db.Create(&period)
	db.Create(&models.ScheduleCommitmentEntity{PeriodID: period.ID, Date: &yesterday})

	GenerateForUser(db, "alice", time.Now())

	var log workoutlogmodels.WorkoutLogEntity
	db.Where("schedule_id = ?", schedule.ID).First(&log)
	if log.Status != workoutlogmodels.WorkoutLogStatusProposed {
		t.Errorf("expected proposed status, got %s", log.Status)
	}
}

func TestInactiveScheduleSkipped(t *testing.T) {
	db := setupTestDB(t)
	workout := createWorkout(t, db)

	// Schedule that hasn't started yet
	tomorrow := startOfDay(time.Now().AddDate(0, 0, 1))
	nextWeek := startOfDay(time.Now().AddDate(0, 0, 8))

	schedule := models.WorkoutScheduleEntity{
		Owner:         "alice",
		WorkoutID:     workout.ID,
		StartDate:     tomorrow.AddDate(0, 0, 7), // starts next week
		InitialStatus: "committed",
	}
	db.Create(&schedule)

	period := models.SchedulePeriodEntity{
		ScheduleID:  schedule.ID,
		PeriodStart: tomorrow,
		PeriodEnd:   nextWeek,
		Type:        models.ScheduleTypeFixedDate,
	}
	db.Create(&period)
	db.Create(&models.ScheduleCommitmentEntity{PeriodID: period.ID, Date: &tomorrow})

	GenerateForUser(db, "alice", time.Now())

	var logCount int64
	db.Model(&workoutlogmodels.WorkoutLogEntity{}).Where("schedule_id = ?", schedule.ID).Count(&logCount)
	if logCount != 0 {
		t.Errorf("inactive schedule should not generate logs, got %d", logCount)
	}
}

func TestFrequency_CommitmentsWithoutDates(t *testing.T) {
	db := setupTestDB(t)
	workout := createWorkout(t, db)

	yesterday := startOfDay(time.Now().AddDate(0, 0, -1))
	nextWeek := startOfDay(time.Now().AddDate(0, 0, 6))

	schedule := models.WorkoutScheduleEntity{
		Owner:         "alice",
		WorkoutID:     workout.ID,
		StartDate:     yesterday,
		InitialStatus: "committed",
	}
	db.Create(&schedule)

	period := models.SchedulePeriodEntity{
		ScheduleID:  schedule.ID,
		PeriodStart: yesterday,
		PeriodEnd:   nextWeek,
		Type:        models.ScheduleTypeFrequency,
	}
	db.Create(&period)

	// 3 commitments without specific dates (frequency type)
	for i := 0; i < 3; i++ {
		db.Create(&models.ScheduleCommitmentEntity{PeriodID: period.ID})
	}

	GenerateForUser(db, "alice", time.Now())

	var logs []workoutlogmodels.WorkoutLogEntity
	db.Where("schedule_id = ?", schedule.ID).Find(&logs)
	if len(logs) != 3 {
		t.Fatalf("expected 3 logs for frequency, got %d", len(logs))
	}

	for _, log := range logs {
		if log.Date != nil {
			t.Error("frequency logs should have nil date")
		}
		if log.DueStart == nil || log.DueEnd == nil {
			t.Error("due window should be set")
		}
	}
}

func ptr(t time.Time) *time.Time { return &t }
