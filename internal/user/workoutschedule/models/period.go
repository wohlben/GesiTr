package models

import (
	"time"

	"gesitr/internal/shared"
)

// ScheduleType determines how commitments within a period are interpreted.
type ScheduleType string

const (
	// ScheduleTypeFixedDate means the user picks specific days within the period.
	// Each commitment has a Date set.
	ScheduleTypeFixedDate ScheduleType = "fixed_date"

	// ScheduleTypeFrequency means the user sets a count (N commitments per period)
	// without picking specific days. Commitments have Date=null.
	ScheduleTypeFrequency ScheduleType = "frequency"
)

// PeriodMode determines how the period duration is computed when cloning.
type PeriodMode string

const (
	// PeriodModeNormal clones the period with the same duration in days.
	// Next period starts the day after the previous one ended.
	PeriodModeNormal PeriodMode = "normal"

	// PeriodModeMonthly clones the period by advancing one calendar month.
	// The day-of-month offset is preserved (e.g. 4th Jan → 4th Feb).
	PeriodModeMonthly PeriodMode = "monthly"
)

// SchedulePeriod is the API-facing DTO for a schedule period.
type SchedulePeriod struct {
	shared.BaseModel `tstype:",extends"`
	ScheduleID       uint         `json:"scheduleId"`
	PeriodStart      time.Time    `json:"periodStart"`
	PeriodEnd        time.Time    `json:"periodEnd"`
	Type             ScheduleType `json:"type"`
	Mode             PeriodMode   `json:"mode"`
}
