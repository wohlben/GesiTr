package handlers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"gesitr/internal/shared"
	schedulemodels "gesitr/internal/user/workoutschedule/models"
)

// TestPeriodStatus verifies that the computed status field is returned
// correctly for planned, active, and archived periods.
func TestPeriodStatus(t *testing.T) {
	setupTestDB(t)
	defer closeDB(t)
	r := newRouter()

	// Create workout + schedule
	doJSONLog(t, r, "POST", "/api/user/workouts", map[string]any{"name": "Status Test"})
	startDate := time.Now().AddDate(0, 0, -10)
	w := doJSONLog(t, r, "POST", "/api/user/workout-schedules", map[string]any{
		"workoutId": 1, "startDate": startDate.Format(time.RFC3339),
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create schedule: %d", w.Code)
	}
	var schedule schedulemodels.WorkoutSchedule
	json.Unmarshal(w.Body.Bytes(), &schedule)

	// Create an archived period (ended yesterday)
	lastWeek := time.Now().AddDate(0, 0, -8).Truncate(24 * time.Hour)
	yesterday := time.Now().AddDate(0, 0, -1).Truncate(24 * time.Hour)
	w = doJSONLog(t, r, "POST", "/api/user/schedule-periods", map[string]any{
		"scheduleId":  schedule.ID,
		"periodStart": lastWeek.Format(time.RFC3339),
		"periodEnd":   yesterday.Format(time.RFC3339),
		"type":        "fixed_date",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create archived period: %d", w.Code)
	}

	// Create an active period (started yesterday, ends next week)
	today := time.Now().Truncate(24 * time.Hour)
	nextWeek := time.Now().AddDate(0, 0, 7).Truncate(24 * time.Hour)
	w = doJSONLog(t, r, "POST", "/api/user/schedule-periods", map[string]any{
		"scheduleId":  schedule.ID,
		"periodStart": today.Format(time.RFC3339),
		"periodEnd":   nextWeek.Format(time.RFC3339),
		"type":        "fixed_date",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create active period: %d", w.Code)
	}

	// Create a planned period (starts next week)
	futureStart := time.Now().AddDate(0, 0, 8).Truncate(24 * time.Hour)
	futureEnd := time.Now().AddDate(0, 0, 15).Truncate(24 * time.Hour)
	w = doJSONLog(t, r, "POST", "/api/user/schedule-periods", map[string]any{
		"scheduleId":  schedule.ID,
		"periodStart": futureStart.Format(time.RFC3339),
		"periodEnd":   futureEnd.Format(time.RFC3339),
		"type":        "fixed_date",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create planned period: %d", w.Code)
	}

	// List periods and verify statuses
	w = doJSONLog(t, r, "GET", fmt.Sprintf("/api/user/schedule-periods?scheduleId=%d", schedule.ID), nil)
	if w.Code != http.StatusOK {
		t.Fatalf("list periods: %d", w.Code)
	}
	var periods []schedulemodels.SchedulePeriod
	json.Unmarshal(w.Body.Bytes(), &periods)
	if len(periods) != 3 {
		t.Fatalf("expected 3 periods, got %d", len(periods))
	}

	// Periods are ordered by period_start
	if periods[0].Status != schedulemodels.PeriodStatusArchived {
		t.Errorf("first period should be archived, got %s", periods[0].Status)
	}
	if periods[1].Status != schedulemodels.PeriodStatusActive {
		t.Errorf("second period should be active, got %s", periods[1].Status)
	}
	if periods[2].Status != schedulemodels.PeriodStatusPlanned {
		t.Errorf("third period should be planned, got %s", periods[2].Status)
	}
}

// TestPeriodPermissions_PlannedIsEditable verifies that a planned period
// returns READ, MODIFY, DELETE permissions.
func TestPeriodPermissions_PlannedIsEditable(t *testing.T) {
	setupTestDB(t)
	defer closeDB(t)
	r := newRouter()

	doJSONLog(t, r, "POST", "/api/user/workouts", map[string]any{"name": "Perms Test"})
	doJSONLog(t, r, "POST", "/api/user/workout-schedules", map[string]any{"workoutId": 1})

	// Create a planned period (starts next week)
	futureStart := time.Now().AddDate(0, 0, 7).Truncate(24 * time.Hour)
	futureEnd := time.Now().AddDate(0, 0, 14).Truncate(24 * time.Hour)
	w := doJSONLog(t, r, "POST", "/api/user/schedule-periods", map[string]any{
		"scheduleId":  1,
		"periodStart": futureStart.Format(time.RFC3339),
		"periodEnd":   futureEnd.Format(time.RFC3339),
		"type":        "fixed_date",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create period: %d", w.Code)
	}
	var period schedulemodels.SchedulePeriod
	json.Unmarshal(w.Body.Bytes(), &period)

	// Get permissions
	w = doJSONLog(t, r, "GET", fmt.Sprintf("/api/user/schedule-periods/%d/permissions", period.ID), nil)
	if w.Code != http.StatusOK {
		t.Fatalf("get permissions: %d", w.Code)
	}
	var perms shared.PermissionsResponse
	json.Unmarshal(w.Body.Bytes(), &perms)

	if len(perms.Permissions) != 3 {
		t.Fatalf("planned period should have 3 permissions, got %d: %v", len(perms.Permissions), perms.Permissions)
	}
	has := func(p shared.Permission) bool {
		for _, v := range perms.Permissions {
			if v == p {
				return true
			}
		}
		return false
	}
	if !has(shared.PermissionRead) {
		t.Error("missing READ")
	}
	if !has(shared.PermissionModify) {
		t.Error("missing MODIFY")
	}
	if !has(shared.PermissionDelete) {
		t.Error("missing DELETE")
	}
}

// TestPeriodPermissions_ActiveIsReadOnly verifies that an active period
// returns only READ permission.
func TestPeriodPermissions_ActiveIsReadOnly(t *testing.T) {
	setupTestDB(t)
	defer closeDB(t)
	r := newRouter()

	doJSONLog(t, r, "POST", "/api/user/workouts", map[string]any{"name": "Active Perms"})
	startDate := time.Now().AddDate(0, 0, -3)
	doJSONLog(t, r, "POST", "/api/user/workout-schedules", map[string]any{
		"workoutId": 1, "startDate": startDate.Format(time.RFC3339),
	})

	// Create an active period
	today := time.Now().Truncate(24 * time.Hour)
	nextWeek := time.Now().AddDate(0, 0, 7).Truncate(24 * time.Hour)
	w := doJSONLog(t, r, "POST", "/api/user/schedule-periods", map[string]any{
		"scheduleId":  1,
		"periodStart": today.Format(time.RFC3339),
		"periodEnd":   nextWeek.Format(time.RFC3339),
		"type":        "fixed_date",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create period: %d", w.Code)
	}
	var period schedulemodels.SchedulePeriod
	json.Unmarshal(w.Body.Bytes(), &period)

	w = doJSONLog(t, r, "GET", fmt.Sprintf("/api/user/schedule-periods/%d/permissions", period.ID), nil)
	if w.Code != http.StatusOK {
		t.Fatalf("get permissions: %d", w.Code)
	}
	var perms shared.PermissionsResponse
	json.Unmarshal(w.Body.Bytes(), &perms)

	if len(perms.Permissions) != 1 {
		t.Fatalf("active period should have 1 permission, got %d: %v", len(perms.Permissions), perms.Permissions)
	}
	if perms.Permissions[0] != shared.PermissionRead {
		t.Errorf("expected READ, got %s", perms.Permissions[0])
	}
}

// TestPeriodPermissions_ArchivedIsReadOnly verifies that an archived period
// returns only READ permission.
func TestPeriodPermissions_ArchivedIsReadOnly(t *testing.T) {
	setupTestDB(t)
	defer closeDB(t)
	r := newRouter()

	doJSONLog(t, r, "POST", "/api/user/workouts", map[string]any{"name": "Archived Perms"})
	startDate := time.Now().AddDate(0, 0, -10)
	doJSONLog(t, r, "POST", "/api/user/workout-schedules", map[string]any{
		"workoutId": 1, "startDate": startDate.Format(time.RFC3339),
	})

	// Create an archived period
	lastWeek := time.Now().AddDate(0, 0, -8).Truncate(24 * time.Hour)
	yesterday := time.Now().AddDate(0, 0, -1).Truncate(24 * time.Hour)
	w := doJSONLog(t, r, "POST", "/api/user/schedule-periods", map[string]any{
		"scheduleId":  1,
		"periodStart": lastWeek.Format(time.RFC3339),
		"periodEnd":   yesterday.Format(time.RFC3339),
		"type":        "fixed_date",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create period: %d", w.Code)
	}
	var period schedulemodels.SchedulePeriod
	json.Unmarshal(w.Body.Bytes(), &period)

	w = doJSONLog(t, r, "GET", fmt.Sprintf("/api/user/schedule-periods/%d/permissions", period.ID), nil)
	if w.Code != http.StatusOK {
		t.Fatalf("get permissions: %d", w.Code)
	}
	var perms shared.PermissionsResponse
	json.Unmarshal(w.Body.Bytes(), &perms)

	if len(perms.Permissions) != 1 {
		t.Fatalf("archived period should have 1 permission, got %d: %v", len(perms.Permissions), perms.Permissions)
	}
	if perms.Permissions[0] != shared.PermissionRead {
		t.Errorf("expected READ, got %s", perms.Permissions[0])
	}
}

// TestPeriodPermissions_NonOwnerForbidden verifies that another user
// cannot access period permissions.
func TestPeriodPermissions_NonOwnerForbidden(t *testing.T) {
	setupTestDB(t)
	defer closeDB(t)
	r := newRouter()

	doJSONLog(t, r, "POST", "/api/user/workouts", map[string]any{"name": "Forbidden Test"})
	doJSONLog(t, r, "POST", "/api/user/workout-schedules", map[string]any{"workoutId": 1})

	futureStart := time.Now().AddDate(0, 0, 7).Truncate(24 * time.Hour)
	futureEnd := time.Now().AddDate(0, 0, 14).Truncate(24 * time.Hour)
	w := doJSONLog(t, r, "POST", "/api/user/schedule-periods", map[string]any{
		"scheduleId":  1,
		"periodStart": futureStart.Format(time.RFC3339),
		"periodEnd":   futureEnd.Format(time.RFC3339),
		"type":        "fixed_date",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create period: %d", w.Code)
	}
	var period schedulemodels.SchedulePeriod
	json.Unmarshal(w.Body.Bytes(), &period)

	// Bob tries to access alice's period permissions
	w = doJSONLogAs(t, r, "GET", fmt.Sprintf("/api/user/schedule-periods/%d/permissions", period.ID), nil, "bob")
	if w.Code != http.StatusForbidden {
		t.Errorf("non-owner should get 403, got %d", w.Code)
	}
}

// TestPeriodList_BySchedule_OnlyOwnerPeriods verifies that listing periods
// by scheduleId only returns periods for the requesting user's schedule.
func TestPeriodList_BySchedule_OnlyOwnerPeriods(t *testing.T) {
	setupTestDB(t)
	defer closeDB(t)
	r := newRouter()

	// Alice creates a workout + schedule + period
	doJSONLog(t, r, "POST", "/api/user/workouts", map[string]any{"name": "Alice Workout"})
	startDate := time.Now().AddDate(0, 0, -3)
	w := doJSONLog(t, r, "POST", "/api/user/workout-schedules", map[string]any{
		"workoutId": 1, "startDate": startDate.Format(time.RFC3339),
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create schedule: %d", w.Code)
	}
	var aliceSchedule schedulemodels.WorkoutSchedule
	json.Unmarshal(w.Body.Bytes(), &aliceSchedule)

	futureStart := time.Now().AddDate(0, 0, 7).Truncate(24 * time.Hour)
	futureEnd := time.Now().AddDate(0, 0, 14).Truncate(24 * time.Hour)
	doJSONLog(t, r, "POST", "/api/user/schedule-periods", map[string]any{
		"scheduleId":  aliceSchedule.ID,
		"periodStart": futureStart.Format(time.RFC3339),
		"periodEnd":   futureEnd.Format(time.RFC3339),
		"type":        "fixed_date",
	})

	// Bob tries to list Alice's schedule periods → 403
	w = doJSONLogAs(t, r, "GET", fmt.Sprintf("/api/user/schedule-periods?scheduleId=%d", aliceSchedule.ID), nil, "bob")
	if w.Code != http.StatusForbidden {
		t.Errorf("bob listing alice's schedule periods should get 403, got %d", w.Code)
	}
}

// TestPeriodList_AllPeriods_OnlyReturnsOwnPeriods verifies that listing all
// periods (no scheduleId) only returns the requesting user's periods.
func TestPeriodList_AllPeriods_OnlyReturnsOwnPeriods(t *testing.T) {
	setupTestDB(t)
	defer closeDB(t)
	r := newRouter()

	// Alice creates a workout + schedule + period
	doJSONLog(t, r, "POST", "/api/user/workouts", map[string]any{"name": "Alice Workout"})
	startDate := time.Now().AddDate(0, 0, -3)
	doJSONLog(t, r, "POST", "/api/user/workout-schedules", map[string]any{
		"workoutId": 1, "startDate": startDate.Format(time.RFC3339),
	})

	futureStart := time.Now().AddDate(0, 0, 7).Truncate(24 * time.Hour)
	futureEnd := time.Now().AddDate(0, 0, 14).Truncate(24 * time.Hour)
	doJSONLog(t, r, "POST", "/api/user/schedule-periods", map[string]any{
		"scheduleId":  1,
		"periodStart": futureStart.Format(time.RFC3339),
		"periodEnd":   futureEnd.Format(time.RFC3339),
		"type":        "fixed_date",
	})

	// Bob creates his own workout + schedule + period
	doJSONLogAs(t, r, "POST", "/api/user/workouts", map[string]any{"name": "Bob Workout"}, "bob")
	doJSONLogAs(t, r, "POST", "/api/user/workout-schedules", map[string]any{
		"workoutId": 2, "startDate": startDate.Format(time.RFC3339),
	}, "bob")

	bobStart := time.Now().AddDate(0, 0, 3).Truncate(24 * time.Hour)
	bobEnd := time.Now().AddDate(0, 0, 10).Truncate(24 * time.Hour)
	doJSONLogAs(t, r, "POST", "/api/user/schedule-periods", map[string]any{
		"scheduleId":  2,
		"periodStart": bobStart.Format(time.RFC3339),
		"periodEnd":   bobEnd.Format(time.RFC3339),
		"type":        "fixed_date",
	}, "bob")

	// Alice lists all periods (no scheduleId) → only sees her own
	w := doJSONLog(t, r, "GET", "/api/user/schedule-periods", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("list all periods: %d", w.Code)
	}
	var alicePeriods []schedulemodels.SchedulePeriod
	json.Unmarshal(w.Body.Bytes(), &alicePeriods)

	if len(alicePeriods) != 1 {
		t.Fatalf("alice should see 1 period, got %d", len(alicePeriods))
	}
	if alicePeriods[0].ScheduleID != 1 {
		t.Errorf("alice's period should belong to schedule 1, got %d", alicePeriods[0].ScheduleID)
	}

	// Bob lists all periods (no scheduleId) → only sees his own
	w = doJSONLogAs(t, r, "GET", "/api/user/schedule-periods", nil, "bob")
	if w.Code != http.StatusOK {
		t.Fatalf("bob list all periods: %d", w.Code)
	}
	var bobPeriods []schedulemodels.SchedulePeriod
	json.Unmarshal(w.Body.Bytes(), &bobPeriods)

	if len(bobPeriods) != 1 {
		t.Fatalf("bob should see 1 period, got %d", len(bobPeriods))
	}
	if bobPeriods[0].ScheduleID != 2 {
		t.Errorf("bob's period should belong to schedule 2, got %d", bobPeriods[0].ScheduleID)
	}
}

// TestCommitmentList_AllCommitments_OnlyReturnsOwnCommitments verifies that
// listing all commitments (no periodId) only returns the requesting user's commitments.
func TestCommitmentList_AllCommitments_OnlyReturnsOwnCommitments(t *testing.T) {
	setupTestDB(t)
	defer closeDB(t)
	r := newRouter()

	// Alice creates workout → schedule → period → commitment
	doJSONLog(t, r, "POST", "/api/user/workouts", map[string]any{"name": "Alice Workout"})
	startDate := time.Now().AddDate(0, 0, -3)
	doJSONLog(t, r, "POST", "/api/user/workout-schedules", map[string]any{
		"workoutId": 1, "startDate": startDate.Format(time.RFC3339),
	})

	futureStart := time.Now().AddDate(0, 0, 7).Truncate(24 * time.Hour)
	futureEnd := time.Now().AddDate(0, 0, 14).Truncate(24 * time.Hour)
	w := doJSONLog(t, r, "POST", "/api/user/schedule-periods", map[string]any{
		"scheduleId": 1, "periodStart": futureStart.Format(time.RFC3339),
		"periodEnd": futureEnd.Format(time.RFC3339), "type": "fixed_date",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create period: %d", w.Code)
	}
	var alicePeriod schedulemodels.SchedulePeriod
	json.Unmarshal(w.Body.Bytes(), &alicePeriod)

	commitDate := time.Now().AddDate(0, 0, 8).Truncate(24 * time.Hour)
	doJSONLog(t, r, "POST", "/api/user/schedule-commitments", map[string]any{
		"periodId": alicePeriod.ID, "date": commitDate.Format(time.RFC3339),
	})

	// Bob creates workout → schedule → period → commitment
	doJSONLogAs(t, r, "POST", "/api/user/workouts", map[string]any{"name": "Bob Workout"}, "bob")
	doJSONLogAs(t, r, "POST", "/api/user/workout-schedules", map[string]any{
		"workoutId": 2, "startDate": startDate.Format(time.RFC3339),
	}, "bob")

	bobStart := time.Now().AddDate(0, 0, 3).Truncate(24 * time.Hour)
	bobEnd := time.Now().AddDate(0, 0, 10).Truncate(24 * time.Hour)
	w = doJSONLogAs(t, r, "POST", "/api/user/schedule-periods", map[string]any{
		"scheduleId": 2, "periodStart": bobStart.Format(time.RFC3339),
		"periodEnd": bobEnd.Format(time.RFC3339), "type": "fixed_date",
	}, "bob")
	if w.Code != http.StatusCreated {
		t.Fatalf("bob create period: %d", w.Code)
	}
	var bobPeriod schedulemodels.SchedulePeriod
	json.Unmarshal(w.Body.Bytes(), &bobPeriod)

	bobDate := time.Now().AddDate(0, 0, 5).Truncate(24 * time.Hour)
	doJSONLogAs(t, r, "POST", "/api/user/schedule-commitments", map[string]any{
		"periodId": bobPeriod.ID, "date": bobDate.Format(time.RFC3339),
	}, "bob")

	// Alice lists all commitments → only sees hers
	w = doJSONLog(t, r, "GET", "/api/user/schedule-commitments", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("alice list all commitments: %d", w.Code)
	}
	var aliceCommitments []schedulemodels.ScheduleCommitment
	json.Unmarshal(w.Body.Bytes(), &aliceCommitments)
	if len(aliceCommitments) != 1 {
		t.Fatalf("alice should see 1 commitment, got %d", len(aliceCommitments))
	}
	if aliceCommitments[0].PeriodID != alicePeriod.ID {
		t.Errorf("alice's commitment should belong to her period %d, got %d", alicePeriod.ID, aliceCommitments[0].PeriodID)
	}

	// Bob lists all commitments → only sees his
	w = doJSONLogAs(t, r, "GET", "/api/user/schedule-commitments", nil, "bob")
	if w.Code != http.StatusOK {
		t.Fatalf("bob list all commitments: %d", w.Code)
	}
	var bobCommitments []schedulemodels.ScheduleCommitment
	json.Unmarshal(w.Body.Bytes(), &bobCommitments)
	if len(bobCommitments) != 1 {
		t.Fatalf("bob should see 1 commitment, got %d", len(bobCommitments))
	}
	if bobCommitments[0].PeriodID != bobPeriod.ID {
		t.Errorf("bob's commitment should belong to his period %d, got %d", bobPeriod.ID, bobCommitments[0].PeriodID)
	}
}
