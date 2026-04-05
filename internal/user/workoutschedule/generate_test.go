package workoutschedule

import (
	"testing"
	"time"

	exercisemodels "gesitr/internal/compendium/exercise/models"
	workoutmodels "gesitr/internal/compendium/workout/models"
	exerciseschememodels "gesitr/internal/user/exercisescheme/models"
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
		&exercisemodels.ExerciseEntity{},
		&exerciseschememodels.ExerciseSchemeEntity{},
		&exerciseschememodels.ExerciseSchemeSectionItemEntity{},
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
	return db
}

func createWorkout(t *testing.T, db *gorm.DB) workoutmodels.WorkoutEntity {
	t.Helper()
	w := workoutmodels.WorkoutEntity{Name: "Test Workout"}
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
		Timezone:      "UTC",
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

	// Create a period that ended yesterday — chain-cloning will produce
	// an active clone (starts today) and then a planned clone (future).
	now := time.Now().UTC()
	lastWeek := startOfDayIn(now.AddDate(0, 0, -8), time.UTC)
	yesterday := startOfDayIn(now.AddDate(0, 0, -1), time.UTC)
	day1 := lastWeek
	day4 := lastWeek.AddDate(0, 0, 3)

	schedule, _ := createScheduleWithPeriodAndCommitments(t, db, workout,
		lastWeek, yesterday, []*time.Time{&day1, &day4})

	// Run generation — should clone the period forward (possibly multiple times)
	if err := GenerateForUser(db, "alice", now); err != nil {
		t.Fatal(err)
	}

	// Should have at least 2 periods; may have 3 if the first clone also started
	var periods []models.SchedulePeriodEntity
	db.Where("schedule_id = ?", schedule.ID).Order("period_start").Find(&periods)
	if len(periods) < 2 {
		t.Fatalf("expected at least 2 periods, got %d", len(periods))
	}

	// First clone should start the day after the original ended
	expectedStart := startOfDayIn(yesterday.AddDate(0, 0, 1), time.UTC)
	if !periods[1].PeriodStart.Equal(expectedStart) {
		t.Errorf("expected first clone to start at %v, got %v", expectedStart, periods[1].PeriodStart)
	}

	// All cloned periods should have same duration in calendar days
	origDays := int(yesterday.Sub(lastWeek).Hours()/24 + 0.5)
	for i := 1; i < len(periods); i++ {
		dur := int(periods[i].PeriodEnd.Sub(periods[i].PeriodStart).Hours()/24 + 0.5)
		if dur != origDays {
			t.Errorf("period %d duration mismatch: %d days vs %d days", i, dur, origDays)
		}
	}

	// First clone should have 2 commitments (same as template)
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

	// The last period should be planned (starts in the future)
	lastClone := periods[len(periods)-1]
	if !time.Now().Before(lastClone.PeriodStart) {
		t.Error("last period should be planned (in the future)")
	}
}

func TestClone_Idempotent(t *testing.T) {
	db := setupTestDB(t)
	workout := createWorkout(t, db)

	// Period ended yesterday — clone starts today (active), so a second clone
	// is needed to produce a planned period. Running twice should not create more.
	lastWeek := startOfDay(time.Now().AddDate(0, 0, -8))
	yesterday := startOfDay(time.Now().AddDate(0, 0, -1))

	createScheduleWithPeriodAndCommitments(t, db, workout,
		lastWeek, yesterday, []*time.Time{&lastWeek})

	// Run twice
	GenerateForUser(db, "alice", time.Now())
	GenerateForUser(db, "alice", time.Now())

	var periodCount int64
	db.Model(&models.SchedulePeriodEntity{}).Count(&periodCount)
	// 3 periods: original (archived) + clone (active, starts today) + clone (planned, future)
	if periodCount != 3 {
		t.Errorf("expected exactly 3 periods (archived + active + planned), got %d", periodCount)
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

func TestClone_TriggeredWhenPeriodBecomesActive(t *testing.T) {
	db := setupTestDB(t)
	workout := createWorkout(t, db)

	// Create a period that started yesterday but ends next week (currently active)
	now := time.Now().UTC()
	yesterday := startOfDayIn(now.AddDate(0, 0, -1), time.UTC)
	nextWeek := startOfDayIn(now.AddDate(0, 0, 6), time.UTC)
	day1 := yesterday

	schedule, _ := createScheduleWithPeriodAndCommitments(t, db, workout,
		yesterday, nextWeek, []*time.Time{&day1})

	// Run generation — should clone forward because the period is now active
	if err := GenerateForUser(db, "alice", now); err != nil {
		t.Fatal(err)
	}

	// Should now have 2 periods (active + planned clone)
	var periods []models.SchedulePeriodEntity
	db.Where("schedule_id = ?", schedule.ID).Order("period_start").Find(&periods)
	if len(periods) != 2 {
		t.Fatalf("expected 2 periods (active + planned clone), got %d", len(periods))
	}

	// The cloned period should start after the current one ends
	expectedStart := startOfDayIn(nextWeek.AddDate(0, 0, 1), time.UTC)
	if !periods[1].PeriodStart.Equal(expectedStart) {
		t.Errorf("expected clone start %v, got %v", expectedStart, periods[1].PeriodStart)
	}

	// The cloned period should be in the future (planned)
	if !now.Before(periods[1].PeriodStart) {
		t.Error("cloned period should be in the future (planned)")
	}
}

func TestClone_NoCloneWhenPeriodStillPlanned(t *testing.T) {
	db := setupTestDB(t)
	workout := createWorkout(t, db)

	// Period starts tomorrow (still planned)
	tomorrow := startOfDay(time.Now().AddDate(0, 0, 1))
	nextWeek := startOfDay(time.Now().AddDate(0, 0, 8))

	schedule, _ := createScheduleWithPeriodAndCommitments(t, db, workout,
		tomorrow, nextWeek, []*time.Time{&tomorrow})

	GenerateForUser(db, "alice", time.Now())

	var periodCount int64
	db.Model(&models.SchedulePeriodEntity{}).
		Where("schedule_id = ?", schedule.ID).Count(&periodCount)
	if periodCount != 1 {
		t.Errorf("expected 1 period (no clone for planned), got %d", periodCount)
	}
}

func TestActivation_SnapshotsWorkoutStructure(t *testing.T) {
	db := setupTestDB(t)

	// Create exercise + scheme with 3 sets of 10 reps at 60kg
	exercise := exercisemodels.ExerciseEntity{Names: []exercisemodels.ExerciseName{{Position: 0, Name: "Bench Press"}}}
	db.Create(&exercise)

	sets := intPtr(3)
	reps := intPtr(10)
	weight := float64Ptr(60.0)
	rest := intPtr(90)
	scheme := exerciseschememodels.ExerciseSchemeEntity{
		Owner:           "alice",
		ExerciseID:      exercise.ID,
		MeasurementType: "REP_BASED",
		Sets:            sets,
		Reps:            reps,
		Weight:          weight,
		RestBetweenSets: rest,
	}
	db.Create(&scheme)

	// Create workout with one section and one exercise item
	workout := workoutmodels.WorkoutEntity{Name: "Push Day"}
	db.Create(&workout)

	section := workoutmodels.WorkoutSectionEntity{
		WorkoutID: workout.ID,
		Type:      "main",
		Position:  0,
	}
	db.Create(&section)

	item := workoutmodels.WorkoutSectionItemEntity{
		WorkoutSectionID: section.ID,
		Type:             workoutmodels.WorkoutSectionItemTypeExercise,
		ExerciseID:       &exercise.ID,
		Position:         0,
	}
	db.Create(&item)

	// Link scheme to section item via join table
	link := exerciseschememodels.ExerciseSchemeSectionItemEntity{
		ExerciseSchemeID:     scheme.ID,
		WorkoutSectionItemID: item.ID,
		Owner:                "alice",
	}
	db.Create(&link)

	// Create schedule + active period + commitment
	yesterday := startOfDay(time.Now().AddDate(0, 0, -1))
	nextWeek := startOfDay(time.Now().AddDate(0, 0, 6))
	commitDate := startOfDay(time.Now().AddDate(0, 0, 2))

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
	db.Create(&models.ScheduleCommitmentEntity{PeriodID: period.ID, Date: &commitDate})

	// Run generation
	if err := GenerateForUser(db, "alice", time.Now()); err != nil {
		t.Fatal(err)
	}

	// Verify the log was created with structure
	var log workoutlogmodels.WorkoutLogEntity
	db.Where("schedule_id = ?", schedule.ID).First(&log)

	if log.Status != workoutlogmodels.WorkoutLogStatusProposed {
		t.Errorf("expected proposed, got %s", log.Status)
	}

	// Verify log sections
	var logSections []workoutlogmodels.WorkoutLogSectionEntity
	db.Where("workout_log_id = ?", log.ID).Find(&logSections)
	if len(logSections) != 1 {
		t.Fatalf("expected 1 log section, got %d", len(logSections))
	}
	if logSections[0].Type != "main" {
		t.Errorf("expected section type main, got %s", logSections[0].Type)
	}

	// Verify log exercises
	var logExercises []workoutlogmodels.WorkoutLogExerciseEntity
	db.Where("workout_log_id = ?", log.ID).Find(&logExercises)
	if len(logExercises) != 1 {
		t.Fatalf("expected 1 log exercise, got %d", len(logExercises))
	}
	if logExercises[0].SourceExerciseSchemeID != scheme.ID {
		t.Errorf("expected scheme ID %d, got %d", scheme.ID, logExercises[0].SourceExerciseSchemeID)
	}
	if logExercises[0].TargetMeasurementType != "REP_BASED" {
		t.Errorf("expected REP_BASED, got %s", logExercises[0].TargetMeasurementType)
	}

	// Verify log sets (3 sets with snapshotted targets)
	var logSets []workoutlogmodels.WorkoutLogExerciseSetEntity
	db.Where("workout_log_id = ?", log.ID).Order("set_number").Find(&logSets)
	if len(logSets) != 3 {
		t.Fatalf("expected 3 log sets, got %d", len(logSets))
	}
	for i, s := range logSets {
		if s.SetNumber != i+1 {
			t.Errorf("set %d: expected set_number %d, got %d", i, i+1, s.SetNumber)
		}
		if s.TargetReps == nil || *s.TargetReps != 10 {
			t.Errorf("set %d: expected target_reps 10", i)
		}
		if s.TargetWeight == nil || *s.TargetWeight != 60.0 {
			t.Errorf("set %d: expected target_weight 60", i)
		}
	}
	// First two sets should have rest, last should not
	if logSets[0].BreakAfterSeconds == nil || *logSets[0].BreakAfterSeconds != 90 {
		t.Error("set 0 should have 90s break")
	}
	if logSets[2].BreakAfterSeconds != nil {
		t.Error("last set should have no break")
	}
}

func TestActivation_LogDateMatchesCommitmentDate(t *testing.T) {
	db := setupTestDB(t)
	workout := createWorkout(t, db)

	yesterday := startOfDay(time.Now().AddDate(0, 0, -1))
	nextWeek := startOfDay(time.Now().AddDate(0, 0, 6))
	// Commitment on day +3 (not the period start)
	commitDate := startOfDay(time.Now().AddDate(0, 0, 3))

	schedule, period := createScheduleWithPeriodAndCommitments(t, db, workout,
		yesterday, nextWeek, []*time.Time{&commitDate})

	if err := GenerateForUser(db, "alice", time.Now()); err != nil {
		t.Fatal(err)
	}

	var log workoutlogmodels.WorkoutLogEntity
	db.Where("schedule_id = ?", schedule.ID).First(&log)

	// Date should be the commitment date, not the period start
	if log.Date == nil {
		t.Fatal("log date should not be nil")
	}
	if !startOfDay(*log.Date).Equal(commitDate) {
		t.Errorf("log date should be commitment date %v, got %v", commitDate, *log.Date)
	}

	// DueStart should be the period start
	if log.DueStart == nil || !startOfDay(*log.DueStart).Equal(period.PeriodStart) {
		t.Errorf("due_start should be period start %v", period.PeriodStart)
	}
}

func TestClone_RespectsTimezoneEastOfUTC(t *testing.T) {
	db := setupTestDB(t)
	workout := createWorkout(t, db)

	// Simulate a user in Europe/Berlin (UTC+2 during CEST).
	// Period "April 6 – April 13" in Berlin = stored as April 5 22:00 – April 12 22:00 UTC.
	berlin, _ := time.LoadLocation("Europe/Berlin")
	periodStart := time.Date(2026, 4, 6, 0, 0, 0, 0, berlin)
	periodEnd := time.Date(2026, 4, 13, 0, 0, 0, 0, berlin)

	schedule := models.WorkoutScheduleEntity{
		Owner:         "alice",
		WorkoutID:     workout.ID,
		StartDate:     periodStart,
		InitialStatus: "committed",
		Timezone:      "Europe/Berlin",
	}
	db.Create(&schedule)

	period := models.SchedulePeriodEntity{
		ScheduleID:  schedule.ID,
		PeriodStart: periodStart,
		PeriodEnd:   periodEnd,
		Type:        models.ScheduleTypeFrequency,
	}
	db.Create(&period)
	db.Create(&models.ScheduleCommitmentEntity{PeriodID: period.ID})

	// Run at a time when the period is active (April 13 10:00 Berlin = April 13 08:00 UTC)
	now := time.Date(2026, 4, 13, 8, 0, 0, 0, time.UTC)
	if err := GenerateForUser(db, "alice", now); err != nil {
		t.Fatal(err)
	}

	var periods []models.SchedulePeriodEntity
	db.Where("schedule_id = ?", schedule.ID).Order("period_start").Find(&periods)
	if len(periods) < 2 {
		t.Fatalf("expected at least 2 periods, got %d", len(periods))
	}

	clone := periods[1]
	// Next period should start April 14 00:00 Berlin (April 13 22:00 UTC)
	expectedStart := time.Date(2026, 4, 14, 0, 0, 0, 0, berlin)
	if !clone.PeriodStart.Equal(expectedStart) {
		t.Errorf("clone start: want %v, got %v", expectedStart, clone.PeriodStart)
	}

	// Should end April 21 00:00 Berlin (same 7-day duration)
	expectedEnd := time.Date(2026, 4, 21, 0, 0, 0, 0, berlin)
	if !clone.PeriodEnd.Equal(expectedEnd) {
		t.Errorf("clone end: want %v, got %v", expectedEnd, clone.PeriodEnd)
	}

	// No overlap: clone starts strictly after original ends
	if !clone.PeriodStart.After(periodEnd) || clone.PeriodStart.Sub(periodEnd) > 25*time.Hour {
		t.Errorf("clone should start exactly 1 day after original end, gap = %v", clone.PeriodStart.Sub(periodEnd))
	}
}

func TestClone_RespectsTimezoneWestOfUTC(t *testing.T) {
	db := setupTestDB(t)
	workout := createWorkout(t, db)

	// Simulate a user in America/New_York (UTC-4 during EDT).
	// Period "April 6 – April 13" in NY = stored as April 6 04:00 – April 13 04:00 UTC.
	ny, _ := time.LoadLocation("America/New_York")
	periodStart := time.Date(2026, 4, 6, 0, 0, 0, 0, ny)
	periodEnd := time.Date(2026, 4, 13, 0, 0, 0, 0, ny)

	schedule := models.WorkoutScheduleEntity{
		Owner:         "alice",
		WorkoutID:     workout.ID,
		StartDate:     periodStart,
		InitialStatus: "committed",
		Timezone:      "America/New_York",
	}
	db.Create(&schedule)

	period := models.SchedulePeriodEntity{
		ScheduleID:  schedule.ID,
		PeriodStart: periodStart,
		PeriodEnd:   periodEnd,
		Type:        models.ScheduleTypeFrequency,
	}
	db.Create(&period)
	db.Create(&models.ScheduleCommitmentEntity{PeriodID: period.ID})

	// Run at a time when the period is active (April 13 10:00 NY = April 13 14:00 UTC)
	now := time.Date(2026, 4, 13, 14, 0, 0, 0, time.UTC)
	if err := GenerateForUser(db, "alice", now); err != nil {
		t.Fatal(err)
	}

	var periods []models.SchedulePeriodEntity
	db.Where("schedule_id = ?", schedule.ID).Order("period_start").Find(&periods)
	if len(periods) < 2 {
		t.Fatalf("expected at least 2 periods, got %d", len(periods))
	}

	clone := periods[1]
	// Next period should start April 14 00:00 NY (April 14 04:00 UTC)
	expectedStart := time.Date(2026, 4, 14, 0, 0, 0, 0, ny)
	if !clone.PeriodStart.Equal(expectedStart) {
		t.Errorf("clone start: want %v, got %v", expectedStart, clone.PeriodStart)
	}

	// Should end April 21 00:00 NY (same 7-day duration)
	expectedEnd := time.Date(2026, 4, 21, 0, 0, 0, 0, ny)
	if !clone.PeriodEnd.Equal(expectedEnd) {
		t.Errorf("clone end: want %v, got %v", expectedEnd, clone.PeriodEnd)
	}
}

func TestClone_UTCBackwardCompatible(t *testing.T) {
	db := setupTestDB(t)
	workout := createWorkout(t, db)

	// Default timezone (empty/UTC) should behave identically to old code
	periodStart := time.Date(2026, 4, 6, 0, 0, 0, 0, time.UTC)
	periodEnd := time.Date(2026, 4, 13, 0, 0, 0, 0, time.UTC)

	schedule := models.WorkoutScheduleEntity{
		Owner:         "alice",
		WorkoutID:     workout.ID,
		StartDate:     periodStart,
		InitialStatus: "committed",
		// Timezone left empty — Location() falls back to UTC
	}
	db.Create(&schedule)

	period := models.SchedulePeriodEntity{
		ScheduleID:  schedule.ID,
		PeriodStart: periodStart,
		PeriodEnd:   periodEnd,
		Type:        models.ScheduleTypeFrequency,
	}
	db.Create(&period)
	db.Create(&models.ScheduleCommitmentEntity{PeriodID: period.ID})

	now := time.Date(2026, 4, 13, 12, 0, 0, 0, time.UTC)
	if err := GenerateForUser(db, "alice", now); err != nil {
		t.Fatal(err)
	}

	var periods []models.SchedulePeriodEntity
	db.Where("schedule_id = ?", schedule.ID).Order("period_start").Find(&periods)
	if len(periods) < 2 {
		t.Fatalf("expected at least 2 periods, got %d", len(periods))
	}

	clone := periods[1]
	expectedStart := time.Date(2026, 4, 14, 0, 0, 0, 0, time.UTC)
	if !clone.PeriodStart.Equal(expectedStart) {
		t.Errorf("clone start: want %v, got %v", expectedStart, clone.PeriodStart)
	}

	expectedEnd := time.Date(2026, 4, 21, 0, 0, 0, 0, time.UTC)
	if !clone.PeriodEnd.Equal(expectedEnd) {
		t.Errorf("clone end: want %v, got %v", expectedEnd, clone.PeriodEnd)
	}
}

func intPtr(v int) *int             { return &v }
func float64Ptr(v float64) *float64 { return &v }
func ptr(t time.Time) *time.Time    { return &t }
